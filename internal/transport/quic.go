package transport

import (
	"context"
	"time"

	"github.com/quic-go/quic-go"
)

func QuicServer(addr string) (quic.Listener, error) {
	return quic.ListenAddr(addr, ServerTLSConfig(), &quic.Config{
		EnableDatagrams: true,
		MaxIdleTimeout:  30 * time.Second,
	})
}

func QuicDial(addr, sni string) (quic.Connection, error) {
	return quic.DialAddr(
		context.Background(),
		addr,
		ClientTLSConfig(sni),
		&quic.Config{
			EnableDatagrams: true,
		},
	)
}
