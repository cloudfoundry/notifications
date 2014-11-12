package cf_test

import (
	"net/http"
	"reflect"

	"github.com/cloudfoundry-incubator/notifications/cf"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetClient", func() {
	It("Gets the global http Client", func() {
		ccClient := cf.NewCloudControllerClient("http://api.example.com", false)
		client1 := cf.GetClient(ccClient)
		client2 := cf.GetClient(ccClient)

		Expect(client1).To(BeAssignableToTypeOf(&http.Client{}))
		Expect(client1).ToNot(BeNil())
		Expect(reflect.ValueOf(client1).Pointer()).To(Equal(reflect.ValueOf(client2).Pointer()))
	})
})
