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

var _ = Describe("NotifyOrganization", func() {
    Describe("Execute", func() {
        var handler handlers.NotifyOrganization
        var writer *httptest.ResponseRecorder
        var request *http.Request
        var fakeNotify *fakes.FakeNotify
        var context stack.Context

        BeforeEach(func() {
            var err error

            writer = httptest.NewRecorder()
            request, err = http.NewRequest("POST", "/organizations/org-001", nil)
            if err != nil {
                panic(err)
            }
            context = stack.NewContext()

            fakeNotify = &fakes.FakeNotify{}
            fakeDatabase := fakes.NewDatabase()
            handler = handlers.NewNotifyOrganization(fakeNotify, nil, nil, fakeDatabase)
        })

        Context("when the notify.Execute returns a successful response", func() {
            It("returns the JSON representation of the response", func() {
                fakeNotify.Response = []byte("whatever")
                recipe := postal.OrganizationRecipe{}

                handler.Execute(writer, request, nil, context, recipe)

                Expect(writer.Code).To(Equal(http.StatusOK))
                Expect(fakeNotify.GUID.String()).To(Equal("org-001"))
                Expect(fakeNotify.GUID.BelongsToOrganization()).To(BeTrue())

                body := string(writer.Body.Bytes())
                Expect(body).To(Equal("whatever"))
            })
        })

        Context("when the notify.Execute returns an error", func() {
            It("propagates the error", func() {
                fakeNotify.Error = errors.New("the error")
                recipe := postal.OrganizationRecipe{}

                err := handler.Execute(writer, request, nil, context, recipe)
                Expect(err).To(Equal(fakeNotify.Error))
            })
        })
    })
})
