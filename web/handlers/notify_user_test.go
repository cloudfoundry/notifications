package handlers_test

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyUser", func() {
    var handler handlers.NotifyUser
    var buffer *bytes.Buffer
    var logger *log.Logger
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var token string
    var mailClient FakeMailClient
    var uaaClient FakeUAAClient

    BeforeEach(func() {
        tokenHeader := map[string]interface{}{
            "alg": "FAST",
        }
        tokenClaims := map[string]interface{}{
            "client_id": "mister-client",
            "exp":       3404281214,
            "scope":     []string{"notifications.write"},
        }
        token = BuildToken(tokenHeader, tokenClaims)

        buffer = bytes.NewBuffer([]byte{})
        logger = log.New(buffer, "", 0)
        writer = httptest.NewRecorder()

        mailClient = FakeMailClient{}
        uaaClient = FakeUAAClient{
            UsersByID: map[string]uaa.User{
                "user-123": uaa.User{
                    ID:     "user-123",
                    Emails: []string{"fake-user@example.com"},
                },
                "user-456": uaa.User{
                    ID:     "user-456",
                    Emails: []string{"bounce@example.com"},
                },
            },
        }

        handler = handlers.NewNotifyUser(logger, &mailClient, &uaaClient, FakeGuidGenerator)
    })

    Context("when the request is valid", func() {
        BeforeEach(func() {
            requestBody, err := json.Marshal(map[string]string{
                "kind":               "forgot_password",
                "kind_description":   "Password reminder",
                "source_description": "Login system",
                "text":               "Please reset your password by clicking on this link...",
                "html":               "<p>Please reset your password by clicking on this link...</p>",
            })
            if err != nil {
                panic(err)
            }

            request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer "+token)
        })

        It("logs the email address of the recipient", func() {
            handler.ServeHTTP(writer, request)

            Expect(buffer.String()).To(ContainSubstring("Sending email to fake-user@example.com"))
        })

        It("logs the message envelope", func() {
            handler.ServeHTTP(writer, request)

            data := []string{
                "From: no-reply@notifications.example.com",
                "To: fake-user@example.com",
                "Subject: CF Notification: Password reminder",
                `The following "Password reminder" notification was sent to you directly by the "Login system" component of Cloud Foundry:`,
                "Please reset your password by clicking on this link...",
            }
            results := strings.Split(buffer.String(), "\n")
            for _, item := range data {
                Expect(results).To(ContainElement(item))
            }
        })

        It("talks to the SMTP server, sending the email", func() {
            handler.ServeHTTP(writer, request)

            Expect(len(mailClient.messages)).To(Equal(1))

            msg := mailClient.messages[0]
            Expect(msg).To(Equal(mail.Message{
                From:    "no-reply@notifications.example.com",
                To:      "fake-user@example.com",
                Subject: "CF Notification: Password reminder",
                Body: `
This is a multi-part message in MIME format...

--our-content-boundary
Content-type: text/plain

The following "Password reminder" notification was sent to you directly by the "Login system" component of Cloud Foundry:

Please reset your password by clicking on this link...
--our-content-boundary
Content-Type: text/html
Content-Disposition: inline
Content-Transfer-Encoding: quoted-printable

<html>
    <body>
        <p>The following "Password reminder" notification was sent to you directly by the "Login system" component of Cloud Foundry:</p>

<p>Please reset your password by clicking on this link...</p>
    </body>
</html>
--our-content-boundary--`,
                Headers: []string{
                    "X-CF-Client-ID: mister-client",
                    "X-CF-Notification-ID: deadbeef-aabb-ccdd-eeff-001122334455",
                },
            }))
        })

        It("returns necessary info in the response for the sent mail", func() {
            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(http.StatusOK))
            parsed := []map[string]string{}
            err := json.Unmarshal(writer.Body.Bytes(), &parsed)
            if err != nil {
                panic(err)
            }

            Expect(parsed[0]["status"]).To(Equal("delivered"))
            Expect(parsed[0]["recipient"]).To(Equal("user-123"))
            Expect(parsed[0]["notification_id"]).NotTo(Equal(""))
        })
    })
})
