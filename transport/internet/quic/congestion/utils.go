package congestion

import (
	"github.com/apernet/quic-go"
	"github.com/v2fly/v2ray-core/v5/transport/internet/quic/congestion/bbr"
	"github.com/v2fly/v2ray-core/v5/transport/internet/quic/congestion/brutal"
)

func UseBBR(conn quic.Connection) {
	conn.SetCongestionControl(bbr.NewBbrSender(
		bbr.DefaultClock{},
		bbr.GetInitialPacketSize(conn.RemoteAddr()),
	))
}

func UseBrutal(conn quic.Connection, tx uint64) {
	conn.SetCongestionControl(brutal.NewBrutalSender(tx))
}
