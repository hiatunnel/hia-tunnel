package transport
import ("github.com/quic-go/quic-go";"time")
func QuicServer(addr string)(quic.Listener,error){
    return quic.ListenAddr(addr,ServerTLSConfig(),&quic.Config{EnableDatagrams:true,KeepAlive:true,MaxIdleTimeout:30*time.Second})
}
func QuicDial(addr,sni string)(quic.Connection,error){
    return quic.DialAddr(addr,ClientTLSConfig(sni),&quic.Config{EnableDatagrams:true,KeepAlive:true})
}
