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
                "kind":               "test_email",
                "text":               "This is the body of the email",
                "subject":            "Your instance is down",
                "source_description": "MySQL Service",
                "kind_description":   "Instance Alert",
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

            body := `The following "Instance Alert" notification was sent to you by the "MySQL Service" component of Cloud Foundry because you are a member of the "production" space in the "pivotaltracker" organization:

This is the body of the email`

            firstMessage := mailClient.messages[0]
            Expect(firstMessage.From).To(Equal("no-reply@notifications.example.com"))
            Expect(firstMessage.To).To(Equal("user-123@example.com"))
            Expect(firstMessage.Subject).To(Equal("CF Notification: Your instance is down"))
            Expect(firstMessage.Body).To(Equal(body))
            Expect(firstMessage.Headers).To(Equal([]string{"X-CF-Client-ID: ", "X-CF-Notification-ID: "}))

            secondMessage := mailClient.messages[1]
            Expect(secondMessage.From).To(Equal("no-reply@notifications.example.com"))
            Expect(secondMessage.To).To(Equal("user-456@example.com"))
            Expect(secondMessage.Subject).To(Equal("CF Notification: Your instance is down"))
            Expect(secondMessage.Body).To(Equal(body))
            Expect(secondMessage.Headers).To(Equal([]string{"X-CF-Client-ID: ", "X-CF-Notification-ID: "}))
        })
    })
})
