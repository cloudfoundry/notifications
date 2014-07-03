package mail_test

import (
    "bufio"
    "net"
    "net/url"
    "strings"
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestMailSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Mail Suite")
}

type SMTPServer struct {
    URL             url.URL
    CurrentDelivery Delivery
    Deliveries      []Delivery
    Listener        *net.TCPListener
}

type Delivery struct {
    Recipient string
    Sender    string
    Data      []string
}

func NewSMTPServer(user, pass string) *SMTPServer {
    addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
    if err != nil {
        panic(err)
    }
    listener, err := net.ListenTCP("tcp", addr)
    if err != nil {
        panic(err)
    }
    listenerURL, err := url.Parse(listener.Addr().String())
    if err != nil {
        panic(err)
    }

    server := SMTPServer{
        URL:      *listenerURL,
        Listener: listener,
    }
    server.Run()
    return &server
}

func (server *SMTPServer) Run() {
    go func() {
        for {
            connection, err := server.Listener.AcceptTCP()
            if err != nil {
                panic(err)
            }

            server.Respond(connection)
        }
    }()
}

func (server *SMTPServer) Respond(conn *net.TCPConn) {
    input := bufio.NewReader(conn)
    output := bufio.NewWriter(conn)
    server.Broadcast(output)

Loop:
    for {
        msg, _ := input.ReadString('\n')
        switch {
        case strings.Contains(msg, "EHLO"):
            server.RespondToEHLO(output)
        case strings.Contains(msg, "AUTH PLAIN"):
            server.RespondToAuthPlain(output)
        case strings.Contains(msg, "MAIL FROM"):
            server.RespondToMailFrom(output, msg)
        case strings.Contains(msg, "RCPT TO"):
            server.RespondToRcptTo(output, msg)
        case strings.Contains(msg, "DATA"):
            server.RespondToData(output)
            server.RecordData(output, input)
        case strings.Contains(msg, "QUIT"):
            server.RespondToQuit(output)
            break Loop
        }
    }
    server.Deliveries = append(server.Deliveries, server.CurrentDelivery)
    server.CurrentDelivery = Delivery{}
}

func (server *SMTPServer) Broadcast(output *bufio.Writer) {
    output.WriteString("220 localhost\r\n")
    output.Flush()
}

func (server *SMTPServer) RespondToEHLO(output *bufio.Writer) {
    output.WriteString("250-localhost Hello\n")
    output.WriteString("250 AUTH PLAIN LOGIN\r\n")
    output.Flush()
}

func (server *SMTPServer) RespondToAuthPlain(output *bufio.Writer) {
    output.WriteString("235 OK, Go ahead\r\n")
    output.Flush()
}

func (server *SMTPServer) RespondToMailFrom(output *bufio.Writer, msg string) {
    sender := strings.TrimSpace(msg)
    sender = strings.TrimPrefix(sender, "MAIL FROM:")
    sender = strings.Trim(sender, "<>")
    server.CurrentDelivery.Sender = sender

    output.WriteString("250 OK\r\n")
    output.Flush()
}

func (server *SMTPServer) RespondToRcptTo(output *bufio.Writer, msg string) {
    recipient := strings.TrimSpace(msg)
    recipient = strings.TrimPrefix(recipient, "RCPT TO:")
    recipient = strings.Trim(recipient, "<>")
    server.CurrentDelivery.Recipient = recipient

    output.WriteString("250 OK\r\n")
    output.Flush()
}

func (server *SMTPServer) RespondToData(output *bufio.Writer) {
    output.WriteString("354 OK\r\n")
    output.Flush()
}

func (server *SMTPServer) RecordData(output *bufio.Writer, input *bufio.Reader) {
    for {
        msg, err := input.ReadString('\n')
        if err != nil {
            panic(err)
        }

        if strings.TrimSpace(msg) == "." {
            break
        }
        server.CurrentDelivery.Data = append(server.CurrentDelivery.Data, strings.TrimSpace(msg))
    }
    output.WriteString("250 Written safely to disk.\r\n")
    output.Flush()
}

func (server *SMTPServer) RespondToQuit(output *bufio.Writer) {
    output.WriteString("221 BYE\r\n")
    output.Flush()
}
