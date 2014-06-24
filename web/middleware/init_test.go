package middleware_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestMiddlewareSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Middleware Suite")
}
