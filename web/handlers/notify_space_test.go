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
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
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

        BeforeEach(func() {
            var err error

            writer = httptest.NewRecorder()
            body, err := json.Marshal(map[string]string{
                "kind": "test_email",
                "text": "This is the body of the email",
            })
            if err != nil {
                panic(err)
            }

            request, err = http.NewRequest("POST", "/spaces/space-001", bytes.NewBuffer(body))
            if err != nil {
                panic(err)
            }

            buffer = bytes.NewBuffer([]byte{})
            logger := log.New(buffer, "", 0)
            fakeCC = NewFakeCloudController()

            fakeUAA := FakeUAAClient{
                ClientToken: uaa.Token{
                    Access: "the-app-token",
                },
                UsersByID: map[string]uaa.User{
                    "user-123": uaa.User{
                        ID:       "user-123",
                        Username: "miss-123",
                        Name: uaa.Name{
                            FamilyName: "123",
                            GivenName:  "Miss",
                        },
                        Emails:   []string{"user-123@example.com"},
                        Active:   true,
                        Verified: false,
                    },
                    "user-456": uaa.User{
                        ID:       "user-456",
                        Username: "mister-456",
                        Name: uaa.Name{
                            FamilyName: "456",
                            GivenName:  "Mister",
                        },
                        Emails:   []string{"user-456@example.com"},
                        Active:   true,
                        Verified: false,
                    },
                },
            }
            mailClient = FakeMailClient{}

            fakeCC.UsersBySpaceGuid["space-001"] = []cf.CloudControllerUser{
                cf.CloudControllerUser{Guid: "user-123"},
                cf.CloudControllerUser{Guid: "user-456"},
            }

            handler = handlers.NewNotifySpace(logger, fakeCC, fakeUAA, &mailClient)
        })

        It("logs the UUIDs of all users in the space", func() {
            handler.ServeHTTP(writer, request)

            Expect(fakeCC.CurrentToken).To(Equal("the-app-token"))

            lines := strings.Split(buffer.String(), "\n")

            Expect(lines).To(ContainElement("user-123"))
            Expect(lines).To(ContainElement("user-456"))
        })

        It("validates the presence of required fields", func() {
            request, err := http.NewRequest("POST", "/spaces/space-001", strings.NewReader(""))
            if err != nil {
                panic(err)
            }

            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(422))
            body := make(map[string]interface{})
            err = json.Unmarshal(writer.Body.Bytes(), &body)
            if err != nil {
                panic(err)
            }

            Expect(body["errors"]).To(ContainElement(`"kind" is a required field`))
            Expect(body["errors"]).To(ContainElement(`"text" is a required field`))
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

            Expect(mailClient.messages).To(ContainElement(mail.Message{
                From:    "no-reply@notifications.example.com",
                To:      "user-123@example.com",
                Subject: "CF Notification: ",
                Body:    "This is the body of the email",
                Headers: []string{"X-CF-Client-ID: ", "X-CF-Notification-ID: "},
            }))

            Expect(mailClient.messages).To(ContainElement(mail.Message{
                From:    "no-reply@notifications.example.com",
                To:      "user-456@example.com",
                Subject: "CF Notification: ",
                Body:    "This is the body of the email",
                Headers: []string{"X-CF-Client-ID: ", "X-CF-Notification-ID: "},
            }))
        })
    })
})
