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
		token = GetClientTokenFor("my-client")
	})

	It("can create and read a new sender", func() {
		var senderID string

		By("creating a sender", func() {
			status, response, err := client.Do("POST", "/senders", map[string]interface{}{
				"name": "My Cool App",
			}, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			Expect(response["name"]).To(Equal("My Cool App"))

			senderID = response["id"].(string)
		})

		By("getting the sender", func() {
			status, response, err := client.Do("GET", fmt.Sprintf("/senders/%s", senderID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(response["id"]).To(Equal(senderID))
			Expect(response["name"]).To(Equal("My Cool App"))
		})
	})

	It("returns a 201 when a sender already exists", func() {
		var senderID string

		By("creating a sender", func() {
			status, response, err := client.Do("POST", "/senders", map[string]interface{}{
				"name": "My Cool App",
			}, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			senderID = response["id"].(string)
		})

		By("creating another sender with the same name", func() {
			status, response, err := client.Do("POST", "/senders", map[string]interface{}{
				"name": "My Cool App",
			}, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))
			Expect(response["id"]).To(Equal(senderID))
		})
	})

	Context("failure states", func() {
		It("returns a 404 when the sender cannot be retrieved", func() {
			status, response, err := client.Do("GET", "/senders/missing-sender-id", nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNotFound))
			Expect(response["errors"]).To(ContainElement("Sender with id \"missing-sender-id\" could not be found"))
		})

		It("returns a 404 when the sender belongs to another client", func() {
			var senderID string

			By("creating a sender for one client", func() {
				status, response, err := client.Do("POST", "/senders", map[string]interface{}{
					"name": "My Cool App",
				}, token.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				senderID = response["id"].(string)
			})

			By("attempting to access the created sender as another client", func() {
				token := GetClientTokenFor("other-client")
				status, response, err := client.Do("GET", fmt.Sprintf("/senders/%s", senderID), nil, token.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement(fmt.Sprintf("Sender with id %q could not be found", senderID)))
			})
		})
	})
})
