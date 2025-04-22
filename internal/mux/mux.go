package mux

import (
	"github.com/xtaci/smux"
	"io"
	"net"
)

func Server(c io.ReadWriteCloser, max int) (*smux.Session, error) {
	cfg := smux.DefaultConfig()
	return smux.Server(c, cfg)
}

func Client(c io.ReadWriteCloser, max int) (*smux.Session, error) {
	cfg := smux.DefaultConfig()
	return smux.Client(c, cfg)
}

func Pipe(s *smux.Stream, target string) {
	dst, err := net.Dial("tcp", target)
	if err != nil {
		_ = s.Close()
		return
	}
	go io.Copy(dst, s)
	io.Copy(s, dst)
}
