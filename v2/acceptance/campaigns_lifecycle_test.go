package acceptance

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/chrj/smtpd"

	"github.com/cloudfoundry-incubator/notifications/v2/acceptance/support"
	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type campaignResponse struct {
	ID             string
	CampaignTypeID string              `json:"campaign_type_id"`
	SendTo         map[string][]string `json:"send_to"`
	Text           string
	Subject        string
	TemplateID     string `json:"template_id"`
	ReplyTo        string `json:"reply_to"`
	Links          struct {
		Self struct {
			Href string
		}
		Template struct {
			Href string
		}
		CampaignType struct {
			Href string
		} `json:"campaign_type"`
		Status struct {
			Href string
		}
	} `json:"_links"`
}

var _ = Describe("Campaign Lifecycle", func() {
	var (
		client   *support.Client
		token    string
		senderID string
	)

	BeforeEach(func() {
		client = support.NewClient(support.Config{
			Host:              Servers.Notifications.URL(),
			Trace:             Trace,
			RoundTripRecorder: roundtripRecorder,
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

	Context("retrieving a campaign", func() {
		var (
			response       campaignResponse
			user           warrant.User
			campaignTypeID string
			templateID     string
			campaignID     string
		)

		BeforeEach(func() {
			var err error
			user, err = warrantClient.Users.Create("user-111", "user-111@example.com", adminToken)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			err := warrantClient.Users.Delete(user.ID, adminToken)
			Expect(err).NotTo(HaveOccurred())
		})

		It("sends a campaign to many audiences and retrieves the campaign details", func() {
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
				var response campaignResponse
				client.Document("campaign-create")
				status, err := client.DoTyped("POST", fmt.Sprintf("/senders/%s/campaigns", senderID), map[string]interface{}{
					"send_to": map[string][]string{
						"emails": {"test@example.com"},
						"spaces": {"space-123"},
						"orgs":   {"org-123"},
						"users":  {user.ID},
					},
					"campaign_type_id": campaignTypeID,
					"text":             "campaign body",
					"subject":          "campaign subject",
					"template_id":      templateID,
					"reply_to":         "no-reply@example.com",
				}, token, &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusAccepted))

				campaignID = response.ID

				Expect(response.ID).To(Equal(campaignID))
				Expect(response.SendTo).To(HaveKeyWithValue("emails", []string{"test@example.com"}))
				Expect(response.SendTo).To(HaveKeyWithValue("spaces", []string{"space-123"}))
				Expect(response.SendTo).To(HaveKeyWithValue("orgs", []string{"org-123"}))
				Expect(response.SendTo).To(HaveKeyWithValue("users", []string{user.ID}))
				Expect(response.CampaignTypeID).To(Equal(campaignTypeID))
				Expect(response.Text).To(Equal("campaign body"))
				Expect(response.Subject).To(Equal("campaign subject"))
				Expect(response.TemplateID).To(Equal(templateID))
				Expect(response.ReplyTo).To(Equal("no-reply@example.com"))
				Expect(response.Links.Self.Href).To(Equal(fmt.Sprintf("/campaigns/%s", campaignID)))
				Expect(response.Links.Template.Href).To(Equal(fmt.Sprintf("/templates/%s", templateID)))
				Expect(response.Links.CampaignType.Href).To(Equal(fmt.Sprintf("/campaign_types/%s", campaignTypeID)))
				Expect(response.Links.Status.Href).To(Equal(fmt.Sprintf("/campaigns/%s/status", campaignID)))
			})

			By("retrieving the campaign details", func() {
				client.Document("campaign-get")
				status, err := client.DoTyped("GET", fmt.Sprintf("/campaigns/%s", campaignID), nil, token, &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
				Expect(response.ID).To(Equal(campaignID))
				Expect(response.SendTo).To(HaveKeyWithValue("emails", []string{"test@example.com"}))
				Expect(response.SendTo).To(HaveKeyWithValue("spaces", []string{"space-123"}))
				Expect(response.SendTo).To(HaveKeyWithValue("orgs", []string{"org-123"}))
				Expect(response.SendTo).To(HaveKeyWithValue("users", []string{user.ID}))
				Expect(response.CampaignTypeID).To(Equal(campaignTypeID))
				Expect(response.Text).To(Equal("campaign body"))
				Expect(response.Subject).To(Equal("campaign subject"))
				Expect(response.TemplateID).To(Equal(templateID))
				Expect(response.ReplyTo).To(Equal("no-reply@example.com"))
				Expect(response.Links.Self.Href).To(Equal(fmt.Sprintf("/campaigns/%s", campaignID)))
				Expect(response.Links.Template.Href).To(Equal(fmt.Sprintf("/templates/%s", templateID)))
				Expect(response.Links.CampaignType.Href).To(Equal(fmt.Sprintf("/campaign_types/%s", campaignTypeID)))
				Expect(response.Links.Status.Href).To(Equal(fmt.Sprintf("/campaigns/%s/status", campaignID)))
			})

			By("seeing that the mail was delivered", func() {
				Eventually(func() []smtpd.Envelope {
					return Servers.SMTP.Deliveries
				}, "5s").Should(HaveLen(3))

				var recipients []string
				for _, delivery := range Servers.SMTP.Deliveries {
					Expect(delivery.Recipients).To(HaveLen(1))
					recipients = append(recipients, delivery.Recipients[0])
				}

				Expect(recipients).To(ConsistOf([]string{
					"test@example.com",
					"user-456@example.com",
					"user-111@example.com",
				}))
			})
		})

		Context("when the campaign uses the default template", func() {
			var campaignTypeID, campaignID string

			It("sends a campaign using the default template", func() {
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
							"emails": {"test@example.com"},
						},
						"campaign_type_id": campaignTypeID,
						"text":             "campaign body",
						"subject":          "campaign subject",
					}, token)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusAccepted))
					Expect(response["id"]).NotTo(BeEmpty())

					campaignID = response["id"].(string)
				})

				By("retrieving the campaign details", func() {
					status, response, err := client.Do("GET", fmt.Sprintf("/campaigns/%s", campaignID), nil, token)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusOK))
					Expect(response["id"]).To(Equal(campaignID))
					Expect(response["send_to"]).To(HaveKeyWithValue("emails", []interface{}{"test@example.com"}))
					Expect(response["campaign_type_id"]).To(Equal(campaignTypeID))
					Expect(response["text"]).To(Equal("campaign body"))
					Expect(response["subject"]).To(Equal("campaign subject"))
					Expect(response["template_id"]).To(Equal("default"))
				})
			})
		})

		Context("failure cases", func() {
			var campaignTypeID, templateID, campaignID string

			BeforeEach(func() {
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

				By("creating a campaign type", func() {
					status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
						"name":        "some-campaign-type-name",
						"description": "acceptance campaign type",
					}, token)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusCreated))

					campaignTypeID = response["id"].(string)
				})

				By("sending a campaign", func() {
					status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaigns", senderID), map[string]interface{}{
						"send_to": map[string][]string{
							"emails": {"test@example.com"},
						},
						"campaign_type_id": campaignTypeID,
						"text":             "campaign body",
						"subject":          "campaign subject",
						"template_id":      templateID,
					}, token)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusAccepted))
					Expect(response["id"]).NotTo(BeEmpty())

					campaignID = response["id"].(string)
				})
			})

			It("verifies that the sender exists", func() {
				By("deleting the sender", func() {
					status, _, err := client.Do("DELETE", fmt.Sprintf("/senders/%s", senderID), map[string]interface{}{
						"name": "my-sender",
					}, token)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusNoContent))
				})

				By("attempting to retrieve the campaign details", func() {
					status, response, err := client.Do("GET", fmt.Sprintf("/campaigns/%s", campaignID), nil, token)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusNotFound))
					Expect(response["errors"]).To(ContainElement(fmt.Sprintf("Sender with id %q could not be found", senderID)))
				})
			})

			It("verifies that the campaign belongs to the authenticated client", func() {
				otherClientToken, err := GetClientTokenWithScopes("notifications.write")
				Expect(err).NotTo(HaveOccurred())

				By("attempting to retrieve the campaign details with a different token", func() {
					status, response, err := client.Do("GET", fmt.Sprintf("/campaigns/%s", campaignID), nil, otherClientToken)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusNotFound))
					Expect(response["errors"]).To(ContainElement(fmt.Sprintf("Campaign with id %q could not be found", campaignID)))
				})
			})

			It("verifies that the campaign exists", func() {
				By("attempting to retrieve an unknown campaign", func() {
					status, response, err := client.Do("GET", "/campaigns/unknown-campaign-id", nil, token)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusNotFound))
					Expect(response["errors"]).To(ContainElement("Campaign with id \"unknown-campaign-id\" could not be found"))
				})
			})

			It("verifies that the notifications.write scope is required", func() {
				otherClientToken, err := GetClientTokenWithScopes()
				Expect(err).NotTo(HaveOccurred())

				By("attempting to retrieve the campaign details with a different token", func() {
					status, response, err := client.Do("GET", fmt.Sprintf("/campaigns/%s", campaignID), nil, otherClientToken)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusForbidden))
					Expect(response["errors"]).To(ContainElement("You are not authorized to perform the requested action"))
				})
			})
		})
	})

	Context("retrieving a campaign status", func() {
		var campaignTypeID, templateID, campaignID string

		BeforeEach(func() {
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

			By("creating a campaign type", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
					"name":        "some-campaign-type-name",
					"description": "acceptance campaign type",
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				campaignTypeID = response["id"].(string)
			})
		})

		It("sends a campaign to an email and retrieves the campaign status once the campaign is complete", func() {
			By("sending the campaign", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaigns", senderID), map[string]interface{}{
					"send_to": map[string][]string{
						"emails": {"test@example.com"},
					},
					"campaign_type_id": campaignTypeID,
					"text":             "campaign body",
					"subject":          "campaign subject",
					"template_id":      templateID,
				}, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusAccepted))
				Expect(response["id"]).NotTo(BeEmpty())

				campaignID = response["id"].(string)
			})

			By("retrieving the campaign status", func() {
				Eventually(func() (interface{}, error) {
					_, response, err := client.Do("GET", fmt.Sprintf("/campaigns/%s/status", campaignID), nil, token)
					return response["status"], err
				}, "5s").Should(Equal("completed"))

				var response struct {
					ID             string
					Status         string
					TotalMessages  int       `json:"total_messages"`
					SentMessages   int       `json:"sent_messages"`
					RetryMessages  int       `json:"retry_messages"`
					FailedMessages int       `json:"failed_messages"`
					QueuedMessages int       `json:"queued_messages"`
					StartTime      time.Time `json:"start_time"`
					CompletedTime  time.Time `json:"completed_time"`
					Links          struct {
						Self struct {
							Href string
						}
						Campaign struct {
							Href string
						}
					} `json:"_links"`
				}

				client.Document("campaign-status")
				status, err := client.DoTyped("GET", fmt.Sprintf("/campaigns/%s/status", campaignID), nil, token, &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusOK))
				Expect(response.ID).To(Equal(campaignID))
				Expect(response.Status).To(Equal("completed"))
				Expect(response.TotalMessages).To(Equal(1))
				Expect(response.SentMessages).To(Equal(1))
				Expect(response.RetryMessages).To(Equal(0))
				Expect(response.FailedMessages).To(Equal(0))
				Expect(response.QueuedMessages).To(Equal(0))
				Expect(response.StartTime).To(BeTemporally("~", time.Now(), 10*time.Second))
				Expect(response.CompletedTime).To(BeTemporally("~", time.Now(), 10*time.Second))
				Expect(response.Links.Self.Href).To(Equal(fmt.Sprintf("/campaigns/%s/status", campaignID)))
				Expect(response.Links.Campaign.Href).To(Equal(fmt.Sprintf("/campaigns/%s", campaignID)))
			})
		})

		It("returns a 404 with an error message when the campaign cannot be found", func() {
			By("retrieving the campaign status", func() {
				status, response, err := client.Do("GET", "/campaigns/missing-campaign-id/status", nil, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement("Campaign with id \"missing-campaign-id\" could not be found"))
			})
		})

		Context("when the SMTP server is erroring", func() {
			BeforeEach(func() {
				Servers.SMTP.HandlerCall.Returns.Error = errors.New("some error")
			})

			AfterEach(func() {
				Servers.Notifications.ResetDatabase()
			})

			It("sends a campaign to an email and retrieves the campaign status before the campaign completes", func() {
				By("sending the campaign", func() {
					status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaigns", senderID), map[string]interface{}{
						"send_to": map[string][]string{
							"emails": {"test@example.com"},
						},
						"campaign_type_id": campaignTypeID,
						"text":             "campaign body",
						"subject":          "campaign subject",
						"template_id":      templateID,
					}, token)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusAccepted))
					Expect(response["id"]).NotTo(BeEmpty())

					campaignID = response["id"].(string)
				})

				By("waiting for the status to indicate retrying", func() {
					Eventually(func() (interface{}, error) {
						_, response, err := client.Do("GET", fmt.Sprintf("/campaigns/%s/status", campaignID), nil, token)
						return response["retry_messages"], err
					}, "5s").Should(Equal(float64(1)))
				})

				By("retrieving the campaign status", func() {
					status, response, err := client.Do("GET", fmt.Sprintf("/campaigns/%s/status", campaignID), nil, token)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusOK))

					Expect(response["id"]).To(Equal(campaignID))
					Expect(response["status"]).To(Equal("sending"))
					Expect(response["total_messages"]).To(Equal(float64(1)))
					Expect(response["sent_messages"]).To(Equal(float64(0)))
					Expect(response["retry_messages"]).To(Equal(float64(1)))
					Expect(response["failed_messages"]).To(Equal(float64(0)))

					startTime, err := time.Parse(time.RFC3339, response["start_time"].(string))
					Expect(err).NotTo(HaveOccurred())
					Expect(startTime).To(BeTemporally("~", time.Now(), 10*time.Second))

					Expect(response["completed_time"]).To(BeNil())
				})
			})
		})

		Context("when there are jobs queueing up", func() {
			BeforeEach(func() {
				Servers.SMTP.HandlerCall.Callback = func() {
					time.Sleep(500 * time.Millisecond)
				}
			})

			It("sends a campaign to an email and retrieves the campaign status before the campaign completes", func() {
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

					campaignID = response["id"].(string)
				})

				By("waiting for the status to indicate queueing", func() {
					Eventually(func() (interface{}, error) {
						_, response, err := client.Do("GET", fmt.Sprintf("/campaigns/%s/status", campaignID), nil, token)
						return response["queued_messages"], err
					}, "5s").Should(BeNumerically(">", 0))
				})
			})
		})
	})
})
