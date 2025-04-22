package main
import("flag";"io";"log";"net";
    "github.com/hiatunnel/hia-tunnel/internal/config";
    "github.com/hiatunnel/hia-tunnel/internal/mux";
    "github.com/hiatunnel/hia-tunnel/internal/transport")
func main(){
    cfgFile:=flag.String("c","/etc/swift-tunnel/server.json","config");flag.Parse()
    var sc config.ServerConf
    if err:=config.Load(*cfgFile,&sc);err!=nil{log.Fatal(err)}
    go forwards(sc.Forwards)
    lis,err:=transport.QuicServer(sc.Listen);if err!=nil{log.Fatal(err)}
    log.Printf("[server] quic %s",sc.Listen)
    for{qc,err:=lis.Accept();if err!=nil{log.Println(err);continue}
        go handle(qc,&sc)}
}
func handle(q transport.QuicConn,sc *config.ServerConf){
    sec,_:=transport.WrapNoiseServer(q.Stream(0),sc.PSK)
    sess,_:=mux.Server(sec,1024)
    for{st,err:=sess.AcceptStream();if err!=nil{return}
        buf:=make([]byte,256);n,_:=st.Read(buf)
        go mux.Pipe(st,string(buf[:n-1]))}
}
func forwards(fs []config.Forward){
    for _,f:=range fs{go func(f config.Forward){
        l,err:=net.Listen("tcp",f.Listen);if err!=nil{return}
        log.Printf("[fwd] %s -> %s",f.Listen,f.Target)
        for{c,_:=l.Accept();go proxy(c,f.Target)}
    }(f)}
}
func proxy(src net.Conn,target string){
    dst,err:=net.Dial("tcp",target);if err!=nil{src.Close();return}
    go io.Copy(dst,src);io.Copy(src,dst)
}
