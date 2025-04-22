package transport

import (
	"crypto/sha256"
	"net"

	"github.com/flynn/noise"
)

func WrapNoiseServer(c net.Conn, psk string) (net.Conn, error) {
	pskBytes := sha256.Sum256([]byte(psk))
	privKey := sha256.Sum256([]byte(psk))
	staticKey := noise.DHKey{Private: privKey[:]}

	hs := noise.HandshakeState{
		Config: noise.Config{
			CipherSuite:   noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashSHA256),
			Pattern:       noise.HandshakeXK,
			Initiator:     false,
			StaticKeypair: staticKey,
			PresharedKey:  pskBytes[:],
		},
	}
	return noise.Wrap(c, &hs)
}

func WrapNoiseClient(c net.Conn, psk string) (net.Conn, error) {
	pskBytes := sha256.Sum256([]byte(psk))

	hs := noise.HandshakeState{
		Config: noise.Config{
			CipherSuite:  noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashSHA256),
			Pattern:      noise.HandshakeXK,
			Initiator:    true,
			PresharedKey: pskBytes[:],
		},
	}
	return noise.Wrap(c, &hs)
}
