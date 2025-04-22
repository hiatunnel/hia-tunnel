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

func handle(qc quic.Connection, sc *config.ServerConf) {
    // ✅ 打开第一个流
    stream, err := qc.AcceptStream(context.Background())
    if err != nil {
        log.Println("accept stream:", err)
        return
    }

    // ✅ 加密流
    sec, err := transport.WrapNoiseServer(stream, sc.PSK)
    if err != nil {
        log.Println("noise handshake:", err)
        return
    }

    sess, err := mux.Server(sec, 1024)
    if err != nil {
        log.Println("mux:", err)
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
