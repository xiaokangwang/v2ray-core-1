package quic

import (
	"context"
	"sync"
	"time"

	"github.com/apernet/quic-go"

	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/task"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
	"github.com/v2fly/v2ray-core/v5/transport/internet/quic/congestion"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tls"
)

type connectionContext struct {
	rawConn *sysConn
	conn    quic.Connection
}

var errConnectionClosed = newError("connection closed")

func (c *connectionContext) openStream(destAddr net.Addr) (*interConn, error) {
	if !isActive(c.conn) {
		return nil, errConnectionClosed
	}

	stream, err := c.conn.OpenStream()
	if err != nil {
		return nil, err
	}

	conn := &interConn{
		stream: stream,
		local:  c.conn.LocalAddr(),
		remote: destAddr,
	}

	return conn, nil
}

type clientConnections struct {
	access             sync.Mutex
	runningConnections map[net.Destination][]*connectionContext
	cleanup            *task.Periodic
}

func isActive(s quic.Connection) bool {
	select {
	case <-s.Context().Done():
		return false
	default:
		return true
	}
}

func removeInactiveConnections(conns []*connectionContext) []*connectionContext {
	activeConnections := make([]*connectionContext, 0, len(conns))
	for _, s := range conns {
		if isActive(s.conn) {
			activeConnections = append(activeConnections, s)
			continue
		}
		if err := s.conn.CloseWithError(0, ""); err != nil {
			newError("failed to close connection").Base(err).WriteToLog()
		}
		if err := s.rawConn.Close(); err != nil {
			newError("failed to close raw connection").Base(err).WriteToLog()
		}
	}

	if len(activeConnections) < len(conns) {
		return activeConnections
	}

	return conns
}

func openStream(conns []*connectionContext, destAddr net.Addr) *interConn {
	for _, s := range conns {
		if !isActive(s.conn) {
			continue
		}

		conn, err := s.openStream(destAddr)
		if err != nil {
			continue
		}

		return conn
	}

	return nil
}

func (s *clientConnections) cleanConnections() error {
	s.access.Lock()
	defer s.access.Unlock()

	if len(s.runningConnections) == 0 {
		return nil
	}

	newConnMap := make(map[net.Destination][]*connectionContext)

	for dest, conns := range s.runningConnections {
		conns = removeInactiveConnections(conns)
		if len(conns) > 0 {
			newConnMap[dest] = conns
		}
	}

	s.runningConnections = newConnMap
	return nil
}

func (s *clientConnections) openConnection(destAddr net.Addr, config *Config, tlsConfig *tls.Config, sockopt *internet.SocketConfig) (internet.Connection, error) {
	s.access.Lock()
	defer s.access.Unlock()

	if s.runningConnections == nil {
		s.runningConnections = make(map[net.Destination][]*connectionContext)
	}

	dest := net.DestinationFromAddr(destAddr)

	var conns []*connectionContext
	if s, found := s.runningConnections[dest]; found {
		conns = s
		conn := openStream(conns, destAddr)
		if conn != nil {
			newError("dialing QUIC stream to ", dest).WriteToLog()
			return conn, nil
		}
	}
	conns = removeInactiveConnections(conns)
	newError("dialing QUIC new connection to ", dest).WriteToLog()

	rawConn, err := internet.ListenSystemPacket(context.Background(), &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 0,
	}, sockopt)
	if err != nil {
		return nil, err
	}

	quicConfig := InitQuicConfig()

	sysConn, err := wrapSysConn(rawConn.(*net.UDPConn), config)
	if err != nil {
		rawConn.Close()
		return nil, err
	}

	tr := quic.Transport{
		Conn: sysConn,
	}

	conn, err := tr.DialEarly(context.Background(), destAddr, tlsConfig.GetTLSConfig(tls.WithDestination(dest)), quicConfig)
	if err != nil {
		sysConn.Close()
		return nil, err
	}
	SetCongestion(conn, config)

	context := &connectionContext{
		conn:    conn,
		rawConn: sysConn,
	}
	s.runningConnections[dest] = append(conns, context)
	return context.openStream(destAddr)
}

func SetCongestion(conn quic.Connection, config *Config) {
	congestionType := config.Congestion.GetType()
	if congestionType == "brutal" && config.Congestion.GetSendMbps() != 0 {
		sendBps := config.Congestion.GetSendMbps() * 1000000
		congestion.UseBrutal(conn, sendBps)
	} else if congestionType == "bbr" {
		congestion.UseBBR(conn)
	}
	newError("[QUIC] using congestion: ", congestionType).WriteToLog()
}

const (
	defaultStreamReceiveWindow  = 8388608                            // 8MB
	defaultConnReceiveWindow    = defaultStreamReceiveWindow * 5 / 2 // 20MB
	defaultMaxIdleTimeout       = 30 * time.Second
	defaultKeepAlivePeriod      = 10 * time.Second
	defaultHandshakeIdleTimeout = 10 * time.Second
)

func InitQuicConfig() *quic.Config {
	quicConfig := &quic.Config{
		HandshakeIdleTimeout:           defaultHandshakeIdleTimeout,
		MaxIdleTimeout:                 defaultMaxIdleTimeout,
		KeepAlivePeriod:                defaultKeepAlivePeriod,
		InitialStreamReceiveWindow:     defaultStreamReceiveWindow,
		InitialConnectionReceiveWindow: defaultConnReceiveWindow,
	}
	return quicConfig
}

var client clientConnections

func init() {
	client.runningConnections = make(map[net.Destination][]*connectionContext)
	client.cleanup = &task.Periodic{
		Interval: time.Minute,
		Execute:  client.cleanConnections,
	}
	common.Must(client.cleanup.Start())
}

func Dial(ctx context.Context, dest net.Destination, streamSettings *internet.MemoryStreamConfig) (internet.Connection, error) {
	tlsConfig := tls.ConfigFromStreamSettings(streamSettings)
	if tlsConfig == nil {
		tlsConfig = &tls.Config{
			ServerName:    internalDomain,
			AllowInsecure: true,
		}
	}

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

	config := streamSettings.ProtocolSettings.(*Config)

	return client.openConnection(destAddr, config, tlsConfig, streamSettings.SocketSettings)
}

func init() {
	common.Must(internet.RegisterTransportDialer(protocolName, Dial))
}
