package hysteria2

import (
	"context"

	hy "github.com/apernet/hysteria/core/server"
	"github.com/apernet/quic-go"

	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/protocol/tls/cert"
	"github.com/v2fly/v2ray-core/v5/common/signal/done"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tls"
)

// Listener is an internet.Listener that listens for TCP connections.
type Listener struct {
	rawConn  *sysConn
	listener *quic.Listener
	done     *done.Instance
	addConn  internet.ConnHandler
	config   *Config
}

// Addr implements internet.Listener.Addr.
func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}

// Close implements internet.Listener.Close.
func (l *Listener) Close() error {
	return nil
}

func GetTLSConfig(streamSettings *internet.MemoryStreamConfig) *hy.TLSConfig {
	tlsSetting := tls.ConfigFromStreamSettings(streamSettings)
	if tlsSetting == nil {
		tlsSetting = &tls.Config{
			Certificate: []*tls.Certificate{
				tls.ParseCertificate(
					cert.MustGenerate(nil, cert.DNSNames(internalDomain), cert.CommonName(internalDomain)),
				),
			},
		}
	}
	return &hy.TLSConfig{Certificates: tlsSetting.GetTLSConfig().Certificates}
}

// Listen creates a new Listener based on configurations.
func Listen(ctx context.Context, address net.Address, port net.Port, streamSettings *internet.MemoryStreamConfig, handler internet.ConnHandler) (internet.Listener, error) {
	if address.Family().IsDomain() {
		return nil, newError("domain address is not allows for listening quic")
	}

	config := streamSettings.ProtocolSettings.(*Config)
	rawConn, err := internet.ListenSystemPacket(context.Background(), &net.UDPAddr{
		IP:   address.IP(),
		Port: int(port),
	}, streamSettings.SocketSettings)
	if err != nil {
		return nil, err
	}

	hyServer, err := hy.NewServer(&hy.Config{
		Conn:      rawConn,
		TLSConfig: *GetTLSConfig(streamSettings),
	})

	err = hyServer.Serve()
	if err != nil {
		rawConn.Close()
		return nil, err
	}

	listener := &Listener{
		done:    done.New(),
		addConn: handler,
		config:  config,
	}

	return listener, nil
}

func init() {
	common.Must(internet.RegisterTransportListener(protocolName, Listen))
}
