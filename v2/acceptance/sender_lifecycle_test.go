package v2

import (
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v2/acceptance/support"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sender lifecycle", func() {
	var (
		client *support.Client
		token  uaa.Token
	)

	BeforeEach(func() {
		client = support.NewClient(support.Config{
			Host:  Servers.Notifications.URL(),
			Trace: Trace,
		})
		token = GetClientTokenFor("my-client", "uaa")
	})

	It("can create and read a new sender", func() {
		var sender struct {
			ID   string
			Name string
		}

		By("creating a sender", func() {
			status, response, err := client.Do("POST", "/senders", map[string]interface{}{
				"name": "My Cool App",
			}, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			sender.ID = response["id"].(string)
			sender.Name = response["name"].(string)

			Expect(sender.Name).To(Equal("My Cool App"))
		})

		By("getting the sender", func() {
			status, response, err := client.Do("GET", fmt.Sprintf("/senders/%s", sender.ID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(response["id"]).To(Equal(sender.ID))
			Expect(response["name"]).To(Equal("My Cool App"))
		})
	})
})
