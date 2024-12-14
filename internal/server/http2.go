package server

import (
    "crypto/tls"
    "golang.org/x/net/http2"
    "net/http"
    "time"
)

type HTTP2Server struct {
    server *http.Server
}

func NewHTTP2Server(addr string, handler http.Handler) *HTTP2Server {
    server := &http.Server{
        Addr:    addr,
        Handler: handler,
        TLSConfig: &tls.Config{
            MinVersion:               tls.VersionTLS12,
            CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
            PreferServerCipherSuites: true,
            CipherSuites: []uint16{
                tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
                tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
            },
        },
    }

    http2.ConfigureServer(server, &http2.Server{
        MaxConcurrentStreams: 250,
        IdleTimeout:         30 * time.Second,
    })

    return &HTTP2Server{server: server}
} 