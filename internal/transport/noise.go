package transport
import ("crypto/sha256";"net";"github.com/flynn/noise")
func WrapNoiseServer(c net.Conn,psk string)(net.Conn,error){
    hs:=noise.HandshakeState{Config:noise.Config{
        CipherSuite:noise.NewCipherSuite(noise.DH25519,noise.CipherChaChaPoly,noise.HashSHA256),
        Pattern:noise.HandshakeXK,Initiator:false,
        StaticKeypair:noise.DHKey{Private:sha256.Sum256([]byte(psk))[:]},
        PresharedKey:sha256.Sum256([]byte(psk))[:},
    }
    return noise.Wrap(c,&hs)
}
func WrapNoiseClient(c net.Conn,psk string)(net.Conn,error){
    hs:=noise.HandshakeState{Config:noise.Config{
        CipherSuite:noise.NewCipherSuite(noise.DH25519,noise.CipherChaChaPoly,noise.HashSHA256),
        Pattern:noise.HandshakeXK,Initiator:true,
        PresharedKey:sha256.Sum256([]byte(psk))[:},
    }
    return noise.Wrap(c,&hs)
}
