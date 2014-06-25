package stack_test

import (
    "errors"
    "net/http"
    "net/http/httptest"

    "github.com/ryanmoran/stack"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Recover", func() {
    It("recovers with a 500 and HTML response for HTML request", func() {
        writer := httptest.NewRecorder()
        request, err := http.NewRequest("GET", "/my/path", nil)
        if err != nil {
            panic(err)
        }

        Expect(func() {
            defer stack.Recover(writer, request)
            panic(errors.New("Random Error"))
        }).NotTo(Panic())
    })
})
