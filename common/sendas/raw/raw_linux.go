//go:build linux

package raw

import (
	"encoding/binary"
	"net"
	"net/netip"
	"sync"
	"unsafe"

	"github.com/v2fly/v2ray-core/v5/common/errors"
	"golang.org/x/sys/unix"
)

type rawIPConn interface {
	WriteToIP(b []byte, addr *net.IPAddr) (int, error)
	Close() error
}

var (
	listenRawIPFunc = func(network string, laddr *net.IPAddr, configure func(fd int) error) (rawIPConn, error) {
		return listenRawIP(network, laddr, configure)
	}
	netInterfaceByNameFunc    = net.InterfaceByName
	openRawEthernetSocketFunc = func(etherType uint16) (int, error) {
		return unix.Socket(unix.AF_PACKET, unix.SOCK_RAW|unix.SOCK_CLOEXEC, int(htons(etherType)))
	}
	closeRawEthernetSocketFunc = unix.Close
	sendRawEthernetFrameFunc   = func(fd int, frame []byte, sockaddr *unix.SockaddrLinklayer) error {
		return unix.Sendto(fd, frame, 0, sockaddr)
	}
)

type rawIPSocketKey struct {
	network       string
	interfaceName string
}

type rawIPSocket struct {
	conn  rawIPConn
	iface *net.Interface
}

type rawEthernetSocketKey struct {
	interfaceName string
	etherType     uint16
}

type rawEthernetSocket struct {
	fd            int
	iface         *net.Interface
	defaultSrcMAC [6]byte
	hasDefaultMAC bool
}

type linuxUDPPacketSenderState struct {
	mu              sync.Mutex
	closed          bool
	rawSockets      map[rawIPSocketKey]*rawIPSocket
	ethernetSockets map[rawEthernetSocketKey]*rawEthernetSocket
	interfaces      map[string]*net.Interface
}

func newUDPPacketSenderState() udpPacketSenderState {
	return &linuxUDPPacketSenderState{}
}

func (s *linuxUDPPacketSenderState) sendRaw(interfaceName string, src netip.AddrPort, dst netip.AddrPort, gateway netip.Addr, payload []byte) error {
	packet, _, err := buildUDPPacket(src, dst, payload)
	if err != nil {
		return err
	}

	routeAddr := dst.Addr()
	if gateway.IsValid() {
		routeAddr = gateway
	}

	network := "ip6:udp"
	if dst.Addr().Is4() {
		network = "ip4:udp"
	}

	socket, err := s.rawIPSocket(network, interfaceName, routeAddr.Zone())
	if err != nil {
		return err
	}

	if dst.Addr().Is4() {
		if _, err := socket.conn.WriteToIP(packet, &net.IPAddr{IP: routeAddr.AsSlice()}); err != nil {
			return newError("failed to send IPv4 raw UDP packet").Base(err)
		}

		return nil
	}

	writeAddr, err := ipv6WriteAddr(routeAddr, socket.iface)
	if err != nil {
		return err
	}

	if _, err := socket.conn.WriteToIP(packet, writeAddr); err != nil {
		return newError("failed to send IPv6 raw UDP packet").Base(err)
	}

	return nil
}

func (s *linuxUDPPacketSenderState) sendRawEthernet(interfaceName string, src netip.AddrPort, dst netip.AddrPort, srcMac []byte, dstMac []byte, payload []byte) error {
	dstHW, err := normalizeMAC(dstMac)
	if err != nil {
		return newError("invalid destination MAC").Base(err)
	}

	etherType := uint16(etherTypeIPv6)
	if dst.Addr().Is4() {
		etherType = etherTypeIPv4
	}

	socket, err := s.rawEthernetSocket(interfaceName, etherType)
	if err != nil {
		return err
	}

	srcHW, err := socket.sourceMAC(srcMac)
	if err != nil {
		return newError("invalid source MAC").Base(err)
	}

	frame, _, err := buildEthernetFrameWithNormalizedMACs(src, dst, srcHW, dstHW, payload)
	if err != nil {
		return err
	}

	sockaddr := &unix.SockaddrLinklayer{
		Protocol: htons(etherType),
		Ifindex:  socket.iface.Index,
		Halen:    6,
	}
	copy(sockaddr.Addr[:], dstHW[:])

	if err := sendRawEthernetFrameFunc(socket.fd, frame, sockaddr); err != nil {
		return newError("failed to send ethernet UDP frame").Base(err)
	}

	return nil
}

func (s *linuxUDPPacketSenderState) close() error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}

	s.closed = true
	rawSockets := s.rawSockets
	ethernetSockets := s.ethernetSockets
	s.rawSockets = nil
	s.ethernetSockets = nil
	s.interfaces = nil
	s.mu.Unlock()

	var closeErrs []error
	for _, socket := range rawSockets {
		if err := socket.conn.Close(); err != nil {
			closeErrs = append(closeErrs, newError("failed to close raw IP socket").Base(err))
		}
	}
	for _, socket := range ethernetSockets {
		if err := closeRawEthernetSocketFunc(socket.fd); err != nil {
			closeErrs = append(closeErrs, newError("failed to close AF_PACKET raw socket").Base(err))
		}
	}

	return errors.Combine(closeErrs...)
}

func (s *linuxUDPPacketSenderState) rawIPSocket(network string, interfaceName string, fallbackZone string) (*rawIPSocket, error) {
	iface, boundInterface, err := s.resolveInterface(interfaceName, fallbackZone)
	if err != nil {
		return nil, err
	}

	key := rawIPSocketKey{
		network:       network,
		interfaceName: boundInterface,
	}

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil, newError("UDP packet sender is closed")
	}
	if socket := s.rawSockets[key]; socket != nil {
		s.mu.Unlock()
		return socket, nil
	}
	s.mu.Unlock()

	conn, err := listenRawIPFunc(network, &net.IPAddr{}, func(fd int) error {
		switch network {
		case "ip4:udp":
			if err := unix.SetsockoptInt(fd, unix.IPPROTO_IP, unix.IP_HDRINCL, 1); err != nil {
				return newError("failed to set IP_HDRINCL").Base(err)
			}
			if err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_BROADCAST, 1); err != nil {
				return newError("failed to set SO_BROADCAST").Base(err)
			}
		case "ip6:udp":
			if err := unix.SetsockoptInt(fd, unix.IPPROTO_IPV6, unix.IPV6_HDRINCL, 1); err != nil {
				return newError("failed to set IPV6_HDRINCL").Base(err)
			}
		default:
			return newError("unsupported raw socket network: ", network)
		}

		return bindToDevice(fd, boundInterface)
	})
	if err != nil {
		return nil, err
	}

	socket := &rawIPSocket{
		conn:  conn,
		iface: iface,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		if err := conn.Close(); err != nil {
			return nil, newError("UDP packet sender is closed").Base(err)
		}
		return nil, newError("UDP packet sender is closed")
	}
	if existing := s.rawSockets[key]; existing != nil {
		if err := conn.Close(); err != nil {
			return nil, newError("failed to close redundant raw IP socket").Base(err)
		}
		return existing, nil
	}
	if s.rawSockets == nil {
		s.rawSockets = make(map[rawIPSocketKey]*rawIPSocket)
	}
	s.rawSockets[key] = socket

	return socket, nil
}

func (s *linuxUDPPacketSenderState) rawEthernetSocket(interfaceName string, etherType uint16) (*rawEthernetSocket, error) {
	iface, _, err := s.resolveInterface(interfaceName, "")
	if err != nil {
		return nil, err
	}
	if iface == nil {
		return nil, newError("raw ethernet sender requires an interface")
	}

	key := rawEthernetSocketKey{
		interfaceName: iface.Name,
		etherType:     etherType,
	}

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil, newError("UDP packet sender is closed")
	}
	if socket := s.ethernetSockets[key]; socket != nil {
		s.mu.Unlock()
		return socket, nil
	}
	s.mu.Unlock()

	fd, err := openRawEthernetSocketFunc(etherType)
	if err != nil {
		return nil, newError("failed to create AF_PACKET raw socket").Base(err)
	}

	socket := &rawEthernetSocket{
		fd:    fd,
		iface: iface,
	}
	if normalized, err := normalizeMAC(iface.HardwareAddr); err == nil {
		socket.defaultSrcMAC = normalized
		socket.hasDefaultMAC = true
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		if closeErr := closeRawEthernetSocketFunc(fd); closeErr != nil {
			return nil, newError("UDP packet sender is closed").Base(closeErr)
		}
		return nil, newError("UDP packet sender is closed")
	}
	if existing := s.ethernetSockets[key]; existing != nil {
		if closeErr := closeRawEthernetSocketFunc(fd); closeErr != nil {
			return nil, newError("failed to close redundant AF_PACKET raw socket").Base(closeErr)
		}
		return existing, nil
	}
	if s.ethernetSockets == nil {
		s.ethernetSockets = make(map[rawEthernetSocketKey]*rawEthernetSocket)
	}
	s.ethernetSockets[key] = socket

	return socket, nil
}

func (s *linuxUDPPacketSenderState) resolveInterface(interfaceName string, fallbackZone string) (*net.Interface, string, error) {
	name := interfaceName
	if name == "" {
		name = fallbackZone
	}
	if name == "" {
		return nil, "", nil
	}

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil, "", newError("UDP packet sender is closed")
	}
	if iface := s.interfaces[name]; iface != nil {
		s.mu.Unlock()
		return iface, iface.Name, nil
	}
	s.mu.Unlock()

	iface, err := netInterfaceByNameFunc(name)
	if err != nil {
		return nil, "", newError("failed to resolve interface ", name).Base(err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil, "", newError("UDP packet sender is closed")
	}
	if existing := s.interfaces[name]; existing != nil {
		return existing, existing.Name, nil
	}
	if s.interfaces == nil {
		s.interfaces = make(map[string]*net.Interface)
	}
	s.interfaces[name] = iface

	return iface, iface.Name, nil
}

func (s *rawEthernetSocket) sourceMAC(srcMac []byte) ([6]byte, error) {
	if len(srcMac) != 0 {
		return normalizeMAC(srcMac)
	}
	if s.hasDefaultMAC {
		return s.defaultSrcMAC, nil
	}
	return normalizeMAC(s.iface.HardwareAddr)
}

func bindToDevice(fd int, interfaceName string) error {
	if interfaceName == "" {
		return nil
	}
	if err := unix.BindToDevice(fd, interfaceName); err != nil {
		return newError("failed to bind raw socket to interface ", interfaceName).Base(err)
	}
	return nil
}

func listenRawIP(network string, laddr *net.IPAddr, configure func(fd int) error) (*net.IPConn, error) {
	conn, err := net.ListenIP(network, laddr)
	if err != nil {
		return nil, newError("failed to create ", network, " raw socket").Base(err)
	}

	rawConn, err := conn.SyscallConn()
	if err != nil {
		conn.Close()
		return nil, newError("failed to access raw socket for ", network).Base(err)
	}

	var controlErr error
	if err := rawConn.Control(func(fd uintptr) {
		controlErr = configure(int(fd))
	}); err != nil {
		conn.Close()
		return nil, newError("failed to configure ", network, " raw socket").Base(err)
	}
	if controlErr != nil {
		conn.Close()
		return nil, controlErr
	}

	return conn, nil
}

func ipv6WriteAddr(addr netip.Addr, iface *net.Interface) (*net.IPAddr, error) {
	writeAddr := &net.IPAddr{IP: addr.AsSlice()}
	if zone := addr.Zone(); zone != "" {
		if iface != nil && iface.Name == zone {
			writeAddr.Zone = zone
			return writeAddr, nil
		}
		if _, err := netInterfaceByNameFunc(zone); err != nil {
			return nil, newError("failed to resolve IPv6 zone ", zone).Base(err)
		}
		writeAddr.Zone = zone
		return writeAddr, nil
	}

	if iface != nil && needsIPv6Zone(addr) {
		writeAddr.Zone = iface.Name
	}

	return writeAddr, nil
}

func needsIPv6Zone(addr netip.Addr) bool {
	return addr.IsLinkLocalUnicast() || addr.IsInterfaceLocalMulticast() || addr.IsLinkLocalMulticast()
}

func htons(value uint16) uint16 {
	if nativeEndian == binary.BigEndian {
		return value
	}
	return (value << 8) | (value >> 8)
}

var nativeEndian = func() binary.ByteOrder {
	var value uint16 = 0x0102
	if *(*byte)(unsafe.Pointer(&value)) == 0x01 {
		return binary.BigEndian
	}
	return binary.LittleEndian
}()
