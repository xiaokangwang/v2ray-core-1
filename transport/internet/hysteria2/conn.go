package hysteria2

import (
	"time"

	hy_client "github.com/apernet/hysteria/core/client"
	hy_server "github.com/apernet/hysteria/core/server"
	"github.com/apernet/quic-go"

	"github.com/v2fly/v2ray-core/v5/common/buf"
	"github.com/v2fly/v2ray-core/v5/common/net"
)

const CanNotUseUdpExtension = "Only hysteria2 proxy protocol can use udpExtension."

type HyConn struct {
	IsUDPExtension   bool
	IsServer         bool
	ClientUDPSession hy_client.HyUDPConn
	ServerUDPSession hy_server.UDPConn
	Target           net.Destination

	stream quic.Stream
	local  net.Addr
	remote net.Addr
}

func (c *HyConn) Read(b []byte) (int, error) {
	if c.IsUDPExtension {
		return 0, newError(CanNotUseUdpExtension)
	}
	return c.stream.Read(b)
}

func (c *HyConn) WriteMultiBuffer(mb buf.MultiBuffer) error {
	mb = buf.Compact(mb)
	mb, err := buf.WriteMultiBuffer(c, mb)
	buf.ReleaseMulti(mb)
	return err
}

func (c *HyConn) Write(b []byte) (int, error) {
	if c.IsUDPExtension {
		return 0, newError(CanNotUseUdpExtension)
	}
	return c.stream.Write(b)
}

func (c *HyConn) WritePacket(b []byte, dest net.Destination) (int, error) {
	if !c.IsUDPExtension || c.ClientUDPSession == nil {
		return 0, newError(CanNotUseUdpExtension)
	}

	if c.IsServer {
		return c.ServerUDPSession.WriteTo(b, dest.NetAddr())
	}
	return len(b), c.ClientUDPSession.Send(b, dest.NetAddr())
}

func (c *HyConn) ReadPacket(b []byte) (int, *net.Destination, error) {
	if !c.IsUDPExtension || c.ClientUDPSession == nil {
		return 0, nil, newError(CanNotUseUdpExtension)
	}

	address := ""
	var err error
	if c.IsServer {
		_, address, err = c.ServerUDPSession.ReadFrom(b)
	} else {
		b, address, err = c.ClientUDPSession.Receive()
	}
	if err != nil {
		return 0, nil, err
	}
	dest, err := net.ParseDestination(address)
	dest.Network = net.Network_UDP
	if err != nil {
		return 0, nil, err
	}
	return len(b), &dest, nil
}

func (c *HyConn) Close() error {
	if c.IsUDPExtension {
		if c.ClientUDPSession == nil {
			return newError(CanNotUseUdpExtension)
		}
		return c.ClientUDPSession.Close()
	}
	return c.stream.Close()
}

func (c *HyConn) LocalAddr() net.Addr {
	return c.local
}

func (c *HyConn) RemoteAddr() net.Addr {
	return c.remote
}

func (c *HyConn) SetDeadline(t time.Time) error {
	if c.IsUDPExtension {
		return nil
	}
	return c.stream.SetDeadline(t)
}

func (c *HyConn) SetReadDeadline(t time.Time) error {
	if c.IsUDPExtension {
		return nil
	}
	return c.stream.SetReadDeadline(t)
}

func (c *HyConn) SetWriteDeadline(t time.Time) error {
	if c.IsUDPExtension {
		return nil
	}
	return c.stream.SetWriteDeadline(t)
}
