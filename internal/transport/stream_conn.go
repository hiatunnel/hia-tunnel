package transport

import (
	"net"
	"time"

	"github.com/quic-go/quic-go"
)

type streamConn struct {
	quic.Stream
	conn quic.Connection
}

func (s *streamConn) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *streamConn) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *streamConn) SetDeadline(t time.Time) error {
	if err := s.SetReadDeadline(t); err != nil {
		return err
	}
	return s.SetWriteDeadline(t)
}

func (s *streamConn) SetReadDeadline(t time.Time) error {
	return s.Stream.SetReadDeadline(t)
}

func (s *streamConn) SetWriteDeadline(t time.Time) error {
	return s.Stream.SetWriteDeadline(t)
}
