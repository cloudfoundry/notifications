package utilities_test

import (
    "errors"
    "net/http"
    "net/http/httptest"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "github.com/pivotal-cf/cf-notifications/web/utilities"
)

var _ = Describe("Recover", func() {
    It("recovers with a 500 and HTML response for HTML request", func() {
        writer := httptest.NewRecorder()
        request, err := http.NewRequest("GET", "/my/path", nil)
        if err != nil {
            panic(err)
        }

        Expect(func() {
            defer utilities.Recover(writer, request)
            panic(errors.New("Random Error"))
        }).NotTo(Panic())
    })
})
