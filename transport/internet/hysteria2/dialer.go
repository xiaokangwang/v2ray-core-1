package hysteria2

import (
	"context"

	hy "github.com/apernet/hysteria/core/client"

	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tls"
)

var RunningClient map[net.Destination]hy.Client

var errConnectionClosed = newError("connection closed")

func InitTLSConifg(streamSettings *internet.MemoryStreamConfig) (*hy.TLSConfig, error) {
	tlsConfig := tls.ConfigFromStreamSettings(streamSettings)
	if tlsConfig == nil {
		tlsConfig = &tls.Config{
			ServerName:    internalDomain,
			AllowInsecure: true,
		}
	}
	res := &hy.TLSConfig{ServerName: tlsConfig.ServerName, InsecureSkipVerify: tlsConfig.AllowInsecure}
	return res, nil
}

func Dial(ctx context.Context, dest net.Destination, streamSettings *internet.MemoryStreamConfig) (internet.Connection, error) {
	var client hy.Client
	client, found := RunningClient[dest]
	if !found {
		var err error
		tlsConfig, err := InitTLSConifg(streamSettings)
		if err != nil {
			return nil, err
		}
		client, err = hy.NewClient(&hy.Config{TLSConfig: *tlsConfig})
		if err != nil {
			return nil, err
		}
		RunningClient[dest] = client
	}

	stream, err := client.OpenStream()
	if err != nil {
		return nil, err
	}

	quicConn := client.GetQuicConn()
	conn := &interConn{
		stream: stream,
		local:  quicConn.LocalAddr(),
		remote: quicConn.RemoteAddr(),
	}
	return conn, nil
}

func init() {
	common.Must(internet.RegisterTransportDialer(protocolName, Dial))
}
