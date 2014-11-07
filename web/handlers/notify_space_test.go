package handlers_test

import (
    "errors"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/ryanmoran/stack"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifySpace", func() {
    Describe("Execute", func() {
        var handler handlers.NotifySpace
        var writer *httptest.ResponseRecorder
        var request *http.Request
        var notify *fakes.Notify
        var context stack.Context

        BeforeEach(func() {
            var err error

            writer = httptest.NewRecorder()
            request, err = http.NewRequest("POST", "/spaces/space-001", nil)
            if err != nil {
                panic(err)
            }
            context = stack.NewContext()

            notify = fakes.NewNotify()
            fakeDatabase := fakes.NewDatabase()
            handler = handlers.NewNotifySpace(notify, nil, nil, fakeDatabase)
        })

        Context("when the notify.Execute returns a successful response", func() {
            It("returns the JSON representation of the response", func() {
                notify.Response = []byte("whatever")
                strategy := postal.SpaceStrategy{}

                handler.Execute(writer, request, nil, context, strategy)

                Expect(writer.Code).To(Equal(http.StatusOK))
                Expect(notify.GUID.String()).To(Equal("space-001"))
                Expect(notify.GUID.BelongsToSpace()).To(BeTrue())

                body := string(writer.Body.Bytes())
                Expect(body).To(Equal("whatever"))
            })
        })

        Context("when the notify.Execute returns an error", func() {
            It("propagates the error", func() {
                notify.Error = errors.New("the error")
                strategy := postal.SpaceStrategy{}

                err := handler.Execute(writer, request, nil, context, strategy)
                Expect(err).To(Equal(notify.Error))
            })
        })
    })
})
