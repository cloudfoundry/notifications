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
        var errorWriter *fakes.ErrorWriter
        var notify *fakes.Notify
        var context stack.Context

        BeforeEach(func() {
            errorWriter = fakes.NewErrorWriter()
            writer = httptest.NewRecorder()
            context = stack.NewContext()
            database := fakes.NewDatabase()

            notify = fakes.NewNotify()
            handler = handlers.NewNotifyEmail(notify, errorWriter, nil, database)
        })

        Context("when notify.Execute returns a proper response", func() {
            It("writes that response", func() {
                notify.Response = []byte("whut")

                handler.Execute(writer, nil, nil, context)

                Expect(writer.Code).To(Equal(http.StatusOK))
                Expect(notify.GUID.IsTypeEmail()).To(BeTrue())

                body := string(writer.Body.Bytes())
                Expect(body).To(Equal("whut"))
            })
        })

        Context("when notify.Execute errors", func() {
            It("propagates the error", func() {
                notify.Error = errors.New("Blambo!")
                err := handler.Execute(writer, nil, nil, context)

                Expect(err).To(Equal(notify.Error))

            })
        })
    })
})
