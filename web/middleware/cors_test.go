package middleware_test

import (
    "net/http"
    "net/http/httptest"
    "os"

    "github.com/cloudfoundry-incubator/notifications/web/middleware"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("CORS", func() {
    var writer *httptest.ResponseRecorder
    var request *http.Request
    var ware middleware.CORS
    var corsOrigin string

    Describe("ServeHTTP", func() {
        BeforeEach(func() {
            corsOrigin = os.Getenv("CORS_ORIGIN")
            os.Setenv("CORS_ORIGIN", "test-cors-origin")
        })

        AfterEach(func() {
            os.Setenv("CORS_ORIGIN", corsOrigin)
        })

        BeforeEach(func() {
            var err error

            writer = httptest.NewRecorder()
            request, err = http.NewRequest("OPTIONS", "/user_preferences", nil)
            if err != nil {
                panic(err)
            }

            ware = middleware.NewCORS()
        })

        It("sets the correct CORS headers", func() {
            result := ware.ServeHTTP(writer, request, nil)

            Expect(result).To(BeTrue())
            Expect(writer.HeaderMap.Get("Access-Control-Allow-Origin")).To(Equal("test-cors-origin"))
            Expect(writer.HeaderMap.Get("Access-Control-Allow-Methods")).To(Equal("GET, PATCH"))
            Expect(writer.HeaderMap.Get("Access-Control-Allow-Headers")).To(Equal("Accept, Authorization, Content-Type"))
        })
    })
})
