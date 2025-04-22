package transport

import (
	"net"
	"time"

	"github.com/quic-go/quic-go"
)

type streamConn struct {
	quic.Stream
	localAddr  net.Addr
	remoteAddr net.Addr
}

func (c *streamConn) Read(b []byte) (int, error) {
	return c.Stream.Read(b)
}

func (c *streamConn) Write(b []byte) (int, error) {
	return c.Stream.Write(b)
}

func (c *streamConn) Close() error {
	return c.Stream.Close()
}

func (c *streamConn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *streamConn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *streamConn) SetDeadline(t time.Time) error {
	_ = c.SetReadDeadline(t)
	_ = c.SetWriteDeadline(t)
	return nil
}

func (c *streamConn) SetReadDeadline(t time.Time) error {
	return c.Stream.SetReadDeadline(t)
}

func (c *streamConn) SetWriteDeadline(t time.Time) error {
	return c.Stream.SetWriteDeadline(t)
}

func WrapStreamAsConn(stream quic.Stream) net.Conn {
	return &streamConn{
		Stream:     stream,
		localAddr:  &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: 0},
		remoteAddr: &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: 0},
	}
}
