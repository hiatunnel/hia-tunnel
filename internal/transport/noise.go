package transport

import (
	"crypto/sha256"
	"net"

	"github.com/flynn/noise"
)

type noiseConn struct {
	net.Conn
	cs *noise.CipherState
}

func (c *noiseConn) Read(b []byte) (int, error) {
	buf := make([]byte, len(b)+1024)
	n, err := c.Conn.Read(buf)
	if err != nil {
		return 0, err
	}
	out, err := c.cs.Decrypt(nil, nil, buf[:n])
	if err != nil {
		return 0, err
	}
	copy(b, out)
	return len(out), nil
}

func (c *noiseConn) Write(b []byte) (int, error) {
	out, err := c.cs.Encrypt(nil, nil, b)
	if err != nil {
		return 0, err
	}
	_, err = c.Conn.Write(out)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func WrapNoiseServer(c net.Conn, psk string) (net.Conn, error) {
	pskBytes := sha256.Sum256([]byte(psk))
	privKey := sha256.Sum256([]byte(psk))

	config := noise.Config{
		CipherSuite:   noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashSHA256),
		Pattern:       noise.HandshakeXK,
		Initiator:     false,
		StaticKeypair: noise.DHKey{Private: privKey[:]},
		PresharedKey:  pskBytes[:],
	}

	hs, err := noise.NewHandshakeState(config)
	if err != nil {
		return nil, err
	}

	// 接收客户端握手
	msg := make([]byte, 512)
	n, err := c.Read(msg)
	if err != nil {
		return nil, err
	}
	_, _, csR, err := hs.ReadMessage(nil, msg[:n])
	if err != nil {
		return nil, err
	}

	// 发送服务端响应
	out, _, _, err := hs.WriteMessage(nil, nil)
	if err != nil {
		return nil, err
	}
	if _, err := c.Write(out); err != nil {
		return nil, err
	}

	return &noiseConn{Conn: c, cs: csR}, nil
}

func WrapNoiseClient(c net.Conn, psk string) (net.Conn, error) {
	pskBytes := sha256.Sum256([]byte(psk))

	config := noise.Config{
		CipherSuite:  noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashSHA256),
		Pattern:      noise.HandshakeXK,
		Initiator:    true,
		PresharedKey: pskBytes[:],
	}

	hs, err := noise.NewHandshakeState(config)
	if err != nil {
		return nil, err
	}

	// 发送客户端握手
	msg, csW, _, err := hs.WriteMessage(nil, nil)
	if err != nil {
		return nil, err
	}
	if _, err := c.Write(msg); err != nil {
		return nil, err
	}

	// 接收服务端响应
	resp := make([]byte, 512)
	n, err := c.Read(resp)
	if err != nil {
		return nil, err
	}
	_, _, _, err = hs.ReadMessage(nil, resp[:n])
	if err != nil {
		return nil, err
	}

	return &noiseConn{Conn: c, cs: csW}, nil
}
