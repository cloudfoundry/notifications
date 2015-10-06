package acceptance

import (
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v2/acceptance/support"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type senderResponse struct {
	ID    string
	Name  string
	Links struct {
		Self struct {
			Href string
		}

		CampaignTypes struct {
			Href string
		} `json:"campaign_types"`

		Campaigns struct {
			Href string
		}
	} `json:"_links"`
}

type sendersListResponse struct {
	Senders []senderResponse
	Links   struct {
		Self struct {
			Href string
		}
	} `json:"_links"`
}

var _ = Describe("Sender lifecycle", func() {
	var (
		client *support.Client
		token  uaa.Token
	)

	BeforeEach(func() {
		client = support.NewClient(support.Config{
			Host:              Servers.Notifications.URL(),
			Trace:             Trace,
			RoundTripRecorder: roundtripRecorder,
		})
		token = GetClientTokenFor("my-client")

	})

	It("can create, list, update and read a new sender", func() {
		var senderID string

		By("creating a sender", func() {
			var response senderResponse
			client.Document("sender-create")
			status, err := client.DoTyped("POST", "/senders", map[string]interface{}{
				"name": "My Cool App",
			}, token.Access, &response)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			Expect(response.Name).To(Equal("My Cool App"))
			Expect(response.Links.Self.Href).To(Equal(fmt.Sprintf("/senders/%s", response.ID)))
			Expect(response.Links.CampaignTypes.Href).To(Equal(fmt.Sprintf("/senders/%s/campaign_types", response.ID)))
			Expect(response.Links.Campaigns.Href).To(Equal(fmt.Sprintf("/senders/%s/campaigns", response.ID)))

			senderID = response.ID
		})

		By("listing all senders", func() {
			var response sendersListResponse
			client.Document("sender-list")
			status, err := client.DoTyped("GET", "/senders", nil, token.Access, &response)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			Expect(response.Senders).To(HaveLen(1))
			Expect(response.Links.Self.Href).To(Equal("/senders"))
			sender := response.Senders[0]
			Expect(sender.ID).To(Equal(senderID))
			Expect(sender.Name).To(Equal("My Cool App"))
			Expect(sender.Links.Self.Href).To(Equal(fmt.Sprintf("/senders/%s", sender.ID)))
			Expect(sender.Links.CampaignTypes.Href).To(Equal(fmt.Sprintf("/senders/%s/campaign_types", sender.ID)))
			Expect(sender.Links.Campaigns.Href).To(Equal(fmt.Sprintf("/senders/%s/campaigns", sender.ID)))
		})

		By("getting the sender", func() {
			var response senderResponse
			client.Document("sender-get")
			status, err := client.DoTyped("GET", fmt.Sprintf("/senders/%s", senderID), nil, token.Access, &response)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			Expect(response.Name).To(Equal("My Cool App"))
			Expect(response.ID).To(Equal(senderID))
			Expect(response.Links.Self.Href).To(Equal(fmt.Sprintf("/senders/%s", response.ID)))
			Expect(response.Links.CampaignTypes.Href).To(Equal(fmt.Sprintf("/senders/%s/campaign_types", response.ID)))
			Expect(response.Links.Campaigns.Href).To(Equal(fmt.Sprintf("/senders/%s/campaigns", response.ID)))
		})

		By("updating the sender", func() {
			var response senderResponse
			client.Document("sender-update")
			status, err := client.DoTyped("PUT", fmt.Sprintf("/senders/%s", senderID),
				map[string]interface{}{
					"name": "My Not Cool App",
				}, token.Access, &response)

			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			Expect(response.ID).To(Equal(senderID))
			Expect(response.Name).To(Equal("My Not Cool App"))
			Expect(response.Links.Self.Href).To(Equal(fmt.Sprintf("/senders/%s", response.ID)))
			Expect(response.Links.CampaignTypes.Href).To(Equal(fmt.Sprintf("/senders/%s/campaign_types", response.ID)))
			Expect(response.Links.Campaigns.Href).To(Equal(fmt.Sprintf("/senders/%s/campaigns", response.ID)))
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
			client.Document("sender-delete")
			status, err := client.DoTyped("DELETE", fmt.Sprintf("/senders/%s", senderID), nil, token.Access, nil)
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

					status, response, err := client.Do("PUT", fmt.Sprintf("/senders/%s", senderID), map[string]interface{}{
						"name": "My Cool App",
					}, otherToken.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusNotFound))
					Expect(response["errors"]).To(ContainElement(fmt.Sprintf("Sender with id %q could not be found", senderID)))
				})
			})

			It("returns a 422 when updating the sender name to a sender name that exists", func() {
				var senderID string
				By("creating a sender named foo", func() {
					status, _, err := client.Do("POST", "/senders", map[string]interface{}{
						"name": "foo",
					}, token.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusCreated))
				})

				By("creating a sender named bar", func() {
					status, response, err := client.Do("POST", "/senders", map[string]interface{}{
						"name": "bar",
					}, token.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusCreated))
					senderID = response["id"].(string)
				})

				By("updating bar to foo", func() {
					status, response, err := client.Do("PUT", fmt.Sprintf("/senders/%s", senderID), map[string]interface{}{
						"name": "foo",
					}, token.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(422))
					Expect(response["errors"]).To(ContainElement("Sender with name \"foo\" already exists"))
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
