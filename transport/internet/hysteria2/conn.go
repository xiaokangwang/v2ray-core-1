package hysteria2

import (
	"time"

	hy "github.com/apernet/hysteria/core/client"
	"github.com/apernet/quic-go"

	"github.com/v2fly/v2ray-core/v5/common/buf"
	"github.com/v2fly/v2ray-core/v5/common/net"
)

type HyConn struct {
	IsUDPExtension bool
	UDPSession     hy.HyUDPConn

	stream quic.Stream
	local  net.Addr
	remote net.Addr
}

func (c *HyConn) Read(b []byte) (int, error) {
	if c.IsUDPExtension {
		data, address, err := c.UDPSession.Receive()
		if err != nil {
			return 0, err
		}
		b = append([]byte(address), data...)
		return len(b), nil
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
		return len(b), c.UDPSession.Send(b, "")
	}
	return c.stream.Write(b)
}

func (c *HyConn) Close() error {
	if c.IsUDPExtension {
		return c.UDPSession.Close()
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
