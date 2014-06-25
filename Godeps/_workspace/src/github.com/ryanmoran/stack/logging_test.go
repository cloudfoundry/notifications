package stack_test

import (
    "bytes"
    "log"
    "net/http"
    "net/http/httptest"
    "os"

    "github.com/ryanmoran/stack"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Logging", func() {
    var ware stack.Logging
    var logFile *bytes.Buffer

    BeforeEach(func() {
        logFile = bytes.NewBuffer([]byte{})
        logger := log.New(logFile, "", log.LstdFlags)
        ware = stack.NewLogging(logger)
    })

    It("logs the incoming requests when the ENABLE_HTTP_LOGGING is set", func() {
        os.Setenv("HTTP_LOGGING_ENABLED", "true")
        writer := httptest.NewRecorder()
        request, err := http.NewRequest("GET", "/some/random/request", nil)
        if err != nil {
            panic(err)
        }

        result := ware.ServeHTTP(writer, request)

        Expect(result).To(BeTrue())
        Expect(logFile.String()).To(ContainSubstring("GET /some/random/request"))
    })

    It("does nothing when the ENABLE_HTTP_LOGGING is not set", func() {
        os.Setenv("HTTP_LOGGING_ENABLED", "")
        writer := httptest.NewRecorder()
        request, err := http.NewRequest("GET", "/some/random/request", nil)
        if err != nil {
            panic(err)
        }

        result := ware.ServeHTTP(writer, request)

        Expect(result).To(BeTrue())
        Expect(logFile.Len()).To(Equal(0))
    })
})
