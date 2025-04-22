module github.com/hiatunnel/hia-tunnel

go 1.22

require (
	// SOCKS5 dial helper
	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5

	// Noise double‑encryption
	github.com/flynn/noise v1.0.0

	// QUIC transport
	github.com/quic-go/quic-go v0.39.2

	// TLS fingerprint randomiser
	github.com/refraction-networking/utls v1.6.6

	// Stream multiplexer
	github.com/xtaci/smux v1.5.34
)
