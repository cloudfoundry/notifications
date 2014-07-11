package handlers_test

import (
    "bytes"
    "log"
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

type FakeCloudController struct {
    UsersBySpaceGuid map[string][]cf.CloudControllerUser
    CurrentToken     string
}

func NewFakeCloudController() *FakeCloudController {
    return &FakeCloudController{
        UsersBySpaceGuid: make(map[string][]cf.CloudControllerUser),
    }
}

func (fake *FakeCloudController) GetUsersBySpaceGuid(guid, token string) ([]cf.CloudControllerUser, error) {
    fake.CurrentToken = token
    if users, ok := fake.UsersBySpaceGuid[guid]; ok {
        return users, nil
    } else {
        return make([]cf.CloudControllerUser, 0), nil
    }
}

var _ = Describe("NotifySpace", func() {
    Describe("ServeHTTP", func() {
        var handler handlers.NotifySpace
        var writer *httptest.ResponseRecorder
        var request *http.Request
        var buffer *bytes.Buffer
        var fakeCC *FakeCloudController

        BeforeEach(func() {
            var err error

            writer = httptest.NewRecorder()
            request, err = http.NewRequest("POST", "/spaces/space-001", nil)
            if err != nil {
                panic(err)
            }
            request.Header.Set("Authorization", "Bearer special-token-value")

            buffer = bytes.NewBuffer([]byte{})
            logger := log.New(buffer, "", 0)
            fakeCC = NewFakeCloudController()

            handler = handlers.NewNotifySpace(logger, fakeCC)
        })

        It("logs the UUIDs of all users in the space", func() {
            fakeCC.UsersBySpaceGuid["space-001"] = []cf.CloudControllerUser{
                cf.CloudControllerUser{Guid: "user-123"},
                cf.CloudControllerUser{Guid: "user-456"},
                cf.CloudControllerUser{Guid: "user-789"},
            }

            handler.ServeHTTP(writer, request)

            Expect(fakeCC.CurrentToken).To(Equal("special-token-value"))

            lines := strings.Split(buffer.String(), "\n")

            Expect(lines).To(ContainElement("user-123"))
            Expect(lines).To(ContainElement("user-456"))
            Expect(lines).To(ContainElement("user-789"))
        })
    })
})
