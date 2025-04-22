package transport

import (
    "crypto/tls"
    "crypto/x509"
    "embed"
)

//go:embed devcert.pem devkey.pem
var devFS embed.FS

func ServerTLSConfig() *tls.Config {
    certPEM,_:=devFS.ReadFile("devcert.pem")
    keyPEM,_:=devFS.ReadFile("devkey.pem")
    cert,_:=tls.X509KeyPair(certPEM,keyPEM)
    return &tls.Config{
        Certificates:[]tls.Certificate{cert},
        NextProtos:[]string{"h3"},
        MinVersion:tls.VersionTLS13,
    }
}
func ClientTLSConfig(sni string)*tls.Config{
    return &tls.Config{
        InsecureSkipVerify:true,
        NextProtos:[]string{"h3"},
        ServerName:sni,
        MinVersion:tls.VersionTLS13,
        RootCAs:x509.NewCertPool(),
    }
}
