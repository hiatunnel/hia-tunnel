package transport

import (
	"errors"
	"io"
	"net"
	"sync"

	"github.com/flynn/noise"
)

type NoiseConn struct {
	conn net.Conn
	hs   *noise.HandshakeState
	rwMu sync.Mutex
}

func NewNoiseConn(c net.Conn, hs *noise.HandshakeState) *NoiseConn {
	return &NoiseConn{
		conn: c,
		hs:   hs,
	}
}

func (nc *NoiseConn) Read(p []byte) (int, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(nc.conn, header); err != nil {
		return 0, err
	}

	size := int(header[0])<<8 | int(header[1])
	buf := make([]byte, size)
	if _, err := io.ReadFull(nc.conn, buf); err != nil {
		return 0, err
	}

	nc.rwMu.Lock()
	defer nc.rwMu.Unlock()
	msg, _, _, err := nc.hs.ReadMessage(nil, buf)
	if err != nil {
		return 0, err
	}

	copy(p, msg)
	return len(msg), nil
}

func (nc *NoiseConn) Write(p []byte) (int, error) {
	nc.rwMu.Lock()
	defer nc.rwMu.Unlock()
	cipher, _, _, err := nc.hs.WriteMessage(nil, p)
	if err != nil {
		return 0, err
	}

	size := len(cipher)
	if size > 65535 {
		return 0, errors.New("message too long")
	}

	header := []byte{byte(size >> 8), byte(size & 0xff)}
	if _, err := nc.conn.Write(header); err != nil {
		return 0, err
	}

	if _, err := nc.conn.Write(cipher); err != nil {
		return 0, err
	}

	return len(p), nil
}

func (nc *NoiseConn) Close() error {
	return nc.conn.Close()
}

func (nc *NoiseConn) LocalAddr() net.Addr {
	return nc.conn.LocalAddr()
}

func (nc *NoiseConn) RemoteAddr() net.Addr {
	return nc.conn.RemoteAddr()
}

func (nc *NoiseConn) SetDeadline(t time.Time) error {
	return nc.conn.SetDeadline(t)
}

func (nc *NoiseConn) SetReadDeadline(t time.Time) error {
	return nc.conn.SetReadDeadline(t)
}

func (nc *NoiseConn) SetWriteDeadline(t time.Time) error {
	return nc.conn.SetWriteDeadline(t)
}
