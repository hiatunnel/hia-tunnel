package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"

	"github.com/armon/go-socks5"
	"github.com/hiatunnel/hia-tunnel/internal/config"
	"github.com/hiatunnel/hia-tunnel/internal/mux"
	"github.com/hiatunnel/hia-tunnel/internal/transport"
)

func main() {
	cfg := flag.String("c", "./client.json", "config")
	flag.Parse()

	var cc config.ClientConf
	if err := config.Load(*cfg, &cc); err != nil {
		log.Fatal(err)
	}

	for _, p := range cc.Peers {
		go dial(p)
	}

	select {}
}

func dial(p config.Peer) {
	qc, err := transport.QuicDial(p.Server, "cdn.cloudflare.com")
	if err != nil {
		log.Printf("[%s] quic dial error: %v", p.Name, err)
		return
	}

	stream, err := qc.OpenStreamSync(context.Background())
	if err != nil {
		log.Printf("[%s] open stream error: %v", p.Name, err)
		return
	}

	conn := transport.WrapStreamAsConn(stream)
	sec, err := transport.WrapNoiseClient(conn, p.PSK)
	
	if err != nil {
		log.Printf("[%s] noise handshake error: %v", p.Name, err)
		return
	}

	sess, err := mux.Client(sec, p.MaxStreams)
	if err != nil {
		log.Printf("[%s] mux error: %v", p.Name, err)
		return
	}

	dialf := func(_, addr string) (net.Conn, error) {
		st, err := sess.OpenStream()
		if err != nil {
			return nil, err
		}
		_, err = st.Write([]byte(addr + "\n"))
		if err != nil {
			return nil, err
		}
		return st, nil
	}

	s, err := socks5.New(&socks5.Config{Dial: dialf})
	if err != nil {
		log.Printf("[%s] socks5 init error: %v", p.Name, err)
		return
	}

	go func() {
		if err := s.ListenAndServe("tcp", p.SocksLocal); err != nil {
			log.Printf("[%s] socks5 error: %v", p.Name, err)
		}
	}()

	log.Printf("[%s] socks %s -> %s", p.Name, p.SocksLocal, p.Server)
}
