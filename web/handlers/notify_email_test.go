package handlers_test

import (
    "errors"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/ryanmoran/stack"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifyEmail", func() {
    Describe("Execute", func() {
        var handler handlers.NotifyEmail
        var writer *httptest.ResponseRecorder
        var errorWriter *fakes.FakeErrorWriter
        var fakeNotify *FakeNotify
        var context stack.Context

        BeforeEach(func() {
            errorWriter = &fakes.FakeErrorWriter{}
            writer = httptest.NewRecorder()
            context = stack.NewContext()
            database := fakes.NewDatabase()

            fakeNotify = &FakeNotify{}
            handler = handlers.NewNotifyEmail(fakeNotify, errorWriter, nil, database)
        })

        Context("when notify.Execute returns a proper response", func() {
            It("writes that response", func() {
                fakeNotify.Response = []byte("whut")

                handler.Execute(writer, nil, nil, context)

                Expect(writer.Code).To(Equal(http.StatusOK))
                Expect(fakeNotify.GUID.IsTypeEmail()).To(BeTrue())

                body := string(writer.Body.Bytes())
                Expect(body).To(Equal("whut"))
            })
        })

        Context("when notify.Execute errors", func() {
            It("propagates the error", func() {
                fakeNotify.Error = errors.New("Blambo!")
                err := handler.Execute(writer, nil, nil, context)

                Expect(err).To(Equal(fakeNotify.Error))

            })
        })
    })
})
