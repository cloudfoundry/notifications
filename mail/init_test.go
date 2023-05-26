package mail_test

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"log"
	"net"
	"net/url"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMailSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mail")
}

const (
	certPEM = `-----BEGIN CERTIFICATE-----
MIID4DCCAsigAwIBAgIJANRk1GTj0oKXMA0GCSqGSIb3DQEBBQUAMH4xCzAJBgNV
BAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRUwEwYDVQQHDAxTYW50YSBNb25p
Y2ExFTATBgNVBAoMDFBpdm90YWwgTGFiczEWMBQGA1UECwwNQ2xvdWQgRm91bmRy
eTEUMBIGA1UEAwwLZGV2ZWxvcG1lbnQwHhcNMTQwNzA3MTgxODE1WhcNMjQwNzA0
MTgxODE1WjB+MQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEVMBMG
A1UEBwwMU2FudGEgTW9uaWNhMRUwEwYDVQQKDAxQaXZvdGFsIExhYnMxFjAUBgNV
BAsMDUNsb3VkIEZvdW5kcnkxFDASBgNVBAMMC2RldmVsb3BtZW50MIIBIjANBgkq
hkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyT+YaveHuISiyHBQhStt4HqXBpAK2arS
F15roOXSI2L6O8LlMaotI4j1OQNf19/6Q/OiBWGGBsVuf0ywde+Ikb9RMuq9FU15
hRXecfFyK1PU7MXMUzzY5tnh9Tcgau6t+O7OgMNGSulU3tt3xV4kQ7oKulv/aHW7
3tTGs6Iib4MlXzp8BM7llA5XmggTL2evwGgO/tKEyoWVOCAaIFfYkf4AGf47xfu6
B6qLpd3o6mdTL2xyZIrRvsCJk+/ToaCRs6ibaM9BiRXktIFRYGlg6fY2KWvxLOL6
YD9PacHvgCCF4ONPEIL69gWks01RpOU60c5MhndNRiRW9+JuRlLMdQIDAQABo2Ew
XzAPBgNVHREECDAGhwR/AAABMB0GA1UdDgQWBBRuNXjTwoXq++HZLTjF4nLhOg9w
2DAfBgNVHSMEGDAWgBRuNXjTwoXq++HZLTjF4nLhOg9w2DAMBgNVHRMEBTADAQH/
MA0GCSqGSIb3DQEBBQUAA4IBAQCv1sk2oJ55l9LfP6bQkR/nADHVZT5SSitAXpVF
PDhk7yMtrokP2SkgOgVLlgs3H/qxaowaqg6zeSPdnAhWM/n0r25zx2HYO1KLcHvF
vRYAb2skOoiYrHo6OOGHfhYj+c0ikgag/0CDy9EZ05/b5xPCMRoRzyp9t2gJpaqz
5TYIwMbDqs0E9pJT/ZjAQauwYUggxmUdhLUBnaKzzjGy7AOAldJVi/N1MMoNQInI
FyGrPuv4+T355ntQ274RGdytyYjMmvBANWX5+xzCJfoKlsfMxQNRgiwhNHjpqgZK
e+/BhTCOY1sHLFwd0eZ/4psN9/ytZRtcH3Y8waIwuQi3MlH5
-----END CERTIFICATE-----`
	keyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAyT+YaveHuISiyHBQhStt4HqXBpAK2arSF15roOXSI2L6O8Ll
MaotI4j1OQNf19/6Q/OiBWGGBsVuf0ywde+Ikb9RMuq9FU15hRXecfFyK1PU7MXM
UzzY5tnh9Tcgau6t+O7OgMNGSulU3tt3xV4kQ7oKulv/aHW73tTGs6Iib4MlXzp8
BM7llA5XmggTL2evwGgO/tKEyoWVOCAaIFfYkf4AGf47xfu6B6qLpd3o6mdTL2xy
ZIrRvsCJk+/ToaCRs6ibaM9BiRXktIFRYGlg6fY2KWvxLOL6YD9PacHvgCCF4ONP
EIL69gWks01RpOU60c5MhndNRiRW9+JuRlLMdQIDAQABAoIBAQDHU3zUXY0Inh5Y
5p1p+OTgVKtnLZ4Bj2Z9DOEPQPHMaMkuDdBSS5pfutQffw8b0tSfHx0XtUs5Q604
2q1gcjpTGSoEg2l6Qv0catejBaCt919KkHLa8sZmh+F8rfgm0XZwu56+/CqQIeEU
xk0vqBnFFuxvPpWPUiUdBKQ14V24EWvpkzIk1PcFXzLAimmfLvrKLnfPqZbXug0L
7kVSqLjNPWdltwWcVQ+uEgVv9TLRbEcxNIt2/DbSxplxjTlmX3O5hzzJMYcjrPpZ
2haoNdFyNb4y3G4F1nNrNzsTkQYfyOZBPmvMAJrx66vtEE0En9nYlq4xfjgjDr0j
zJMhMa31AoGBAO7SV+ORuEAQ53wQU7ccvn4F8FQBaeNbSpTiHmARW0iEIvHkvp0N
0aEdbo8aunFXNneEULP0I8tn0pBqeXjJIuh1epk3sqQg6xGodOsDK3ivIj6oaYtk
ifIVTcCL+dsCLXLMrdjAvbUWGtgS1SfZU4aRxehlH8UISaftnOE1NzdHAoGBANe5
YR2eLb0JPwRluFeiOhiEwlkrdai8vlk8EBHSUmku/sG7a8UuttqL1dGaUuedGfUk
0igc3WWGU0M7JNrrk2c6BuhqwQd57A3FlvEXon4kzQAWXGWfsU2KVq51HG+/8Cfp
7pIybTr6ysulVtNSh1NxPwO0wgWYnmarC1kKUjRjAoGAXB5Ydlgb6OJcV9d4YxY8
SCIETHLrJB5vizQZIVcwja0iSYnBGJVe+bV/ksVtixBn2vv3oSIXuHrIlpnrVvLG
e0HtUzJPvs1PvtTqnEfxubBcFi0h4Pmb1/vtrMqRSq/xVemrWQMnabUoD5ZcD+3d
MPgDjZuMAJUszBB0Rc4gCTsCgYBtCRkKJFpP8u10JonfWXLt06R795h32jaH2fDx
YRIgcg14FGgreSoZGpbPY6ZFxUVKf/rtJXHOD+/jynAdavbNNSoqrVK1ma1zZIyf
fWe3RJiNU8AN6YJvg92+PhlKboRPWFEqeex15C8+cWqKU2ttBI9qKyHqPDLMB+Yr
cikMqwKBgDLBKiVLqSeli1BrlURWGMyl7j+NgYNdih+M2Ra8dAuWt6BTQfjYW53u
EK8URF8KvO4+PR5pRJDCNx6+uOLoTsBE7KBYEiLzK9rTpEBQIv35h4hmF75SNiF/
gMirbaXT377nSX0oPon0P1iUgl5tUJNqYnYdA+qcpoeCvXuObAzm
-----END RSA PRIVATE KEY-----`
	StateUnknown   = "unknown"
	StateConnected = "connected"
	StateClosed    = "closed"
)

type SMTPServer struct {
	URL             url.URL
	CurrentDelivery Delivery
	Deliveries      []Delivery
	Listener        *net.TCPListener
	SupportsTLS     bool
	ConnectWait     time.Duration
	halt            chan bool
	ConnectionState string
	FailsHello      bool
}

type Delivery struct {
	Recipient string
	Sender    string
	Data      []string
	UsedTLS   bool
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
	listenerURL, err := url.Parse("//" + listener.Addr().String())
	if err != nil {
		panic(err)
	}

	server := SMTPServer{
		URL:             *listenerURL,
		Listener:        listener,
		ConnectionState: StateUnknown,
		halt:            make(chan bool),
	}
	server.Run()
	return &server
}

func (server *SMTPServer) Run() {
	go func() {
		for {
			select {
			case <-server.halt:
				return
			case connection := <-server.Accept():
				go server.Respond(connection)
			}
		}
	}()
}

func (server *SMTPServer) Accept() chan *net.TCPConn {
	connectionChan := make(chan *net.TCPConn)

	go func() {
		connection, err := server.Listener.AcceptTCP()
		if err != nil {
			return
		}
		connectionChan <- connection
	}()

	return connectionChan
}

func (server *SMTPServer) Close() {
	server.halt <- true
	server.Listener.Close()
}

func (server *SMTPServer) Respond(conn net.Conn) {
	<-time.After(server.ConnectWait)
	server.ConnectionState = StateConnected

	input := bufio.NewReader(conn)
	output := bufio.NewWriter(conn)
	server.Broadcast(output)

Loop:
	for {
		msg, _ := input.ReadString('\n')
		switch {
		case strings.Contains(msg, "EHLO"):
			server.RespondToEHLO(output)
		case strings.Contains(msg, "STARTTLS"):
			conn, input, output = server.RespondToStartTLS(conn, input, output)
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
	if server.FailsHello {
		output.WriteString("550-FOOBAR\n")
		output.WriteString("550 CRAZYBANANA\r\n")
		output.Flush()
		return
	}

	output.WriteString("250-localhost Hello\n")
	if server.SupportsTLS {
		output.WriteString("250-STARTTLS\n")
		output.WriteString("250 AUTH PLAIN LOGIN\r\n")
	} else {
		output.WriteString("250 AUTH LOGIN\r\n")
	}
	output.Flush()
}

func (server *SMTPServer) RespondToStartTLS(conn net.Conn, input *bufio.Reader, output *bufio.Writer) (*tls.Conn, *bufio.Reader, *bufio.Writer) {
	output.WriteString("220 Go ahead\r\n")
	output.Flush()

	server.CurrentDelivery.UsedTLS = true

	cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	tlsConn := tls.Server(conn, &config)

	return tlsConn, bufio.NewReader(tlsConn), bufio.NewWriter(tlsConn)
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
	server.ConnectionState = StateClosed
}
