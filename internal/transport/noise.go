package transport

import (
	"crypto/sha256"
	"net"

	"github.com/flynn/noise"
)

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

	msg := make([]byte, 512)
	n, err := c.Read(msg)
	if err != nil {
		return nil, err
	}
	_, _, _, err = hs.ReadMessage(nil, msg[:n])
	if err != nil {
		return nil, err
	}

	out, _, _, err := hs.WriteMessage(nil, nil)
	if err != nil {
		return nil, err
	}
	if _, err := c.Write(out); err != nil {
		return nil, err
	}

	return NewNoiseConn(c, hs), nil
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

	msg, _, _, err := hs.WriteMessage(nil, nil)
	if err != nil {
		return nil, err
	}
	if _, err := c.Write(msg); err != nil {
		return nil, err
	}

	resp := make([]byte, 512)
	n, err := c.Read(resp)
	if err != nil {
		return nil, err
	}
	_, _, _, err = hs.ReadMessage(nil, resp[:n])
	if err != nil {
		return nil, err
	}

	return NewNoiseConn(c, hs), nil
}
