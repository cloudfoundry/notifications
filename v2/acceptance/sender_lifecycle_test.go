package acceptance

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

	It("can create, list, update and read a new sender", func() {
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

		By("listing all senders", func() {
			status, response, err := client.Do("GET", "/senders", nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			senders := response["senders"].([]interface{})
			Expect(len(senders)).To(Equal(1))
			sender := senders[0].(map[string]interface{})
			Expect(sender["id"]).To(Equal(senderID))
		})

		By("getting the sender", func() {
			status, response, err := client.Do("GET", fmt.Sprintf("/senders/%s", senderID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(response["id"]).To(Equal(senderID))
			Expect(response["name"]).To(Equal("My Cool App"))
		})

		By("updating the sender", func() {
			status, response, err := client.Do("PUT", fmt.Sprintf("/senders/%s", senderID),
				map[string]interface{}{
					"name": "My Not Cool App",
				}, token.Access)

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			Expect(response["id"]).To(Equal(senderID))
			Expect(response["name"]).To(Equal("My Not Cool App"))
		})

		By("getting the updated sender", func() {
			status, response, err := client.Do("GET", fmt.Sprintf("/senders/%s", senderID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(response["id"]).To(Equal(senderID))
			Expect(response["name"]).To(Equal("My Not Cool App"))
		})
	})

	It("can delete a sender and associated campaign types", func() {
		var senderID, campaignTypeID string
		By("creating a sender", func() {
			status, response, err := client.Do("POST", "/senders", map[string]interface{}{
				"name": "My Cool App",
			}, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			Expect(response["name"]).To(Equal("My Cool App"))

			senderID = response["id"].(string)
		})

		By("creating a campaign type", func() {
			status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
				"name":        "some-campaign-type",
				"description": "a great campaign type",
			}, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			campaignTypeID = response["id"].(string)
		})

		By("deleting the sender", func() {
			status, _, err := client.Do("DELETE", fmt.Sprintf("/senders/%s", senderID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("getting the deleted sender", func() {
			status, _, err := client.Do("GET", fmt.Sprintf("/senders/%s", senderID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNotFound))
		})

		By("getting the deleted campaign type", func() {
			status, _, err := client.Do("GET", fmt.Sprintf("/senders/%s/campaign_types/%s", senderID, campaignTypeID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNotFound))
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
		Context("when getting a sender", func() {
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

		Context("when updating a sender", func() {
			It("returns a 404 when the sender cannot be retrieved", func() {
				status, response, err := client.Do("PUT", "/senders/missing-sender-id", map[string]interface{}{
					"name": "My Not Cool App",
				}, token.Access)
				Expect(err).NotTo(HaveOccurred())

				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement("Sender with id \"missing-sender-id\" could not be found"))
			})

			It("returns a 404 when the client does not own the sender", func() {
				var senderID string

				By("creating a sender with a different token", func() {
					status, response, err := client.Do("POST", "/senders", map[string]interface{}{
						"name": "My Cool App",
					}, token.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusCreated))

					senderID = response["id"].(string)
				})

				By("attempting to get the sender with a different token", func() {
					otherToken := GetClientTokenFor("otherclient")

					status, response, err := client.Do("GET", fmt.Sprintf("/senders/%s", senderID), nil, otherToken.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusNotFound))
					Expect(response["errors"]).To(ContainElement(fmt.Sprintf("Sender with id %q could not be found", senderID)))
				})
			})
		})

		Context("when deleting a sender", func() {
			It("returns a 404 when the sender cannot be retrieved", func() {
				status, response, err := client.Do("DELETE", "/senders/missing-sender-id", nil, token.Access)
				Expect(err).NotTo(HaveOccurred())

				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement("Sender with id \"missing-sender-id\" could not be found"))
			})

			It("returns a 404 when the client does not own the sender", func() {
				var senderID string

				By("creating a sender with a different token", func() {
					status, response, err := client.Do("POST", "/senders", map[string]interface{}{
						"name": "My Cool App",
					}, token.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusCreated))

					senderID = response["id"].(string)
				})

				By("attempting to deletd the sender with a different token", func() {
					otherToken := GetClientTokenFor("otherclient")

					status, response, err := client.Do("DELETE", fmt.Sprintf("/senders/%s", senderID), nil, otherToken.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusNotFound))
					Expect(response["errors"]).To(ContainElement(fmt.Sprintf("Sender with id %q could not be found", senderID)))
				})
			})
		})
	})
})
