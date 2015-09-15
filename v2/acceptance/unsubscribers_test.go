package acceptance

import (
	"fmt"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v2/acceptance/support"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unsubscribers", func() {
	var (
		client      *support.Client
		clientToken uaa.Token
		userToken   uaa.Token

		senderID, userGUID, campaignTypeID, templateID, campaignID string
	)

	BeforeEach(func() {
		client = support.NewClient(support.Config{
			Host:  Servers.Notifications.URL(),
			Trace: Trace,
		})
		clientToken = GetClientTokenFor("my-client")

		status, response, err := client.Do("POST", "/senders", map[string]interface{}{
			"name": "my-sender",
		}, clientToken.Access)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusCreated))

		userGUID = "user-123"
		userToken = GetUserTokenFor(userGUID)
		senderID = response["id"].(string)

		By("creating a template", func() {
			status, response, err := client.Do("POST", "/templates", map[string]interface{}{
				"name":    "Acceptance Template",
				"text":    "campaign template {{.Text}}",
				"html":    "{{.HTML}}",
				"subject": "{{.Subject}}",
			}, clientToken.Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusCreated))

			templateID = response["id"].(string)
		})
	})

	Context("managing subscription with a client token", func() {
		It("delivers or not based on the unsubscribe state", func() {
			By("creating a campaign type", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
					"name":        "some-campaign-type-name",
					"description": "acceptance campaign type",
					"template_id": templateID,
				}, clientToken.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				campaignTypeID = response["id"].(string)
			})

			By("unsubscribing from the campaign type", func() {
				path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, campaignTypeID, userGUID)
				status, _, err := client.Do("PUT", path, nil, clientToken.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNoContent))
			})

			By("sending the campaign", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaigns", senderID), map[string]interface{}{
					"send_to": map[string]interface{}{
						"user": userGUID,
					},
					"campaign_type_id": campaignTypeID,
					"text":             "campaign body",
					"subject":          "campaign subject",
					"template_id":      templateID,
				}, clientToken.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusAccepted))
				Expect(response["campaign_id"]).NotTo(BeEmpty())

				campaignID = response["campaign_id"].(string)
			})

			By("waiting for the email to arrive", func() {
				Eventually(func() (interface{}, error) {
					_, response, err := client.Do("GET", fmt.Sprintf("/senders/%s/campaigns/%s/status", senderID, campaignID), nil, clientToken.Access)
					return response["status"], err
				}).Should(Equal("completed"))

				Expect(Servers.SMTP.Deliveries).To(HaveLen(0))
			})

			By("deleting the unsubscribe", func() {
				path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, campaignTypeID, userGUID)
				status, _, err := client.Do("DELETE", path, nil, clientToken.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNoContent))
			})

			var secondCampaignID string

			By("sending another campaign", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaigns", senderID), map[string]interface{}{
					"send_to": map[string]interface{}{
						"user": userGUID,
					},
					"campaign_type_id": campaignTypeID,
					"text":             "campaign body",
					"subject":          "campaign subject",
					"template_id":      templateID,
				}, clientToken.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusAccepted))
				Expect(response["campaign_id"]).NotTo(BeEmpty())

				secondCampaignID = response["campaign_id"].(string)
			})

			By("confirming that the email is received", func() {
				Eventually(func() (interface{}, error) {
					_, response, err := client.Do("GET", fmt.Sprintf("/senders/%s/campaigns/%s/status", senderID, secondCampaignID), nil, clientToken.Access)
					return response["status"], err
				}).Should(Equal("completed"))

				Expect(Servers.SMTP.Deliveries).To(HaveLen(1))

				Expect(Servers.SMTP.Deliveries[0].Recipients).To(ConsistOf([]string{
					"user-123@example.com",
				}))
			})
		})

		Context("when attempting to unsubscribe from a critical notification", func() {
			It("returns a 403 status code and reports an error message as JSON", func() {
				By("creating a campaign type", func() {
					status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
						"name":        "some-campaign-type-name",
						"description": "acceptance campaign type",
						"template_id": templateID,
						"critical":    true,
					}, clientToken.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusCreated))

					campaignTypeID = response["id"].(string)
				})

				By("unsubscribing from the campaign type", func() {
					path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, campaignTypeID, userGUID)
					status, response, err := client.Do("PUT", path, nil, clientToken.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusForbidden))
					Expect(response["errors"]).To(ContainElement(fmt.Sprintf("Campaign type %q cannot be unsubscribed from", campaignTypeID)))
				})
			})
		})

		Context("when the API client lacks the required scopes", func() {
			It("returns a 403 status code for PUTting an unsubscribe", func() {
				clientToken = GetClientTokenFor("non-admin-client")

				path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, "some-campaign-type-id", userGUID)
				status, response, err := client.Do("PUT", path, nil, clientToken.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusForbidden))
				Expect(response["errors"]).To(ContainElement("You are not authorized to perform the requested action"))
			})

			It("returns a 403 status code for DELETEing an unsubscribe", func() {
				clientToken = GetClientTokenFor("non-admin-client")

				path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, "some-campaign-type-id", userGUID)
				status, response, err := client.Do("DELETE", path, nil, clientToken.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusForbidden))
				Expect(response["errors"]).To(ContainElement("You are not authorized to perform the requested action"))
			})
		})

		Context("when attempting to manage subscriptions for a non-existent user", func() {
			It("returns a 404 status code and reports the error message as JSON", func() {
				var campaignTypeID, templateID string
				By("creating a campaign type", func() {
					status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
						"name":        "some-campaign-type-name",
						"description": "acceptance campaign type",
						"template_id": templateID,
					}, clientToken.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusCreated))

					campaignTypeID = response["id"].(string)
				})

				By("unsubscribing from the campaign type", func() {
					path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, campaignTypeID, "not-a-user")
					status, response, err := client.Do("PUT", path, nil, clientToken.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusNotFound))
					Expect(response["errors"]).To(ContainElement("User \"not-a-user\" not found"))
				})

				By("removing an unsubscribe from the campaign type", func() {
					path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, campaignTypeID, "not-a-user")
					status, response, err := client.Do("DELETE", path, nil, clientToken.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusNotFound))
					Expect(response["errors"]).To(ContainElement("User \"not-a-user\" not found"))
				})
			})
		})

		Context("when attempting to manage subscriptions for a non-existent campaign type", func() {
			It("returns a 404 status code and reports the error message as JSON for PUTs", func() {
				path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, "not-a-campaign-type", userGUID)
				status, response, err := client.Do("PUT", path, nil, clientToken.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement("Campaign type with id \"not-a-campaign-type\" could not be found"))
			})
			It("returns a 404 status code and reports the error message as JSON for DELETEs", func() {
				path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, "not-a-campaign-type", userGUID)
				status, response, err := client.Do("DELETE", path, nil, clientToken.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement("Campaign type with id \"not-a-campaign-type\" could not be found"))
			})
		})
	})

	Context("managing subscription with a user token", func() {
		It("delivers or not based on the unsubscribe state", func() {
			By("creating a campaign type using the client token", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
					"name":        "some-campaign-type-name",
					"description": "acceptance campaign type",
					"template_id": templateID,
				}, clientToken.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusCreated))

				campaignTypeID = response["id"].(string)
			})

			By("unsubscribing from the campaign type using the user token", func() {
				path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, campaignTypeID, userGUID)
				status, _, err := client.Do("PUT", path, nil, userToken.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNoContent))
			})

			By("sending the campaign with the client token", func() {
				status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaigns", senderID), map[string]interface{}{
					"send_to": map[string]interface{}{
						"user": userGUID,
					},
					"campaign_type_id": campaignTypeID,
					"text":             "campaign body",
					"subject":          "campaign subject",
					"template_id":      templateID,
				}, clientToken.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusAccepted))
				Expect(response["campaign_id"]).NotTo(BeEmpty())

				campaignID = response["campaign_id"].(string)
			})

			By("waiting for the email to arrive", func() {
				Eventually(func() (interface{}, error) {
					_, response, err := client.Do("GET", fmt.Sprintf("/senders/%s/campaigns/%s/status", senderID, campaignID), nil, clientToken.Access)
					return response["status"], err
				}).Should(Equal("completed"))

				Expect(Servers.SMTP.Deliveries).To(HaveLen(0))
			})
		})

		Context("when attempting to unsubscribe for a different user", func() {
			It("returns a 403 status code and reports the error message as JSON", func() {
				By("creating a campaign type using the client token", func() {
					status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
						"name":        "some-campaign-type-name",
						"description": "acceptance campaign type",
						"template_id": templateID,
					}, clientToken.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusCreated))

					campaignTypeID = response["id"].(string)
				})

				By("attempting to unsubscribe from the campaign type using the user token", func() {
					path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, campaignTypeID, "some-other-user")
					status, response, err := client.Do("PUT", path, nil, userToken.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusForbidden))
					Expect(response["errors"]).To(ContainElement("You are not authorized to perform the requested action"))
				})
			})
		})

		Context("when attempting to unsubscribe without any token", func() {
			It("returns a 403 status code and reports the error message as JSON", func() {
				By("creating a campaign type using the client token", func() {
					status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
						"name":        "some-campaign-type-name",
						"description": "acceptance campaign type",
						"template_id": templateID,
					}, clientToken.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusCreated))

					campaignTypeID = response["id"].(string)
				})

				By("attempting to unsubscribe from the campaign type using the user token", func() {
					path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, campaignTypeID, userGUID)
					status, response, err := client.Do("PUT", path, nil, "")
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusUnauthorized))
					Expect(response["errors"]).To(ContainElement("Authorization header is invalid: missing"))
				})
			})
		})

		Context("when the user token lacks the required scopes", func() {
			It("returns a 403 status code for PUTting an unsubscribe", func() {
				unauthorizedUserToken := GetUserTokenFor("unauthorized-user")

				path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, "some-campaign-type-id", userGUID)
				status, response, err := client.Do("PUT", path, nil, unauthorizedUserToken.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusForbidden))
				Expect(response["errors"]).To(ContainElement("You are not authorized to perform the requested action"))
			})
		})

		Context("when attempting to manage subscriptions for a non-existent campaign type", func() {
			It("returns a 404 status code and reports the error message as JSON for PUTs", func() {
				path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, "not-a-campaign-type", userGUID)
				status, response, err := client.Do("PUT", path, nil, userToken.Access)
				Expect(err).NotTo(HaveOccurred())
				Expect(status).To(Equal(http.StatusNotFound))
				Expect(response["errors"]).To(ContainElement("Campaign type with id \"not-a-campaign-type\" could not be found"))
			})
		})

		Context("when attempting to unsubscribe from a critical notification", func() {
			It("returns a 403 status code and reports an error message as JSON", func() {
				By("creating a campaign type", func() {
					status, response, err := client.Do("POST", fmt.Sprintf("/senders/%s/campaign_types", senderID), map[string]interface{}{
						"name":        "some-campaign-type-name",
						"description": "acceptance campaign type",
						"template_id": templateID,
						"critical":    true,
					}, clientToken.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusCreated))

					campaignTypeID = response["id"].(string)
				})

				By("unsubscribing from the campaign type", func() {
					path := fmt.Sprintf("/senders/%s/campaign_types/%s/unsubscribers/%s", senderID, campaignTypeID, userGUID)
					status, response, err := client.Do("PUT", path, nil, userToken.Access)
					Expect(err).NotTo(HaveOccurred())
					Expect(status).To(Equal(http.StatusForbidden))
					Expect(response["errors"]).To(ContainElement(fmt.Sprintf("Campaign type %q cannot be unsubscribed from", campaignTypeID)))
				})
			})
		})
	})
})
