package handlers_test

import (
    "bytes"
    "encoding/json"
    "errors"
    "log"
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/nu7hatch/gouuid"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifySpace", func() {
    Describe("ServeHTTP", func() {
        var handler handlers.NotifySpace
        var writer *httptest.ResponseRecorder
        var request *http.Request
        var buffer *bytes.Buffer
        var fakeCC *FakeCloudController
        var mailClient FakeMailClient
        var token string
        var logger *log.Logger
        var fakeUAA uaa.UAAInterface

        BeforeEach(func() {
            var err error

            writer = httptest.NewRecorder()
            body, err := json.Marshal(map[string]string{
                "kind":               "test_email",
                "text":               "This is the plain text body of the email",
                "html":               "<p>This is the HTML Body of the email</p>",
                "subject":            "Your instance is down",
                "source_description": "MySQL Service",
                "kind_description":   "Instance Alert",
            })
            if err != nil {
                panic(err)
            }

            tokenHeader := map[string]interface{}{
                "alg": "FAST",
            }
            tokenClaims := map[string]interface{}{
                "client_id": "mister-client",
                "exp":       3404281214,
                "scope":     []string{"notifications.write"},
            }
            token = BuildToken(tokenHeader, tokenClaims)

            request, err = http.NewRequest("POST", "/spaces/space-001", bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer "+token)

            buffer = bytes.NewBuffer([]byte{})
            logger = log.New(buffer, "", 0)

            fakeUAA = FakeUAAClient{
                ClientToken: uaa.Token{
                    Access: token,
                },
                UsersByID: map[string]uaa.User{
                    "user-123": uaa.User{
                        ID:     "user-123",
                        Emails: []string{"user-123@example.com"},
                    },
                    "user-456": uaa.User{
                        ID:     "user-456",
                        Emails: []string{"user-456@example.com"},
                    },
                },
            }
            mailClient = FakeMailClient{}

            fakeCC = NewFakeCloudController()
            fakeCC.UsersBySpaceGuid["space-001"] = []cf.CloudControllerUser{
                cf.CloudControllerUser{Guid: "user-123"},
                cf.CloudControllerUser{Guid: "user-456"},
            }
            fakeCC.Spaces = map[string]cf.CloudControllerSpace{
                "space-001": cf.CloudControllerSpace{
                    Name:             "production",
                    Guid:             "space-001",
                    OrganizationGuid: "org-001",
                },
            }
            fakeCC.Orgs = map[string]cf.CloudControllerOrganization{
                "org-001": cf.CloudControllerOrganization{
                    Name: "pivotaltracker",
                },
            }

            handler = handlers.NewNotifySpace(logger, fakeCC, fakeUAA, &mailClient, FakeGuidGenerator)
        })

        It("logs the UUIDs of all users in the space", func() {
            handler.ServeHTTP(writer, request)

            Expect(fakeCC.CurrentToken).To(Equal(token))

            lines := strings.Split(buffer.String(), "\n")

            Expect(lines).To(ContainElement("user-123"))
            Expect(lines).To(ContainElement("user-456"))
        })

        It("validates the presence of required fields", func() {
            request, err := http.NewRequest("POST", "/spaces/space-001", strings.NewReader(""))
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer "+token)

            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(422))
            body := make(map[string]interface{})
            err = json.Unmarshal(writer.Body.Bytes(), &body)
            if err != nil {
                panic(err)
            }

            Expect(body["errors"]).To(ContainElement(`"kind" is a required field`))
            Expect(body["errors"]).To(ContainElement(`"text" or "html" fields must be supplied`))
        })

        It("returns a 502 when CloudController fails to respond", func() {
            fakeCC.GetUsersBySpaceGuidError = errors.New("BOOM!")

            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(http.StatusBadGateway))
            body := make(map[string]interface{})
            err := json.Unmarshal(writer.Body.Bytes(), &body)
            if err != nil {
                panic(err)
            }

            Expect(body["errors"]).To(ContainElement("Cloud Controller is unavailable"))
        })

        It("sends mail to the users in the space", func() {
            handler.ServeHTTP(writer, request)

            Expect(len(mailClient.messages)).To(Equal(2))

            body := `
This is a multi-part message in MIME format...

--our-content-boundary
Content-type: text/plain

The following "Instance Alert" notification was sent to you by the "MySQL Service" component of Cloud Foundry because you are a member of the "production" space in the "pivotaltracker" organization:

This is the plain text body of the email
--our-content-boundary
Content-Type: text/html
Content-Disposition: inline
Content-Transfer-Encoding: quoted-printable

<html>
    <body>
        <p>The following "Instance Alert" notification was sent to you by the "MySQL Service" component of Cloud Foundry because you are a member of the "production" space in the "pivotaltracker" organization:</p>

<p>This is the HTML Body of the email</p>
    </body>
</html>
--our-content-boundary--`

            firstMessage := mailClient.messages[0]
            Expect(firstMessage.From).To(Equal("no-reply@notifications.example.com"))
            Expect(firstMessage.To).To(Equal("user-123@example.com"))
            Expect(firstMessage.Subject).To(Equal("CF Notification: Your instance is down"))
            Expect(firstMessage.Body).To(Equal(body))
            Expect(firstMessage.Headers).To(Equal([]string{
                "X-CF-Client-ID: mister-client",
                "X-CF-Notification-ID: deadbeef-aabb-ccdd-eeff-001122334455",
            }))

            secondMessage := mailClient.messages[1]
            Expect(secondMessage.From).To(Equal("no-reply@notifications.example.com"))
            Expect(secondMessage.To).To(Equal("user-456@example.com"))
            Expect(secondMessage.Subject).To(Equal("CF Notification: Your instance is down"))
            Expect(secondMessage.Body).To(Equal(body))
            Expect(secondMessage.Headers).To(Equal([]string{
                "X-CF-Client-ID: mister-client",
                "X-CF-Notification-ID: deadbeef-aabb-ccdd-eeff-001122334455",
            }))
        })

        It("returns necessary info in the response for the sent mail", func() {

            handler = handlers.NewNotifySpace(logger, fakeCC, fakeUAA, &mailClient, func() (*uuid.UUID, error) {
                guid, err := uuid.NewV4()
                if err != nil {
                    panic(err)
                }
                return guid, nil
            })

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

            Expect(parsed[1]["status"]).To(Equal("delivered"))
            Expect(parsed[1]["recipient"]).To(Equal("user-456"))
            Expect(parsed[1]["notification_id"]).NotTo(Equal(parsed[0]["notification_id"]))
        })

        Context("when the SMTP server fails to deliver the mail", func() {
            It("returns a status indicating that delivery failed", func() {
                mailClient.errorOnSend = true
                handler.ServeHTTP(writer, request)

                Expect(writer.Code).To(Equal(http.StatusOK))
                parsed := []map[string]string{}
                err := json.Unmarshal(writer.Body.Bytes(), &parsed)
                if err != nil {
                    panic(err)
                }

                Expect(parsed[0]["status"]).To(Equal("failed"))
                Expect(parsed[1]["status"]).To(Equal("failed"))
            })
        })

        Context("when the SMTP server cannot be reached", func() {
            It("returns a status indicating that the server is unavailable", func() {
                mailClient.errorOnConnect = true
                handler.ServeHTTP(writer, request)

                Expect(writer.Code).To(Equal(http.StatusOK))
                parsed := []map[string]string{}
                err := json.Unmarshal(writer.Body.Bytes(), &parsed)
                if err != nil {
                    panic(err)
                }

                Expect(parsed[0]["status"]).To(Equal("unavailable"))
                Expect(parsed[1]["status"]).To(Equal("unavailable"))
            })
        })
    })
})
