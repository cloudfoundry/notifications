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
    It("recovers from panics", func() {
        writer := httptest.NewRecorder()
        request, err := http.NewRequest("GET", "/my/path", nil)
        if err != nil {
            panic(err)
        }

        Expect(func() {
            defer stack.Recover(writer, request, nil)
            panic(errors.New("Random Error"))
        }).NotTo(Panic())
    })

    It("returns a 500 status code and error message", func() {
        writer := httptest.NewRecorder()
        request, err := http.NewRequest("GET", "/my/path", nil)
        if err != nil {
            panic(err)
        }

        func() {
            defer stack.Recover(writer, request, nil)
            panic(errors.New("Random Error"))
        }()

        Expect(writer.Code).To(Equal(http.StatusInternalServerError))
        Expect(writer.Body.String()).To(ContainSubstring("Internal Server Error"))
    })
})
