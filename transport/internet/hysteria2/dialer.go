package hysteria2

import (
	"context"

	hy_client "github.com/apernet/hysteria/core/client"
	hyProtocol "github.com/apernet/hysteria/core/international/protocol"
	"github.com/apernet/quic-go/quicvarint"

	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/session"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tls"
)

const (
	FrameTypeTCPRequest = 0x401
)

var RunningClient map[net.Destination](hy_client.Client)

func InitTLSConifg(streamSettings *internet.MemoryStreamConfig) (*hy_client.TLSConfig, error) {
	tlsSetting := CheckTLSConfig(streamSettings, true)
	if tlsSetting == nil {
		tlsSetting = &tls.Config{
			ServerName:    internalDomain,
			AllowInsecure: true,
		}
	}
	res := &hy_client.TLSConfig{
		ServerName:         tlsSetting.ServerName,
		InsecureSkipVerify: tlsSetting.AllowInsecure,
	}
	return res, nil
}

func InitAddress(dest net.Destination) (net.Addr, error) {
	var destAddr *net.UDPAddr
	if dest.Address.Family().IsIP() {
		destAddr = &net.UDPAddr{
			IP:   dest.Address.IP(),
			Port: int(dest.Port),
		}
	} else {
		addr, err := net.ResolveUDPAddr("udp", dest.NetAddr())
		if err != nil {
			return nil, err
		}
		destAddr = addr
	}
	return destAddr, nil
}

func NewHyClient(dest net.Destination, streamSettings *internet.MemoryStreamConfig) (hy_client.Client, error) {
	tlsConfig, err := InitTLSConifg(streamSettings)
	if err != nil {
		return nil, err
	}

	serverAddr, err := InitAddress(dest)
	if err != nil {
		return nil, err
	}

	config := streamSettings.ProtocolSettings.(*Config)
	client, _, err := hy_client.NewClient(&hy_client.Config{
		TLSConfig:  *tlsConfig,
		Auth:       config.GetPassword(),
		ServerAddr: serverAddr,
	})
	if err != nil {
		return nil, err
	}

	RunningClient[dest] = client
	return client, nil
}

func CloseHyClient(dest net.Destination) error {
	client, found := RunningClient[dest]
	if found {
		delete(RunningClient, dest)
		return client.Close()
	}
	return nil
}

func Dial(ctx context.Context, dest net.Destination, streamSettings *internet.MemoryStreamConfig) (internet.Connection, error) {
	config := streamSettings.ProtocolSettings.(*Config)

	var client hy_client.Client
	var err error
	client, found := RunningClient[dest]
	if !found {
		// TODO: Clean idle connections
		client, err = NewHyClient(dest, streamSettings)
		if err != nil {
			return nil, err
		}
	}

	quicConn := client.GetQuicConn()
	conn := &HyConn{
		local:  quicConn.LocalAddr(),
		remote: quicConn.RemoteAddr(),
	}

	outbound := session.OutboundFromContext(ctx)
	network := net.Network_TCP
	if outbound != nil {
		network = outbound.Target.Network
		conn.Target = outbound.Target
	}

	if network == net.Network_UDP && config.GetUseUdpExtension() { // only hysteria2 can use udpExtension
		conn.IsUDPExtension = true
		conn.IsServer = false
		conn.ClientUDPSession, err = client.UDP()
		if err != nil {
			CloseHyClient(dest)
			return nil, err
		}
		return conn, nil
	}

	conn.stream, err = client.OpenStream()
	if err != nil {
		CloseHyClient(dest)
		return nil, err
	}

	// write TCP frame type
	frameSize := int(quicvarint.Len(FrameTypeTCPRequest))
	buf := make([]byte, frameSize)
	hyProtocol.VarintPut(buf, FrameTypeTCPRequest)
	conn.stream.Write(buf)
	return conn, nil
}

func init() {
	RunningClient = make(map[net.Destination]hy_client.Client)
	common.Must(internet.RegisterTransportDialer(protocolName, Dial))
}
