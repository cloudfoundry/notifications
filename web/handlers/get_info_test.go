package handlers_test

import (
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/web/handlers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("GetInfo", func() {
    Describe("ServeHTTP", func() {
        var handler handlers.GetInfo

        BeforeEach(func() {
            handler = handlers.NewGetInfo()
        })

        It("returns a 200 response code and an empty JSON body", func() {
            writer := httptest.NewRecorder()
            request, err := http.NewRequest("GET", "/info", nil)
            if err != nil {
                panic(err)
            }

            handler.ServeHTTP(writer, request)

            Expect(writer.Code).To(Equal(http.StatusOK))
            Expect(writer.Body.String()).To(Equal("{}"))
        })
    })
})
