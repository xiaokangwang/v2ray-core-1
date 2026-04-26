//go:build linux

package raw

import (
	"net"
	"net/netip"
	"strings"
	"testing"

	"github.com/v2fly/v2ray-core/v5/common/sendas"
	"golang.org/x/sys/unix"
)

type stubRawIPConn struct {
	writes int
	closed bool
	addrs  []string
}

func (c *stubRawIPConn) WriteToIP(b []byte, addr *net.IPAddr) (int, error) {
	c.writes++
	if addr != nil {
		c.addrs = append(c.addrs, addr.String())
	}
	return len(b), nil
}

func (c *stubRawIPConn) Close() error {
	c.closed = true
	return nil
}

func TestUDPPacketSenderCachesRawIPv4Socket(t *testing.T) {
	originalListen := listenRawIPFunc
	originalInterfaceLookup := netInterfaceByNameFunc
	t.Cleanup(func() {
		listenRawIPFunc = originalListen
		netInterfaceByNameFunc = originalInterfaceLookup
	})

	var listenCalls int
	var resolveCalls int
	conn := &stubRawIPConn{}

	listenRawIPFunc = func(network string, laddr *net.IPAddr, configure func(fd int) error) (rawIPConn, error) {
		listenCalls++
		if network != "ip4:udp" {
			t.Fatalf("unexpected raw socket network: %s", network)
		}
		return conn, nil
	}
	netInterfaceByNameFunc = func(name string) (*net.Interface, error) {
		resolveCalls++
		return &net.Interface{Name: name, Index: 7}, nil
	}

	sender := NewUDPPacketSender(&sendas.Config{
		Interface: "eth0",
		Method:    sendas.Method_RAW,
	})

	src := netip.MustParseAddrPort("198.51.100.10:12345")
	dst := netip.MustParseAddrPort("203.0.113.20:53")
	if err := sender.SendTo(src, dst, nil, nil, []byte("first")); err != nil {
		t.Fatalf("first SendTo failed: %v", err)
	}
	if err := sender.SendTo(src, dst, nil, nil, []byte("second")); err != nil {
		t.Fatalf("second SendTo failed: %v", err)
	}

	if listenCalls != 1 {
		t.Fatalf("unexpected raw socket opens: got %d want 1", listenCalls)
	}
	if resolveCalls != 1 {
		t.Fatalf("unexpected interface resolutions: got %d want 1", resolveCalls)
	}
	if conn.writes != 2 {
		t.Fatalf("unexpected raw socket writes: got %d want 2", conn.writes)
	}
	if len(conn.addrs) != 2 || conn.addrs[0] != "203.0.113.20" || conn.addrs[1] != "203.0.113.20" {
		t.Fatalf("unexpected raw route destinations: got %v", conn.addrs)
	}

	if err := sender.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	if !conn.closed {
		t.Fatal("expected cached raw socket to be closed")
	}
}

func TestUDPPacketSenderCachesEthernetSocket(t *testing.T) {
	originalInterfaceLookup := netInterfaceByNameFunc
	originalOpenSocket := openRawEthernetSocketFunc
	originalCloseSocket := closeRawEthernetSocketFunc
	originalSendFrame := sendRawEthernetFrameFunc
	t.Cleanup(func() {
		netInterfaceByNameFunc = originalInterfaceLookup
		openRawEthernetSocketFunc = originalOpenSocket
		closeRawEthernetSocketFunc = originalCloseSocket
		sendRawEthernetFrameFunc = originalSendFrame
	})

	var resolveCalls int
	var openCalls int
	var sendCalls int
	var closedFDs []int
	var usedFD int

	netInterfaceByNameFunc = func(name string) (*net.Interface, error) {
		resolveCalls++
		return &net.Interface{
			Name:         name,
			Index:        11,
			HardwareAddr: []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x01},
		}, nil
	}
	openRawEthernetSocketFunc = func(etherType uint16) (int, error) {
		openCalls++
		if etherType != etherTypeIPv4 {
			t.Fatalf("unexpected etherType: %#x", etherType)
		}
		return 100 + openCalls, nil
	}
	closeRawEthernetSocketFunc = func(fd int) error {
		closedFDs = append(closedFDs, fd)
		return nil
	}
	sendRawEthernetFrameFunc = func(fd int, frame []byte, sockaddr *unix.SockaddrLinklayer) error {
		sendCalls++
		if usedFD == 0 {
			usedFD = fd
		}
		if fd != usedFD {
			t.Fatalf("expected ethernet sends to reuse the same fd, got %d then %d", usedFD, fd)
		}
		return nil
	}

	sender := NewUDPPacketSender(&sendas.Config{
		Interface: "eth0",
		Method:    sendas.Method_RAW_ETHERNET,
	})

	src := netip.MustParseAddrPort("198.51.100.10:12345")
	dst := netip.MustParseAddrPort("203.0.113.20:53")
	dstMAC := []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x02}
	if err := sender.SendTo(src, dst, nil, dstMAC, []byte("first")); err != nil {
		t.Fatalf("first SendTo failed: %v", err)
	}
	if err := sender.SendTo(src, dst, nil, dstMAC, []byte("second")); err != nil {
		t.Fatalf("second SendTo failed: %v", err)
	}

	if resolveCalls != 1 {
		t.Fatalf("unexpected interface resolutions: got %d want 1", resolveCalls)
	}
	if openCalls != 1 {
		t.Fatalf("unexpected ethernet socket opens: got %d want 1", openCalls)
	}
	if sendCalls != 2 {
		t.Fatalf("unexpected ethernet sends: got %d want 2", sendCalls)
	}

	if err := sender.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	if len(closedFDs) != 1 || closedFDs[0] != usedFD {
		t.Fatalf("unexpected ethernet socket closes: got %v want [%d]", closedFDs, usedFD)
	}
}

func TestUDPPacketSenderUsesPerPacketGatewayWithoutReopeningSocket(t *testing.T) {
	originalListen := listenRawIPFunc
	t.Cleanup(func() {
		listenRawIPFunc = originalListen
	})

	var listenCalls int
	conn := &stubRawIPConn{}
	listenRawIPFunc = func(network string, laddr *net.IPAddr, configure func(fd int) error) (rawIPConn, error) {
		listenCalls++
		return conn, nil
	}

	sender := NewUDPPacketSender(&sendas.Config{
		Method: sendas.Method_RAW,
	})

	src := netip.MustParseAddrPort("198.51.100.10:12345")
	dst := netip.MustParseAddrPort("203.0.113.20:53")
	gatewayA := netip.MustParseAddr("192.0.2.1")
	gatewayB := netip.MustParseAddr("192.0.2.2")
	if err := sender.SendToWithGateway(src, dst, gatewayA, nil, nil, []byte("first")); err != nil {
		t.Fatalf("first SendToWithGateway failed: %v", err)
	}
	if err := sender.SendToWithGateway(src, dst, gatewayB, nil, nil, []byte("second")); err != nil {
		t.Fatalf("second SendToWithGateway failed: %v", err)
	}

	if listenCalls != 1 {
		t.Fatalf("unexpected raw socket opens: got %d want 1", listenCalls)
	}
	if conn.writes != 2 {
		t.Fatalf("unexpected raw socket writes: got %d want 2", conn.writes)
	}
	if len(conn.addrs) != 2 || conn.addrs[0] != "192.0.2.1" || conn.addrs[1] != "192.0.2.2" {
		t.Fatalf("unexpected raw route destinations: got %v", conn.addrs)
	}
}

func TestUDPPacketSenderRejectsGatewayOnRawEthernet(t *testing.T) {
	sender := NewUDPPacketSender(&sendas.Config{
		Method: sendas.Method_RAW_ETHERNET,
	})

	src := netip.MustParseAddrPort("198.51.100.10:12345")
	dst := netip.MustParseAddrPort("203.0.113.20:53")
	gateway := netip.MustParseAddr("192.0.2.1")
	err := sender.SendToWithGateway(src, dst, gateway, nil, nil, []byte("payload"))
	if err == nil {
		t.Fatal("expected raw ethernet gateway use to fail")
	}
	if !strings.Contains(err.Error(), "gateway") {
		t.Fatalf("unexpected raw ethernet gateway error: %v", err)
	}
}

func TestUDPPacketSenderRejectsSendAfterClose(t *testing.T) {
	originalListen := listenRawIPFunc
	t.Cleanup(func() {
		listenRawIPFunc = originalListen
	})

	conn := &stubRawIPConn{}
	var listenCalls int
	listenRawIPFunc = func(network string, laddr *net.IPAddr, configure func(fd int) error) (rawIPConn, error) {
		listenCalls++
		return conn, nil
	}

	sender := NewUDPPacketSender(&sendas.Config{
		Method: sendas.Method_RAW,
	})

	src := netip.MustParseAddrPort("198.51.100.10:12345")
	dst := netip.MustParseAddrPort("203.0.113.20:53")
	if err := sender.SendTo(src, dst, nil, nil, []byte("first")); err != nil {
		t.Fatalf("first SendTo failed: %v", err)
	}
	if err := sender.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	err := sender.SendTo(src, dst, nil, nil, []byte("second"))
	if err == nil {
		t.Fatal("expected SendTo after Close to fail")
	}
	if !strings.Contains(err.Error(), "closed") {
		t.Fatalf("unexpected SendTo error after Close: %v", err)
	}
	if listenCalls != 1 {
		t.Fatalf("unexpected raw socket opens after Close: got %d want 1", listenCalls)
	}
}
