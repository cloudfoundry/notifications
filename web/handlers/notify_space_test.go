package handlers_test

import (
    "errors"
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotifySpace", func() {
    Describe("Execute", func() {
        var handler handlers.NotifySpace
        var writer *httptest.ResponseRecorder
        var request *http.Request
        var fakeNotify *FakeNotify

        BeforeEach(func() {
            var err error

            writer = httptest.NewRecorder()
            request, err = http.NewRequest("POST", "/spaces/space-001", nil)
            if err != nil {
                panic(err)
            }

            fakeNotify = &FakeNotify{}
            handler = handlers.NewNotifySpace(fakeNotify, nil, nil)
        })

        Context("when the notify.Execute returns a successful response", func() {
            It("returns the JSON representation of the response", func() {
                fakeNotify.Response = []byte("whatever")

                handler.Execute(writer, request, nil)

                Expect(writer.Code).To(Equal(http.StatusOK))
                Expect(fakeNotify.GUID.String()).To(Equal("space-001"))
                Expect(fakeNotify.GUID.BelongsToSpace()).To(BeTrue())

                body := string(writer.Body.Bytes())
                Expect(body).To(Equal("whatever"))
            })
        })

        Context("when the notify.Execute returns an error", func() {
            It("propagates the error", func() {
                fakeNotify.Error = errors.New("the error")

                err := handler.Execute(writer, request, nil)
                Expect(err).To(Equal(fakeNotify.Error))
            })
        })
    })
})
