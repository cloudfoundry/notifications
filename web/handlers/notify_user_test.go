package handlers_test

import (
    "errors"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyUser", func() {
    Context("Execute", func() {
        var handler handlers.NotifyUser
        var writer *httptest.ResponseRecorder
        var request *http.Request
        var errorWriter *FakeErrorWriter
        var fakeNotify *FakeNotify

        BeforeEach(func() {
            var err error
            errorWriter = &FakeErrorWriter{}

            writer = httptest.NewRecorder()

            request, err = http.NewRequest("POST", "/users/user-123", nil)
            if err != nil {
                panic(err)
            }

            fakeNotify = &FakeNotify{}
            handler = handlers.NewNotifyUser(fakeNotify, errorWriter, nil)
        })

        Context("when notify.Execute returns a successful response", func() {
            It("returns the JSON representation of the response", func() {
                fakeNotify.Response = []byte("whut")

                handler.Execute(writer, request, nil)

                Expect(writer.Code).To(Equal(http.StatusOK))
                body := string(writer.Body.Bytes())
                Expect(body).To(Equal("whut"))

                Expect(fakeNotify.GUID.String()).To(Equal("user-123"))
                Expect(fakeNotify.GUID.BelongsToSpace()).To(BeFalse())
                Expect(fakeNotify.GUID.IsTypeEmail()).To(BeFalse())

            })
        })

        Context("when notify.Execute returns an error", func() {
            It("propagates the error", func() {
                fakeNotify.Error = errors.New("BOOM!")

                err := handler.Execute(writer, request, nil)

                Expect(err).To(Equal(fakeNotify.Error))

            })
        })
    })
})
