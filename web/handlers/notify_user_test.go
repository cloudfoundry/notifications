package handlers_test

import (
    "bufio"
    "bytes"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "log"
    "net"
    "net/http"
    "net/http/httptest"
    "os"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

type Envelope struct {
    AuthUser string
    AuthPass string
    From     string
    To       string
    Data     string
}

func (envelope *Envelope) Respond(request string) (string, bool) {
    switch {
    case strings.Contains(request, "EHLO"):
        return "250-localhost Hello\n250-SIZE 52428800\n250-PIPELINING\n250-AUTH PLAIN LOGIN\n250 HELP", false
    case strings.Contains(request, "AUTH"):
        auth := strings.TrimPrefix(request, "AUTH PLAIN ")
        decoded, _ := base64.StdEncoding.DecodeString(strings.TrimSpace(auth))
        parts := strings.SplitN(string(decoded), "\x00", 3)
        envelope.AuthUser = parts[1]
        envelope.AuthPass = parts[2]
        return "235 OK, Go ahead", false
    case strings.Contains(request, "MAIL FROM"):
        from := strings.TrimPrefix(request, "MAIL FROM:")
        envelope.From = strings.TrimSpace(from)
        return "250 OK", false
    case strings.Contains(request, "RCPT TO"):
        to := strings.TrimPrefix(request, "RCPT TO:")
        envelope.To = strings.TrimSpace(to)
        return "250 OK", false
    case strings.Contains(request, "DATA"):
        return "354 Go ahead", false
    case strings.TrimSpace(request) == ".":
        return "250 Written safely to disk.", false
    case strings.Contains(request, "QUIT"):
        return "221 localhost saying goodbye", true
    default:
        envelope.Data += strings.TrimSpace(request)
        return "", false
    }
}

var _ = Describe("NotifyUser", func() {
    var fakeUAA *httptest.Server
    var buffer *bytes.Buffer
    var logger *log.Logger
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var handler handlers.NotifyUser
    var envelope Envelope

    BeforeEach(func() {
        fakeUAA = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
            if req.URL.Path == "/Users/user-123" &&
                req.Method == "GET" &&
                req.Header.Get("Authorization") == "Bearer a-special-token" {

                response := `{
                      "id": "user-123",
                      "meta": {
                        "version": 6,
                        "created": "2014-05-22T22:36:36.941Z",
                        "lastModified": "2014-06-25T23:10:03.845Z"
                      },
                      "userName": "admin",
                      "name": {
                        "familyName": "Admin",
                        "givenName": "Mister"
                      },
                      "emails": [
                        {
                          "value": "fake-user@example.com"
                        }
                      ],
                      "groups": [
                        {
                          "value": "e7f74565-4c7e-44ba-b068-b16072cbf08f",
                          "display": "clients.read",
                          "type": "DIRECT"
                        }
                      ],
                      "approvals": [],
                      "active": true,
                      "verified": false,
                      "schemas": [
                        "urn:scim:schemas:core:1.0"
                      ]
                    }`

                w.WriteHeader(http.StatusOK)
                w.Write([]byte(response))
            } else {
                w.WriteHeader(http.StatusNotFound)
            }
        }))

        os.Setenv("UAA_HOST", fakeUAA.URL)
        os.Setenv("UAA_CLIENT_ID", "notifications")
        os.Setenv("UAA_CLIENT_SECRET", "secret")

        addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
        if err != nil {
            panic(err)
        }

        smtpServer, err := net.ListenTCP("tcp", addr)
        if err != nil {
            panic(err)
        }

        parts := strings.SplitN(smtpServer.Addr().String(), ":", 2)
        host := parts[0]
        port := parts[1]

        go func() {
            conn, err := smtpServer.AcceptTCP()
            if err != nil {
                panic(err)
            }
            input := bufio.NewReader(conn)
            output := bufio.NewWriter(conn)

            output.WriteString(fmt.Sprintf("220 %s\r\n", host))
            output.Flush()

            for {
                request, _ := input.ReadString('\n')
                response, exit := envelope.Respond(request)
                if response != "" {
                    output.WriteString(response + "\r\n")
                    output.Flush()
                }
                if exit {
                    break
                }
            }

            conn.Close()
        }()

        os.Setenv("SMTP_USER", "smtp-user")
        os.Setenv("SMTP_PASS", "smtp-pass")
        os.Setenv("SMTP_HOST", host)
        os.Setenv("SMTP_PORT", port)

        envelope = Envelope{}

        buffer = bytes.NewBuffer([]byte{})
        logger = log.New(buffer, "", 0)

        writer = httptest.NewRecorder()

        request, err = http.NewRequest("POST", "/users/user-123", nil)
        if err != nil {
            panic(err)
        }
        request.Header.Set("Authorization", "Bearer a-special-token")

        handler = handlers.NewNotifyUser(logger)
    })

    AfterEach(func() {
        fakeUAA.Close()
    })

    It("logs the email address of the recipient", func() {
        handler.ServeHTTP(writer, request)

        Expect(buffer.String()).To(ContainSubstring("Sending email to fake-user@example.com"))
    })

    It("logs the message envelope", func() {
        handler.ServeHTTP(writer, request)

        data := `From: no-reply@notifications.example.com\nTo: fake-user@example.com\n`
        Expect(buffer.String()).To(ContainSubstring(data))
    })

    It("talks to the SMTP server, sending the email", func() {
        handler.ServeHTTP(writer, request)

        Expect(envelope).To(Equal(Envelope{
            AuthUser: "smtp-user",
            AuthPass: "smtp-pass",
            From:     "<no-reply@notifications.example.com>",
            To:       "<fake-user@example.com>",
            Data:     `From: no-reply@notifications.example.com\nTo: fake-user@example.com\n`,
        }))
    })

    It("returns a status response for the sent mail", func() {
        handler.ServeHTTP(writer, request)

        Expect(writer.Code).To(Equal(http.StatusOK))
        parsed := []map[string]string{}
        err := json.Unmarshal(writer.Body.Bytes(), &parsed)
        if err != nil {
            panic(err)
        }

        Expect(parsed).To(Equal([]map[string]string{
            map[string]string{
                "status": "delivered",
            },
        }))
    })
})
