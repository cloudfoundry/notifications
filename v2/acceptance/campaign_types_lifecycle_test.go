package v2

import (
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v2/acceptance/support"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Campaign types lifecycle", func() {
	var (
		client   *support.Client
		token    uaa.Token
		senderID string
	)

	BeforeEach(func() {
		client = support.NewClient(support.Config{
			Host:  Servers.Notifications.URL(),
			Trace: Trace,
		})
		token = GetClientTokenFor("my-client")

		status, response, err := client.Do("POST", "/senders", map[string]interface{}{
			"name": "my-sender",
		}, token.Access)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusCreated))

		senderID = response["id"].(string)
	})

	It("can create, update, show, and delete a new campaign type", func() {
		var campaignTypeID string

		By("creating a campaign type", func() {
			status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
				"name":        "some-campaign-type",
				"description": "a great campaign type",
			}, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			campaignTypeID = response["id"].(string)

			Expect(response["name"]).To(Equal("some-campaign-type"))
			Expect(response["description"]).To(Equal("a great campaign type"))
			Expect(response["critical"]).To(BeFalse())
			Expect(response["template_id"]).To(BeEmpty())
		})

		By("showing the newly created campaign type", func() {
			status, response, err := client.Do("GET", fmt.Sprintf("/senders/%s/campaign_types/%s", senderID, campaignTypeID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			Expect(response["id"]).To(Equal(campaignTypeID))
			Expect(response["name"]).To(Equal("some-campaign-type"))
			Expect(response["description"]).To(Equal("a great campaign type"))
			Expect(response["critical"]).To(BeFalse())
			Expect(response["template_id"]).To(BeEmpty())
		})

		By("creating it again with the same name", func() {
			status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
				"name":        "some-campaign-type",
				"description": "another great campaign type",
			}, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			Expect(response["id"]).To(Equal(campaignTypeID))
			Expect(response["name"]).To(Equal("some-campaign-type"))
			Expect(response["description"]).To(Equal("a great campaign type"))
			Expect(response["critical"]).To(BeFalse())
			Expect(response["template_id"]).To(BeEmpty())
		})

		By("updating it with different information", func() {
			status, response, err := client.Do("PUT", fmt.Sprintf("/senders/%s/campaign_types/%s", senderID, campaignTypeID), map[string]interface{}{
				"name":        "updated-campaign-type",
				"description": "still the same great campaign type",
				"critical":    true,
			}, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			Expect(response["name"]).To(Equal("updated-campaign-type"))
			Expect(response["description"]).To(Equal("still the same great campaign type"))
			Expect(response["critical"]).To(BeTrue())
			Expect(response["template_id"]).To(BeEmpty())
		})

		By("showing the updated campaign type", func() {
			status, response, err := client.Do("GET", fmt.Sprintf("/senders/%s/campaign_types/%s", senderID, campaignTypeID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))

			Expect(response["id"]).To(Equal(campaignTypeID))
			Expect(response["name"]).To(Equal("updated-campaign-type"))
			Expect(response["description"]).To(Equal("still the same great campaign type"))
			Expect(response["critical"]).To(BeTrue())
			Expect(response["template_id"]).To(BeEmpty())
		})

		By("deleting the campaign type", func() {
			status, _, err := client.Do("DELETE", fmt.Sprintf("/senders/%s/campaign_types/%s", senderID, campaignTypeID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))

			status, _, err = client.Do("GET", fmt.Sprintf("/senders/%s/campaign_types/%s", senderID, campaignTypeID), nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNotFound))
		})
	})

	Context("failure cases", func() {
		PIt("returns a 403 when the client does not have access to create a critical campaign type", func() {})

		It("returns a 404 when the campaign type cannot be retrieved", func() {
			var senderID string

			By("creating a sender", func() {
				status, response, err := client.Do("POST", "/senders", map[string]interface{}{
					"name": "My Sender",
				}, token.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				senderID = response["id"].(string)
			})

			By("attempting to retrieve a non-existent campaign type", func() {
				status, response, err := client.Do("GET", fmt.Sprintf("/senders/%s/campaign_types/missing-campaign-type-id", senderID), nil, token.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement(`Campaign type with id "missing-campaign-type-id" could not be found`))
			})
		})

		It("returns a 404 when the campaign type is associated with a different sender", func() {
			var (
				campaignTypeID string
				otherSenderID  string
			)

			By("creating a campaign type belonging to 'my-sender'", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
					"name":        "some-campaign-type",
					"description": "a great campaign type",
				}, token.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				campaignTypeID = response["id"].(string)
			})

			By("creating a sender that is not 'my-sender'", func() {
				status, response, err := client.Do("POST", "/senders", map[string]interface{}{
					"name": "some-other-sender",
				}, token.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				otherSenderID = response["id"].(string)
			})

			By("verifying that you cannot get a campaign type belonging to a different sender", func() {
				status, response, err := client.Do("GET", fmt.Sprintf("/senders/%s/campaign_types/%s", otherSenderID, campaignTypeID), nil, token.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement(fmt.Sprintf("Campaign type with id %q could not be found", campaignTypeID)))
			})
		})

		It("returns a 404 when the sender cannot be retrieved", func() {
			status, response, err := client.Do("GET", "/senders/missing-sender-id/campaign_types/missing-campaign-type-id", nil, token.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNotFound))
			Expect(response["errors"]).To(ContainElement(`Sender with id "missing-sender-id" could not be found`))
		})

		It("returns a 404 when attempting to create a campaign type for a sender that belongs to a different client", func() {
			var senderID string

			By("creating a sender as one client", func() {
				status, response, err := client.Do("POST", "/senders", map[string]interface{}{
					"name": "My Sender",
				}, token.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				senderID = response["id"].(string)
			})

			By("attempting to create a campaign type as a different client", func() {
				token := GetClientTokenFor("other-client")
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
					"name":        "some-campaign-type",
					"description": "a great campaign type",
				}, token.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement(fmt.Sprintf("Sender with id %q could not be found", senderID)))
			})
		})

		It("returns a 404 when attempting to retrieve a campaign type for a sender that belongs to a different client", func() {
			var (
				senderID       string
				campaignTypeID string
			)

			By("creating a sender as one client", func() {
				status, response, err := client.Do("POST", "/senders", map[string]interface{}{
					"name": "My Sender",
				}, token.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				senderID = response["id"].(string)
			})

			By("creating a campaign type belonging to 'My sender'", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
					"name":        "some-campaign-type",
					"description": "a great campaign type",
				}, token.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				campaignTypeID = response["id"].(string)
			})

			By("attempting to retrieve the campaign type as a different client", func() {
				token := GetClientTokenFor("other-client")
				status, response, err := client.Do("GET", fmt.Sprintf("/senders/%s/campaign_types/%s", senderID, campaignTypeID), nil, token.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement(fmt.Sprintf("Sender with id %q could not be found", senderID)))
			})
		})
	})
})
