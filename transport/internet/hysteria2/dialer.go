package hysteria2

import (
	"context"
	"sync"

	hyClient "github.com/apernet/hysteria/core/v2/client"
	hyProtocol "github.com/apernet/hysteria/core/v2/international/protocol"
	"github.com/apernet/quic-go/quicvarint"

	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/session"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tls"
)

var RunningClient map[net.Destination](hyClient.Client)
var ClientMutex sync.Mutex

func GetClientTLSConfig(streamSettings *internet.MemoryStreamConfig) (*hyClient.TLSConfig, error) {
	config := tls.ConfigFromStreamSettings(streamSettings)
	if config == nil {
		return nil, newError(Hy2MustNeedTLS)
	}
	tlsConfig := config.GetTLSConfig()

	return &hyClient.TLSConfig{
		RootCAs:               tlsConfig.RootCAs,
		ServerName:            tlsConfig.ServerName,
		InsecureSkipVerify:    tlsConfig.InsecureSkipVerify,
		VerifyPeerCertificate: tlsConfig.VerifyPeerCertificate,
	}, nil
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

type connFactory struct {
	hyClient.ConnFactory

	NewFunc func(addr net.Addr) (net.PacketConn, error)
}

func (f *connFactory) New(addr net.Addr) (net.PacketConn, error) {
	return f.NewFunc(addr)
}

func NewHyClient(dest net.Destination, streamSettings *internet.MemoryStreamConfig) (hyClient.Client, error) {
	tlsConfig, err := GetClientTLSConfig(streamSettings)
	if err != nil {
		return nil, err
	}

	serverAddr, err := InitAddress(dest)
	if err != nil {
		return nil, err
	}

	config := streamSettings.ProtocolSettings.(*Config)
	client, _, err := hyClient.NewClient(&hyClient.Config{
		TLSConfig:  *tlsConfig,
		Auth:       config.GetPassword(),
		ServerAddr: serverAddr,
		ConnFactory: &connFactory{
			NewFunc: func(addr net.Addr) (net.PacketConn, error) {
				rawConn, err := internet.ListenSystemPacket(context.Background(), &net.UDPAddr{
					IP:   []byte{0, 0, 0, 0},
					Port: 0,
				}, streamSettings.SocketSettings)
				if err != nil {
					return nil, err
				}
				return rawConn.(*net.UDPConn), nil
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CloseHyClient(dest net.Destination) error {
	ClientMutex.Lock()
	defer ClientMutex.Unlock()

	client, found := RunningClient[dest]
	if found {
		delete(RunningClient, dest)
		return client.Close()
	}
	return nil
}

func GetHyClient(dest net.Destination, streamSettings *internet.MemoryStreamConfig) (hyClient.Client, error) {
	ClientMutex.Lock()
	defer ClientMutex.Unlock()

	var client hyClient.Client
	var err error
	client, found := RunningClient[dest]
	if !found {
		// TODO: Clean idle connections
		client, err = NewHyClient(dest, streamSettings)
		if err != nil {
			return nil, err
		}
		RunningClient[dest] = client
	}
	return client, nil
}

func Dial(ctx context.Context, dest net.Destination, streamSettings *internet.MemoryStreamConfig) (internet.Connection, error) {
	config := streamSettings.ProtocolSettings.(*Config)

	client, err := GetHyClient(dest, streamSettings)
	if err != nil {
		CloseHyClient(dest)
		return nil, err
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
	frameSize := int(quicvarint.Len(hyProtocol.FrameTypeTCPRequest))
	buf := make([]byte, frameSize)
	hyProtocol.VarintPut(buf, hyProtocol.FrameTypeTCPRequest)
	conn.stream.Write(buf)
	return conn, nil
}

func init() {
	RunningClient = make(map[net.Destination]hyClient.Client)
	common.Must(internet.RegisterTransportDialer(protocolName, Dial))
}
