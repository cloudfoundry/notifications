package servers

import (
    "net"
    "os"

    "bitbucket.org/chrj/smtpd"
)

type SMTPServer struct {
    server     *smtpd.Server
    Deliveries []smtpd.Envelope
}

func NewSMTPServer() *SMTPServer {
    return &SMTPServer{
        server: &smtpd.Server{
            Addr: "127.0.0.1:0",
        },
        Deliveries: make([]smtpd.Envelope, 0),
    }
}

func (s *SMTPServer) Boot() {
    listener, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        panic(err)
    }

    s.server.Handler = s.Handler

    go s.server.Serve(listener)

    addr := listener.Addr().String()
    host, port, err := net.SplitHostPort(addr)
    if err != nil {
        panic(err)
    }
    os.Setenv("SMTP_HOST", host)
    os.Setenv("SMTP_PORT", port)
}

func (s *SMTPServer) Handler(peer smtpd.Peer, env smtpd.Envelope) error {
    s.Deliveries = append(s.Deliveries, env)
    return nil
}
