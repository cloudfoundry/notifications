package handlers_test

import (
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/ryanmoran/stack"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("OptionsPreferences", func() {
    var handler handlers.OptionsPreferences
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var context stack.Context

    BeforeEach(func() {
        var err error
        writer = httptest.NewRecorder()
        request, err = http.NewRequest("OPTIONS", "/user_preferences", nil)
        if err != nil {
            panic(err)
        }
        context = stack.NewContext()
        handler = handlers.NewOptionsPreferences()
    })

    Describe("ServeHTTP", func() {
        It("returns a 204 status code", func() {
            handler.ServeHTTP(writer, request, context)

            Expect(writer.Code).To(Equal(http.StatusNoContent))
        })
    })
})
