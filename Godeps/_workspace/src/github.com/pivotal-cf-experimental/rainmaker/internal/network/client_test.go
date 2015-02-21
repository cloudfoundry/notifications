package network_test

import (
	"net/http"
	"reflect"
	"time"

	"github.com/pivotal-cf-experimental/rainmaker/internal/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("using a global HTTP client", func() {
	BeforeEach(func() {
		network.ClearHTTPClient()
	})

	It("retrieves the exact same client reference", func() {
		client1 := network.GetClient(network.Config{})
		client2 := network.GetClient(network.Config{})

		Expect(client1).To(BeAssignableToTypeOf(&http.Client{}))
		Expect(client1).ToNot(BeNil())
		Expect(reflect.ValueOf(client1).Pointer()).To(Equal(reflect.ValueOf(client2).Pointer()))

	})

	It("uses the configuration to configure the HTTP client", func() {
		config := network.Config{
			SkipVerifySSL: true,
		}
		transport := network.GetClient(config).Transport.(*http.Transport)

		Expect(transport.TLSClientConfig.InsecureSkipVerify).To(BeTrue())
		Expect(reflect.ValueOf(transport.Proxy).Pointer()).To(Equal(reflect.ValueOf(http.ProxyFromEnvironment).Pointer()))
		Expect(transport.TLSHandshakeTimeout).To(Equal(10 * time.Second))
	})
})
