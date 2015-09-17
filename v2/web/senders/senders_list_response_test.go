package senders_test

import (
	"encoding/json"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/senders"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SendersListResponse", func() {
	var response senders.SendersListResponse

	BeforeEach(func() {
		response = senders.NewSendersListResponse([]collections.Sender{
			{
				ID:   "some-sender-id",
				Name: "some-sender",
			},
			{
				ID:   "another-sender-id",
				Name: "another-sender",
			},
		})
	})

	It("provides a JSON representation of a list of sender resources", func() {
		Expect(response).To(Equal(senders.SendersListResponse{
			Senders: []senders.SenderResponse{
				{
					ID:   "some-sender-id",
					Name: "some-sender",
					Links: senders.SenderResponseLinks{
						Self:          senders.Link{"/senders/some-sender-id"},
						CampaignTypes: senders.Link{"/senders/some-sender-id/campaign_types"},
						Campaigns:     senders.Link{"/senders/some-sender-id/campaigns"},
					},
				},
				{
					ID:   "another-sender-id",
					Name: "another-sender",
					Links: senders.SenderResponseLinks{
						Self:          senders.Link{"/senders/another-sender-id"},
						CampaignTypes: senders.Link{"/senders/another-sender-id/campaign_types"},
						Campaigns:     senders.Link{"/senders/another-sender-id/campaigns"},
					},
				},
			},
			Links: senders.SendersListResponseLinks{
				Self: senders.Link{"/senders"},
			},
		}))
	})

	It("can marshal into JSON", func() {
		output, err := json.Marshal(response)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(MatchJSON(`{
			"senders": [
				{
					"id":   "some-sender-id",
					"name": "some-sender",
					"_links": {
						"self": {
							"href": "/senders/some-sender-id"
						},
						"campaign_types": {
							"href": "/senders/some-sender-id/campaign_types"
						},
						"campaigns": {
							"href": "/senders/some-sender-id/campaigns"
						}
					}
				},
				{
					"id":   "another-sender-id",
					"name": "another-sender",
					"_links": {
						"self": {
							"href": "/senders/another-sender-id"
						},
						"campaign_types": {
							"href": "/senders/another-sender-id/campaign_types"
						},
						"campaigns": {
							"href": "/senders/another-sender-id/campaigns"
						}
					}
				}
			],
			"_links": {
				"self": {
					"href": "/senders"
				}
			}
		}`))
	})
})
