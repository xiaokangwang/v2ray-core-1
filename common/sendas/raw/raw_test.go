package raw

import (
	"bytes"
	"net/netip"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func TestBuildIPv4UDPPacket(t *testing.T) {
	src := netip.MustParseAddrPort("198.51.100.10:12345")
	dst := netip.MustParseAddrPort("203.0.113.20:53")
	payload := []byte("payload")

	packet, etherType, err := buildUDPPacket(src, dst, payload)
	if err != nil {
		t.Fatalf("buildUDPPacket failed: %v", err)
	}
	if etherType != etherTypeIPv4 {
		t.Fatalf("unexpected etherType: got %#x want %#x", etherType, etherTypeIPv4)
	}

	parsed := gopacket.NewPacket(packet, layers.LayerTypeIPv4, gopacket.Default)
	ipLayer := parsed.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		t.Fatalf("missing IPv4 layer: %v", parsed.ErrorLayer())
	}
	ip := ipLayer.(*layers.IPv4)
	if got := ip.SrcIP.String(); got != "198.51.100.10" {
		t.Fatalf("unexpected IPv4 source: %s", got)
	}
	if got := ip.DstIP.String(); got != "203.0.113.20" {
		t.Fatalf("unexpected IPv4 destination: %s", got)
	}
	if ip.Protocol != layers.IPProtocolUDP {
		t.Fatalf("unexpected IPv4 protocol: %v", ip.Protocol)
	}

	udpLayer := parsed.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		t.Fatalf("missing UDP layer: %v", parsed.ErrorLayer())
	}
	udp := udpLayer.(*layers.UDP)
	if got := uint16(udp.SrcPort); got != 12345 {
		t.Fatalf("unexpected UDP source port: %d", got)
	}
	if got := uint16(udp.DstPort); got != 53 {
		t.Fatalf("unexpected UDP destination port: %d", got)
	}
	if !bytes.Equal(udp.Payload, payload) {
		t.Fatalf("unexpected UDP payload: %q", udp.Payload)
	}
}

func TestBuildIPv6UDPPacket(t *testing.T) {
	src := netip.MustParseAddrPort("[2001:db8::10]:12345")
	dst := netip.MustParseAddrPort("[2001:db8::20]:53")
	payload := []byte("payload")

	packet, etherType, err := buildUDPPacket(src, dst, payload)
	if err != nil {
		t.Fatalf("buildUDPPacket failed: %v", err)
	}
	if etherType != etherTypeIPv6 {
		t.Fatalf("unexpected etherType: got %#x want %#x", etherType, etherTypeIPv6)
	}

	parsed := gopacket.NewPacket(packet, layers.LayerTypeIPv6, gopacket.Default)
	ipLayer := parsed.Layer(layers.LayerTypeIPv6)
	if ipLayer == nil {
		t.Fatalf("missing IPv6 layer: %v", parsed.ErrorLayer())
	}
	ip := ipLayer.(*layers.IPv6)
	if got := ip.SrcIP.String(); got != "2001:db8::10" {
		t.Fatalf("unexpected IPv6 source: %s", got)
	}
	if got := ip.DstIP.String(); got != "2001:db8::20" {
		t.Fatalf("unexpected IPv6 destination: %s", got)
	}
	if ip.NextHeader != layers.IPProtocolUDP {
		t.Fatalf("unexpected IPv6 next header: %v", ip.NextHeader)
	}

	udpLayer := parsed.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		t.Fatalf("missing UDP layer: %v", parsed.ErrorLayer())
	}
	udp := udpLayer.(*layers.UDP)
	if got := uint16(udp.SrcPort); got != 12345 {
		t.Fatalf("unexpected UDP source port: %d", got)
	}
	if got := uint16(udp.DstPort); got != 53 {
		t.Fatalf("unexpected UDP destination port: %d", got)
	}
	if !bytes.Equal(udp.Payload, payload) {
		t.Fatalf("unexpected UDP payload: %q", udp.Payload)
	}
	if udp.Checksum == 0 {
		t.Fatal("expected a non-zero IPv6 UDP checksum")
	}
}

func TestBuildEthernetFrame(t *testing.T) {
	src := netip.MustParseAddrPort("198.51.100.10:12345")
	dst := netip.MustParseAddrPort("203.0.113.20:53")
	srcMac := []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x01}
	dstMac := []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x02}
	payload := []byte("payload")

	frame, _, err := buildEthernetFrame(src, dst, srcMac, dstMac, payload)
	if err != nil {
		t.Fatalf("buildEthernetFrame failed: %v", err)
	}

	parsed := gopacket.NewPacket(frame, layers.LayerTypeEthernet, gopacket.Default)
	ethernetLayer := parsed.Layer(layers.LayerTypeEthernet)
	if ethernetLayer == nil {
		t.Fatalf("missing ethernet layer: %v", parsed.ErrorLayer())
	}
	ethernet := ethernetLayer.(*layers.Ethernet)
	if !bytes.Equal(ethernet.SrcMAC, srcMac) {
		t.Fatalf("unexpected source MAC: %v", ethernet.SrcMAC)
	}
	if !bytes.Equal(ethernet.DstMAC, dstMac) {
		t.Fatalf("unexpected destination MAC: %v", ethernet.DstMAC)
	}
	if ethernet.EthernetType != layers.EthernetTypeIPv4 {
		t.Fatalf("unexpected ethernet type: %v", ethernet.EthernetType)
	}
}
