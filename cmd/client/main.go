package main
import("flag";"io";"log";"net";
    "github.com/armon/go-socks5";
    "github.com/hiatunnel/hia-tunnel/internal/config";
    "github.com/hiatunnel/hia-tunnel/internal/mux";
    "github.com/hiatunnel/hia-tunnel/internal/transport")
func main(){
    cfg:=flag.String("c","./client.json","config");flag.Parse()
    var cc config.ClientConf
    if err:=config.Load(*cfg,&cc);err!=nil{log.Fatal(err)}
    for _,p:=range cc.Peers{go dial(p)}
    select{}
}
func dial(p config.Peer){
    qc,err:=transport.QuicDial(p.Server,"cdn.cloudflare.com");if err!=nil{log.Println(err);return}
    sec,_:=transport.WrapNoiseClient(qc.Stream(0),p.PSK)
    sess,_:=mux.Client(sec,p.MaxStreams)
    dialf:=func(_,addr string)(net.Conn,error){
        st,err:=sess.OpenStream();if err!=nil{return nil,err}
        st.Write([]byte(addr+"\n"));return st,nil}
    s,_:=socks5.New(&socks5.Config{Dial:dialf})
    go s.ListenAndServe("tcp",p.SocksLocal)
    log.Printf("[%s] socks %s -> %s",p.Name,p.SocksLocal,p.Server)
}
