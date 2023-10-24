package popcount

import (
	"bytes"
	"context"
	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
	"io"
	"net"
)

//go:generate go run github.com/v2fly/v2ray-core/v5/common/errors/errorgen

type popCountConnection struct {
	config *Config
}

func (p popCountConnection) Client(conn net.Conn) net.Conn {
	return &popCountClientConnectionWrapper{
		Conn:   conn,
		config: p.config,
	}
}

func (p popCountConnection) Server(conn net.Conn) net.Conn {
	return &popCountServerConnectionWrapper{
		Conn:   conn,
		config: p.config,
	}
}

type popCountClientConnectionWrapper struct {
	net.Conn
	config  *Config
	written bool
}

func (p *popCountClientConnectionWrapper) Write(b []byte) (n int, err error) {
	if p.written {
		return p.Conn.Write(b)
	}
	if len(b) <= int(p.config.ExpansionPayloadLength) {
		return 0, newError("payload length shorter than expansion payload length")
	}
	var expansionPayload = make([]byte, p.config.ExpansionPayloadLength)
	copy(expansionPayload, b[:p.config.ExpansionPayloadLength])
	encodedData, err := popcountEncode(expansionPayload, p.config)
	if err != nil {
		return 0, newError("failed to popcount encode payload").Base(err)
	}
	resultData := append(encodedData, b[p.config.ExpansionPayloadLength:]...)
	_, err = p.Conn.Write(resultData)
	if err != nil {
		return 0, newError("failed to write encoded payload").Base(err)
	}
	p.written = true
	return len(b), nil
}

type popCountServerConnectionWrapper struct {
	net.Conn
	config *Config
	read   bool
	reader *bytes.Reader
}

func (p *popCountServerConnectionWrapper) Read(b []byte) (n int, err error) {
	if p.read {
		return p.Conn.Read(b)
	}

	if p.reader != nil {
		n, err = p.reader.Read(b)
		if err == nil {
			return n, nil
		}
		if err == io.EOF {
			p.read = true
			return p.Conn.Read(b)
		}
	}

	popCountLength := p.config.ExpansionPayloadLength + p.config.DiversifierLength + p.config.AddedLength
	var popCountPayload = make([]byte, popCountLength)
	_, err = p.Conn.Read(popCountPayload)
	if err != nil {
		return 0, newError("failed to read payload").Base(err)
	}
	decodedData, err := popcountDecode(popCountPayload, p.config)
	if err != nil {
		return 0, newError("failed to popcount decode payload").Base(err)
	}
	p.reader = bytes.NewReader(decodedData)

	n, err = p.reader.Read(b)
	return n, err
}

func newPopcountConnectionAuthenticator(config *Config) (internet.ConnectionAuthenticator, error) {
	return popCountConnection{
		config: config,
	}, nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return newPopcountConnectionAuthenticator(config.(*Config))
	}))
}
