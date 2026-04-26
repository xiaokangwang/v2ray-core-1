//go:build !linux

package raw

import "net/netip"

type unsupportedUDPPacketSenderState struct{}

func newUDPPacketSenderState() udpPacketSenderState {
	return unsupportedUDPPacketSenderState{}
}

func (unsupportedUDPPacketSenderState) sendRaw(_ string, _ netip.AddrPort, _ netip.AddrPort, _ netip.Addr, _ []byte) error {
	return newError("linux raw UDP sender is not supported on this platform")
}

func (unsupportedUDPPacketSenderState) sendRawEthernet(_ string, _ netip.AddrPort, _ netip.AddrPort, _ []byte, _ []byte, _ []byte) error {
	return newError("linux raw ethernet sender is not supported on this platform")
}

func (unsupportedUDPPacketSenderState) close() error {
	return nil
}
