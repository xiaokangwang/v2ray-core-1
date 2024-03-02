package hysteria2

import (
	"time"

	"github.com/apernet/quic-go"

	"github.com/v2fly/v2ray-core/v5/common/buf"
	"github.com/v2fly/v2ray-core/v5/common/net"
)

type TCPConn struct {
	stream quic.Stream
	local  net.Addr
	remote net.Addr
}

const (
	FrameTypeTCPRequest = 0x401
)

func (c *TCPConn) Read(b []byte) (int, error) {
	return c.stream.Read(b)
}

func (c *TCPConn) WriteMultiBuffer(mb buf.MultiBuffer) error {
	mb = buf.Compact(mb)
	mb, err := buf.WriteMultiBuffer(c, mb)
	buf.ReleaseMulti(mb)
	return err
}

func (c *TCPConn) Write(b []byte) (int, error) {
	return c.stream.Write(b)
}

func (c *TCPConn) Close() error {
	return c.stream.Close()
}

func (c *TCPConn) LocalAddr() net.Addr {
	return c.local
}

func (c *TCPConn) RemoteAddr() net.Addr {
	return c.remote
}

func (c *TCPConn) SetDeadline(t time.Time) error {
	return c.stream.SetDeadline(t)
}

func (c *TCPConn) SetReadDeadline(t time.Time) error {
	return c.stream.SetReadDeadline(t)
}

func (c *TCPConn) SetWriteDeadline(t time.Time) error {
	return c.stream.SetWriteDeadline(t)
}
