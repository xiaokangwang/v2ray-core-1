package raw

import (
	"encoding/binary"
	"net/netip"

	"github.com/v2fly/v2ray-core/v5/common/errors"
	"github.com/v2fly/v2ray-core/v5/common/sendas"
)

const (
	ipv4HeaderLen = 20
	ipv6HeaderLen = 40
	udpHeaderLen  = 8
	defaultHopTTL = 64

	etherTypeIPv4 = 0x0800
	etherTypeIPv6 = 0x86DD
)

type errPathObjHolder struct{}

func newError(values ...interface{}) *errors.Error {
	return errors.New(values...).WithPathObj(errPathObjHolder{})
}

type udpPacketSenderState interface {
	sendRaw(interfaceName string, src netip.AddrPort, dst netip.AddrPort, gateway netip.Addr, payload []byte) error
	sendRawEthernet(interfaceName string, src netip.AddrPort, dst netip.AddrPort, srcMac []byte, dstMac []byte, payload []byte) error
	close() error
}

type UDPPacketSender struct {
	config *sendas.Config
	state  udpPacketSenderState
}

func NewUDPPacketSender(config *sendas.Config) *UDPPacketSender {
	return &UDPPacketSender{
		config: config,
		state:  newUDPPacketSenderState(),
	}
}

func (u *UDPPacketSender) Close() error {
	if u == nil || u.state == nil {
		return nil
	}

	return u.state.close()
}

func (u *UDPPacketSender) SendTo(src netip.AddrPort, dst netip.AddrPort, srcMac []byte, dstMac []byte, msg []byte) error {
	return u.SendToWithGateway(src, dst, netip.Addr{}, srcMac, dstMac, msg)
}

func (u *UDPPacketSender) SendToWithGateway(src netip.AddrPort, dst netip.AddrPort, gateway netip.Addr, srcMac []byte, dstMac []byte, msg []byte) error {
	if u == nil {
		return newError("nil UDP packet sender")
	}
	if u.config == nil {
		return newError("nil sendas config")
	}

	src, dst, err := validateAddrPorts(src, dst)
	if err != nil {
		return err
	}
	gateway, err = validateGateway(dst, gateway)
	if err != nil {
		return err
	}

	if u.state == nil {
		u.state = newUDPPacketSenderState()
	}

	switch u.config.GetMethod() {
	case sendas.Method_RAW:
		return u.state.sendRaw(u.config.GetInterface(), src, dst, gateway, msg)
	case sendas.Method_RAW_ETHERNET:
		if gateway.IsValid() {
			return newError("gateway is not supported for raw ethernet sender")
		}
		return u.state.sendRawEthernet(u.config.GetInterface(), src, dst, srcMac, dstMac, msg)
	case sendas.Method_FAIL:
		return newError("sendas method is disabled")
	default:
		return newError("unsupported sendas method: ", u.config.GetMethod())
	}
}

func validateAddrPorts(src netip.AddrPort, dst netip.AddrPort) (netip.AddrPort, netip.AddrPort, error) {
	src = normalizeAddrPort(src)
	dst = normalizeAddrPort(dst)

	if !src.IsValid() {
		return netip.AddrPort{}, netip.AddrPort{}, newError("invalid source address")
	}
	if !dst.IsValid() {
		return netip.AddrPort{}, netip.AddrPort{}, newError("invalid destination address")
	}
	if src.Addr().Is4() != dst.Addr().Is4() {
		return netip.AddrPort{}, netip.AddrPort{}, newError("address family mismatch between source and destination")
	}

	return src, dst, nil
}

func validateGateway(dst netip.AddrPort, gateway netip.Addr) (netip.Addr, error) {
	if !gateway.IsValid() {
		return netip.Addr{}, nil
	}

	gateway = normalizeAddr(gateway)
	if gateway.Is4() != dst.Addr().Is4() {
		return netip.Addr{}, newError("address family mismatch between destination and gateway")
	}

	return gateway, nil
}

func normalizeAddrPort(addrPort netip.AddrPort) netip.AddrPort {
	if !addrPort.IsValid() {
		return addrPort
	}

	return netip.AddrPortFrom(normalizeAddr(addrPort.Addr()), addrPort.Port())
}

func normalizeAddr(addr netip.Addr) netip.Addr {
	if addr.Is4In6() {
		return addr.Unmap()
	}
	return addr
}

func buildUDPPacket(src netip.AddrPort, dst netip.AddrPort, payload []byte) ([]byte, uint16, error) {
	src, dst, err := validateAddrPorts(src, dst)
	if err != nil {
		return nil, 0, err
	}

	if src.Addr().Is4() {
		packet, buildErr := buildIPv4UDPPacket(src, dst, payload)
		return packet, etherTypeIPv4, buildErr
	}

	packet, buildErr := buildIPv6UDPPacket(src, dst, payload)
	return packet, etherTypeIPv6, buildErr
}

func buildEthernetFrame(src netip.AddrPort, dst netip.AddrPort, srcMac []byte, dstMac []byte, payload []byte) ([]byte, uint16, error) {
	srcHW, err := normalizeMAC(srcMac)
	if err != nil {
		return nil, 0, newError("invalid source MAC").Base(err)
	}
	dstHW, err := normalizeMAC(dstMac)
	if err != nil {
		return nil, 0, newError("invalid destination MAC").Base(err)
	}

	return buildEthernetFrameWithNormalizedMACs(src, dst, srcHW, dstHW, payload)
}

func buildEthernetFrameWithNormalizedMACs(src netip.AddrPort, dst netip.AddrPort, srcHW [6]byte, dstHW [6]byte, payload []byte) ([]byte, uint16, error) {
	ipPacket, etherType, err := buildUDPPacket(src, dst, payload)
	if err != nil {
		return nil, 0, err
	}

	frame := make([]byte, 14+len(ipPacket))
	copy(frame[0:6], dstHW[:])
	copy(frame[6:12], srcHW[:])
	binary.BigEndian.PutUint16(frame[12:14], etherType)
	copy(frame[14:], ipPacket)

	return frame, etherType, nil
}

func buildIPv4UDPPacket(src netip.AddrPort, dst netip.AddrPort, payload []byte) ([]byte, error) {
	udpLen := udpHeaderLen + len(payload)
	totalLen := ipv4HeaderLen + udpLen
	if totalLen > 0xFFFF {
		return nil, newError("IPv4 UDP packet too large: ", totalLen)
	}

	packet := make([]byte, totalLen)
	srcAddr := src.Addr().As4()
	dstAddr := dst.Addr().As4()

	packet[0] = 0x45
	packet[8] = defaultHopTTL
	packet[9] = 17
	binary.BigEndian.PutUint16(packet[2:4], uint16(totalLen))
	copy(packet[12:16], srcAddr[:])
	copy(packet[16:20], dstAddr[:])

	udp := packet[ipv4HeaderLen:]
	binary.BigEndian.PutUint16(udp[0:2], src.Port())
	binary.BigEndian.PutUint16(udp[2:4], dst.Port())
	binary.BigEndian.PutUint16(udp[4:6], uint16(udpLen))
	copy(udp[udpHeaderLen:], payload)

	pseudoHeader := make([]byte, 12)
	copy(pseudoHeader[0:4], srcAddr[:])
	copy(pseudoHeader[4:8], dstAddr[:])
	pseudoHeader[9] = 17
	binary.BigEndian.PutUint16(pseudoHeader[10:12], uint16(udpLen))
	binary.BigEndian.PutUint16(udp[6:8], transportChecksum(pseudoHeader, udp))
	binary.BigEndian.PutUint16(packet[10:12], internetChecksum(packet[:ipv4HeaderLen]))

	return packet, nil
}

func buildIPv6UDPPacket(src netip.AddrPort, dst netip.AddrPort, payload []byte) ([]byte, error) {
	udpLen := udpHeaderLen + len(payload)
	if udpLen > 0xFFFF {
		return nil, newError("IPv6 UDP payload too large: ", udpLen)
	}

	packet := make([]byte, ipv6HeaderLen+udpLen)
	srcAddr := src.Addr().As16()
	dstAddr := dst.Addr().As16()

	packet[0] = 0x60
	packet[6] = 17
	packet[7] = defaultHopTTL
	binary.BigEndian.PutUint16(packet[4:6], uint16(udpLen))
	copy(packet[8:24], srcAddr[:])
	copy(packet[24:40], dstAddr[:])

	udp := packet[ipv6HeaderLen:]
	binary.BigEndian.PutUint16(udp[0:2], src.Port())
	binary.BigEndian.PutUint16(udp[2:4], dst.Port())
	binary.BigEndian.PutUint16(udp[4:6], uint16(udpLen))
	copy(udp[udpHeaderLen:], payload)

	pseudoHeader := make([]byte, 40)
	copy(pseudoHeader[0:16], srcAddr[:])
	copy(pseudoHeader[16:32], dstAddr[:])
	binary.BigEndian.PutUint32(pseudoHeader[32:36], uint32(udpLen))
	pseudoHeader[39] = 17
	binary.BigEndian.PutUint16(udp[6:8], transportChecksum(pseudoHeader, udp))

	return packet, nil
}

func normalizeMAC(mac []byte) ([6]byte, error) {
	var normalized [6]byte
	if len(mac) != len(normalized) {
		return normalized, newError("expected 6-byte MAC address, got ", len(mac))
	}

	copy(normalized[:], mac)
	return normalized, nil
}

func internetChecksum(parts ...[]byte) uint16 {
	sum := checksumSum(parts...)
	for sum>>16 != 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum)
}

func transportChecksum(parts ...[]byte) uint16 {
	sum := checksumSum(parts...)
	for sum>>16 != 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	checksum := ^uint16(sum)
	if checksum == 0 {
		return 0xFFFF
	}
	return checksum
}

func checksumSum(parts ...[]byte) uint32 {
	var sum uint32

	for _, part := range parts {
		for i := 0; i+1 < len(part); i += 2 {
			sum += uint32(binary.BigEndian.Uint16(part[i : i+2]))
		}
		if len(part)%2 != 0 {
			sum += uint32(part[len(part)-1]) << 8
		}
	}

	return sum
}
