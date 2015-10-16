package acceptance

import (
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/chrj/smtpd"

	"github.com/cloudfoundry-incubator/notifications/v2/acceptance/support"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Space Campaigns", func() {
	var (
		client   *support.Client
		token    string
		senderID string
	)

	BeforeEach(func() {
		client = support.NewClient(support.Config{
			Host:  Servers.Notifications.URL(),
			Trace: Trace,
		})
		var err error
		token, err = GetClientTokenWithScopes("notifications.write")
		Expect(err).NotTo(HaveOccurred())

		status, response, err := client.Do("POST", "/senders", map[string]interface{}{
			"name": "my-sender",
		}, token)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusCreated))

		senderID = response["id"].(string)
	})

	It("sends a campaign to a space", func() {
		var campaignTypeID, templateID, campaignTypeTemplateID string

		By("creating a template", func() {
			status, response, err := client.Do("POST", "/templates", map[string]interface{}{
				"name":    "Acceptance Template",
				"text":    "campaign template {{.Text}}",
				"html":    "{{.HTML}}",
				"subject": "{{.Subject}}",
			}, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			templateID = response["id"].(string)
		})

		By("creating a campaign type template", func() {
			status, response, err := client.Do("POST", "/templates", map[string]interface{}{
				"name":    "CampaignType Template",
				"text":    "campaign type template {{.Text}}",
				"html":    "{{.HTML}}",
				"subject": "{{.Subject}}",
			}, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			campaignTypeTemplateID = response["id"].(string)
		})

		By("creating a campaign type", func() {
			status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
				"name":        "some-campaign-type-name",
				"description": "acceptance campaign type",
				"template_id": campaignTypeTemplateID,
			}, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			campaignTypeID = response["id"].(string)
		})

		By("sending the campaign", func() {
			status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaigns", senderID), map[string]interface{}{
				"send_to": map[string][]string{
					"spaces": {"space-123"},
				},
				"campaign_type_id": campaignTypeID,
				"text":             "campaign body",
				"subject":          "campaign subject",
				"template_id":      templateID,
			}, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusAccepted))
			Expect(response["id"]).NotTo(BeEmpty())
		})

		By("seeing that the mail was delivered", func() {
			Eventually(func() []smtpd.Envelope {
				return Servers.SMTP.Deliveries
			}, "5s").Should(HaveLen(1))

			delivery := Servers.SMTP.Deliveries[0]

			Expect(delivery.Recipients).To(HaveLen(1))
			Expect(delivery.Recipients).To(ConsistOf([]string{
				"user-456@example.com",
			}))

			data := strings.Split(string(delivery.Data), "\n")
			Expect(data).To(ContainElement("campaign template campaign body"))
		})
	})

	It("sends a campaign to a list of spaces", func() {
		var campaignTypeID, templateID, campaignTypeTemplateID string

		By("creating a template", func() {
			status, response, err := client.Do("POST", "/templates", map[string]interface{}{
				"name":    "Acceptance Template",
				"text":    "campaign template {{.Text}}",
				"html":    "{{.HTML}}",
				"subject": "{{.Subject}}",
			}, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			templateID = response["id"].(string)
		})

		By("creating a campaign type template", func() {
			status, response, err := client.Do("POST", "/templates", map[string]interface{}{
				"name":    "CampaignType Template",
				"text":    "campaign type template {{.Text}}",
				"html":    "{{.HTML}}",
				"subject": "{{.Subject}}",
			}, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			campaignTypeTemplateID = response["id"].(string)
		})

		By("creating a campaign type", func() {
			status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
				"name":        "some-campaign-type-name",
				"description": "acceptance campaign type",
				"template_id": campaignTypeTemplateID,
			}, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			campaignTypeID = response["id"].(string)
		})

		By("sending the campaign", func() {
			status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaigns", senderID), map[string]interface{}{
				"send_to": map[string][]string{
					"spaces": {"space-123", "space-456"},
				},
				"campaign_type_id": campaignTypeID,
				"text":             "campaign body",
				"subject":          "campaign subject",
				"template_id":      templateID,
			}, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusAccepted))
			Expect(response["id"]).NotTo(BeEmpty())
		})

		By("seeing that the mail was delivered", func() {
			Eventually(func() []smtpd.Envelope {
				return Servers.SMTP.Deliveries
			}, "5s").Should(HaveLen(2))

			var recipients []string
			for _, delivery := range Servers.SMTP.Deliveries {
				Expect(delivery.Recipients).To(HaveLen(1))
				recipients = append(recipients, delivery.Recipients[0])
			}

			Expect(recipients).To(ConsistOf([]string{
				"user-456@example.com",
				"user-123@example.com",
			}))
		})
	})

	Context("when lacking the critical scope", func() {
		It("returns a 403 forbidden", func() {
			var campaignTypeID, templateID, campaignTypeTemplateID string

			By("creating a template", func() {
				status, response, err := client.Do("POST", "/templates", map[string]interface{}{
					"name":    "Acceptance Template",
					"text":    "campaign template {{.Text}}",
					"html":    "{{.HTML}}",
					"subject": "{{.Subject}}",
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				templateID = response["id"].(string)
			})

			By("creating a campaign type template", func() {
				status, response, err := client.Do("POST", "/templates", map[string]interface{}{
					"name":    "CampaignType Template",
					"text":    "campaign type template {{.Text}}",
					"html":    "{{.HTML}}",
					"subject": "{{.Subject}}",
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				campaignTypeTemplateID = response["id"].(string)
			})

			By("adding critical_notifications.write scope to client", func() {
				var err error
				token, err = UpdateClientTokenWithDifferentScopes(token, "notifications.write", "critical_notifications.write")
				Expect(err).NotTo(HaveOccurred())
			})

			By("creating a campaign type", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
					"name":        "some-campaign-type-name",
					"description": "acceptance campaign type",
					"template_id": campaignTypeTemplateID,
					"critical":    true,
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				campaignTypeID = response["id"].(string)
			})

			By("removing critical_notifications.write scope from the client", func() {
				var err error
				token, err = UpdateClientTokenWithDifferentScopes(token, "notifications.write")
				Expect(err).NotTo(HaveOccurred())
			})

			By("sending the campaign", func() {
				status, response, err := client.Do("POST", "/senders", map[string]interface{}{
					"name": "my-sender",
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				senderID = response["id"].(string)

				status, response, err = client.Do("POST", fmt.Sprintf("/senders/%s/campaigns", senderID), map[string]interface{}{
					"send_to": map[string][]string{
						"spaces": {"space-123"},
					},
					"campaign_type_id": campaignTypeID,
					"text":             "campaign body",
					"subject":          "campaign subject",
					"template_id":      templateID,
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusForbidden))
				Expect(response["errors"]).To(Equal([]interface{}{"Scope critical_notifications.write is required"}))
			})
		})
	})

	Context("when there is a space that is non-existent", func() {
		It("returns a 202 and sends no emails to that space", func() {
			var campaignID, campaignTypeID, templateID string

			By("creating a template", func() {
				status, response, err := client.Do("POST", "/templates", map[string]interface{}{
					"name":    "Acceptance Template",
					"text":    "{{.Text}}",
					"html":    "{{.HTML}}",
					"subject": "{{.Subject}}",
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				templateID = response["id"].(string)
			})

			By("creating a campaign type", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
					"name":        "some-campaign-type-name",
					"description": "acceptance campaign type",
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				campaignTypeID = response["id"].(string)
			})

			By("sending the campaign", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaigns", senderID), map[string]interface{}{
					"send_to": map[string][]string{
						"spaces": {
							"missing-space",
							"space-123",
						},
					},
					"campaign_type_id": campaignTypeID,
					"text":             "campaign body",
					"subject":          "campaign subject",
					"template_id":      templateID,
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusAccepted))

				campaignID = response["id"].(string)
			})

			By("retrieving the campaign status", func() {
				getDeliveries := func() []string {
					var recipients []string
					for _, delivery := range Servers.SMTP.Deliveries {
						recipients = append(recipients, delivery.Recipients...)
					}
					return recipients
				}

				Eventually(getDeliveries, "5s").Should(ConsistOf([]string{
					"user-456@example.com",
				}))

				status, response, err := client.Do("GET", fmt.Sprintf("/campaigns/%s/status", campaignID), nil, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
				Expect(response["id"]).To(Equal(campaignID))
				Expect(response["sent_messages"]).To(Equal(float64(1)))
			})
		})
	})

	Context("when the template ID doesn't exist", func() {
		It("returns a 404 with an error message", func() {
			var campaignTypeID string

			By("creating a campaign type", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
					"name":        "some-campaign-type-name",
					"description": "acceptance campaign type",
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				campaignTypeID = response["id"].(string)
			})

			By("sending the campaign", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaigns", senderID), map[string]interface{}{
					"send_to": map[string][]string{
						"spaces": {"space-123"},
					},
					"campaign_type_id": campaignTypeID,
					"text":             "campaign body",
					"subject":          "campaign subject",
					"template_id":      "missing-template-id",
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement("Template with id \"missing-template-id\" could not be found"))
			})
		})
	})

	Context("when the campaign type ID does not exist", func() {
		It("returns a 404 with an error message", func() {
			var templateID string

			By("creating a template", func() {
				status, response, err := client.Do("POST", "/templates", map[string]interface{}{
					"name":    "Acceptance Template",
					"text":    "{{.Text}} campaign type template",
					"html":    "{{.HTML}}",
					"subject": "{{.Subject}}",
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				templateID = response["id"].(string)
			})

			By("sending the campaign", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaigns", senderID), map[string]interface{}{
					"send_to": map[string][]string{
						"spaces": {"space-123"},
					},
					"campaign_type_id": "missing-campaign-type-id",
					"text":             "campaign body",
					"subject":          "campaign subject",
					"template_id":      templateID,
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement("Campaign type with id \"missing-campaign-type-id\" could not be found"))
			})
		})
	})

	Context("when omitting the template_id", func() {
		It("uses the template assigned to the campaign type", func() {
			var campaignTypeID, templateID string

			By("creating a template", func() {
				status, response, err := client.Do("POST", "/templates", map[string]interface{}{
					"name":    "Acceptance Template",
					"text":    "{{.Text}} campaign type template",
					"html":    "{{.HTML}}",
					"subject": "{{.Subject}}",
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				templateID = response["id"].(string)
			})

			By("creating a campaign type", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
					"name":        "some-campaign-type-name",
					"description": "acceptance campaign type",
					"template_id": templateID,
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				campaignTypeID = response["id"].(string)
			})

			By("sending the campaign", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaigns", senderID), map[string]interface{}{
					"send_to": map[string][]string{
						"spaces": {"space-123"},
					},
					"campaign_type_id": campaignTypeID,
					"text":             "campaign body",
					"subject":          "campaign subject",
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusAccepted))
				Expect(response["id"]).NotTo(BeEmpty())
			})

			By("seeing that the mail was delivered", func() {
				Eventually(func() []smtpd.Envelope {
					return Servers.SMTP.Deliveries
				}, "5s").Should(HaveLen(1))

				delivery := Servers.SMTP.Deliveries[0]

				Expect(delivery.Recipients).To(HaveLen(1))
				Expect(delivery.Recipients[0]).To(Equal("user-456@example.com"))

				data := strings.Split(string(delivery.Data), "\n")
				Expect(data).To(ContainElement("campaign body campaign type template"))
			})
		})
	})
})
