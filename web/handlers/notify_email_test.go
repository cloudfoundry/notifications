package handlers_test

import (
    "errors"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyEmail", func() {
    Describe("Execute", func() {
        var handler handlers.NotifyEmail
        var writer *httptest.ResponseRecorder
        var errorWriter *FakeErrorWriter
        var fakeNotify *FakeNotify

        BeforeEach(func() {
            errorWriter = &FakeErrorWriter{}
            writer = httptest.NewRecorder()

            fakeNotify = &FakeNotify{}
            handler = handlers.NewNotifyEmail(fakeNotify, errorWriter, nil)
        })

        Context("when notify.execute returns a proper response", func() {
            It("writes that response", func() {
                fakeNotify.Response = []byte("whut")

                handler.Execute(writer, nil, nil)

                Expect(writer.Code).To(Equal(http.StatusOK))
                Expect(fakeNotify.GUID.IsTypeEmail()).To(BeTrue())

                body := string(writer.Body.Bytes())
                Expect(body).To(Equal("whut"))
            })
        })

        Context("when notify.execute errors", func() {
            It("propagates the error", func() {
                fakeNotify.Error = errors.New("Blambo!")
                err := handler.Execute(writer, nil, nil)

                Expect(err).To(Equal(fakeNotify.Error))

            })
        })
    })
})
