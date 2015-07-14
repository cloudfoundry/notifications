package network_test

import (
	"net/http"
	"reflect"
	"time"

	"github.com/pivotal-cf-experimental/warrant/internal/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("using a global HTTP client", func() {
	It("retrieves the exact same client reference for a given value of config.SkipVerifySSL", func() {
		transport1 := network.BuildTransport(true)
		transport2 := network.BuildTransport(true)

		transportPointer1 := reflect.ValueOf(transport1).Pointer()
		transportPointer2 := reflect.ValueOf(transport2).Pointer()
		Expect(transportPointer1).To(Equal(transportPointer2))
	})

	It("retrieves difference client references for different values of config.SkipVerifySSL", func() {
		transport1 := network.BuildTransport(true)
		transport2 := network.BuildTransport(false)

		transportPointer1 := reflect.ValueOf(transport1).Pointer()
		transportPointer2 := reflect.ValueOf(transport2).Pointer()
		Expect(transportPointer1).NotTo(Equal(transportPointer2))
	})

	It("builds a correctly configured transport", func() {
		transport := network.BuildTransport(true).(*http.Transport)

		Expect(transport.TLSClientConfig.InsecureSkipVerify).To(BeTrue())
		Expect(reflect.ValueOf(transport.Proxy).Pointer()).To(Equal(reflect.ValueOf(http.ProxyFromEnvironment).Pointer()))
		Expect(transport.TLSHandshakeTimeout).To(Equal(10 * time.Second))
	})
})
