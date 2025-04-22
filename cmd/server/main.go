package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"

	"github.com/hiatunnel/hia-tunnel/internal/config"
	"github.com/hiatunnel/hia-tunnel/internal/mux"
	"github.com/hiatunnel/hia-tunnel/internal/transport"
	"github.com/quic-go/quic-go"
)

func main() {
	cfgFile := flag.String("c", "/etc/swift-tunnel/server.json", "config")
	flag.Parse()

	var sc config.ServerConf
	if err := config.Load(*cfgFile, &sc); err != nil {
		log.Fatal(err)
	}

	go forwards(sc.Forwards)

	lis, err := transport.QuicServer(sc.Listen)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[server] quic %s", sc.Listen)

	for {
		qc, err := lis.Accept(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		go handle(qc, &sc)
	}
}

func handle(q quic.Connection, sc *config.ServerConf) {
	stream, err := q.AcceptStream(context.Background())
	if err != nil {
		log.Println("stream error:", err)
		return
	}

	sec, err := transport.WrapNoiseServer(stream, sc.PSK)
	if err != nil {
		log.Println("noise handshake failed:", err)
		return
	}

	sess, err := mux.Server(sec, 1024)
	if err != nil {
		log.Println("mux server error:", err)
		return
	}

	for {
		st, err := sess.AcceptStream()
		if err != nil {
			return
		}
		buf := make([]byte, 256)
		n, _ := st.Read(buf)
		go mux.Pipe(st, string(buf[:n-1]))
	}
}

func forwards(fs []config.Forward) {
	for _, f := range fs {
		go func(f config.Forward) {
			l, err := net.Listen("tcp", f.Listen)
			if err != nil {
				log.Printf("listen error: %s -> %s: %v", f.Listen, f.Target, err)
				return
			}
			log.Printf("[fwd] %s -> %s", f.Listen, f.Target)
			for {
				c, err := l.Accept()
				if err != nil {
					continue
				}
				go proxy(c, f.Target)
			}
		}(f)
	}
}

func proxy(src net.Conn, target string) {
	dst, err := net.Dial("tcp", target)
	if err != nil {
		src.Close()
		return
	}
	go io.Copy(dst, src)
	io.Copy(src, dst)
}
