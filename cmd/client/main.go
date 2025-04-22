package main

import (
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
		log.Println("quic dial:", err)
		return
	}

	stream, err := qc.OpenStream()
	if err != nil {
		log.Println("open stream:", err)
		return
	}

	// ✅ Noise 加密通道
	sec, err := transport.WrapNoiseClient(stream, p.PSK)
	if err != nil {
		log.Println("noise handshake:", err)
		return
	}

	sess, err := mux.Client(sec, p.MaxStreams)
	if err != nil {
		log.Println("smux client:", err)
		return
	}

	dialFn := func(_, addr string) (net.Conn, error) {
		st, err := sess.OpenStream()
		if err != nil {
			return nil, err
		}
		st.Write([]byte(addr + "\n"))
		return st, nil
	}

	socks, _ := socks5.New(&socks5.Config{Dial: dialFn})
	go socks.ListenAndServe("tcp", p.SocksLocal)

	log.Printf("[%s] socks %s -> %s", p.Name, p.SocksLocal, p.Server)
}
