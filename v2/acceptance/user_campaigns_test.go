package acceptance

import (
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/chrj/smtpd"

	"github.com/cloudfoundry-incubator/notifications/v2/acceptance/support"
	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Campaigns", func() {
	var (
		client   *support.Client
		token    string
		senderID string
		user     warrant.User
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

		user, err = warrantClient.Users.Create("user-111", "user-111@example.com", adminToken)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := warrantClient.Users.Delete(user.ID, adminToken)
		Expect(err).NotTo(HaveOccurred())
	})

	It("sends a campaign to a user", func() {
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
					"users": {user.ID},
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
			Expect(delivery.Recipients[0]).To(Equal("user-111@example.com"))

			data := strings.Split(string(delivery.Data), "\n")
			Expect(data).To(ContainElement("campaign template campaign body"))
		})
	})

	It("sends a campaign to multiple users", func() {
		var (
			campaignTypeID         string
			templateID             string
			campaignTypeTemplateID string
			anotherUser            warrant.User
		)

		By("creating another user", func() {
			var err error
			anotherUser, err = warrantClient.Users.Create("user-abc", "user-abc@example.com", adminToken)
			Expect(err).NotTo(HaveOccurred())
		})

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
					"users": {
						user.ID,
						anotherUser.ID,
						user.ID,
					},
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
				"user-111@example.com",
				"user-abc@example.com",
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
						"users": {user.ID},
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

	Context("when the audience key is invalid", func() {
		It("returns a 422 and an error message", func() {
			var campaignTypeID, templateID string

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
						"bananas": {user.ID},
					},
					"campaign_type_id": campaignTypeID,
					"text":             "campaign body",
					"subject":          "campaign subject",
					"template_id":      templateID,
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(422))
				Expect(response["errors"]).To(Equal([]interface{}{"\"bananas\" is not a valid audience"}))
			})
		})
	})

	Context("when the audience is not a list", func() {
		It("returns a 400 and an error message", func() {
			var campaignTypeID, templateID string

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
					"send_to": map[string]string{
						"users": user.ID,
					},
					"campaign_type_id": campaignTypeID,
					"text":             "campaign body",
					"subject":          "campaign subject",
					"template_id":      templateID,
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusBadRequest))
				Expect(response["errors"]).To(Equal([]interface{}{"invalid json body"}))
			})
		})
	})

	Context("when the audience is non-existent", func() {
		It("returns a 404 with an error message", func() {
			var campaignTypeID, templateID string

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
						"users": {"missing-user"},
					},
					"campaign_type_id": campaignTypeID,
					"text":             "campaign body",
					"subject":          "campaign subject",
					"template_id":      templateID,
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(404))
				Expect(response["errors"]).To(Equal([]interface{}{"The user \"missing-user\" cannot be found"}))
			})
		})
	})

	Context("when the sender ID doesn't exist", func() {
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
				status, response, err := client.Do("POST", "/senders/missing-sender-id/campaigns", map[string]interface{}{
					"send_to": map[string][]string{
						"users": {user.ID},
					},
					"campaign_type_id": campaignTypeID,
					"text":             "campaign body",
					"subject":          "campaign subject",
					"template_id":      "missing-template-id",
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement("Sender with id \"missing-sender-id\" could not be found"))
			})
		})
	})

	Context("when the sender ID belongs to a different client", func() {
		It("returns a 404 with an error message", func() {
			var campaignTypeID, differentSenderID string

			By("creating a sender belonging to a different client", func() {
				otherClientToken, err := GetClientTokenWithScopes("notifications.write")
				Expect(err).NotTo(HaveOccurred())

				status, response, err := client.Do("POST", "/senders", map[string]interface{}{
					"name": "my-sender",
				}, otherClientToken)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				differentSenderID = response["id"].(string)
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
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaigns", differentSenderID), map[string]interface{}{
					"send_to": map[string][]string{
						"users": {user.ID},
					},
					"campaign_type_id": campaignTypeID,
					"text":             "campaign body",
					"subject":          "campaign subject",
					"template_id":      "missing-template-id",
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement(fmt.Sprintf("Sender with id %q could not be found", differentSenderID)))
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
						"users": {user.ID},
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
						"users": {user.ID},
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
						"users": {user.ID},
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
				Expect(delivery.Recipients[0]).To(Equal("user-111@example.com"))

				data := strings.Split(string(delivery.Data), "\n")
				Expect(data).To(ContainElement("campaign body campaign type template"))
			})
		})
	})
})
