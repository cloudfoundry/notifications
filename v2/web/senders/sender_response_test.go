package senders_test

import (
	"encoding/json"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/senders"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SenderResponse", func() {
	It("provides a JSON representation of a sender resource", func() {
		sender := collections.Sender{
			ID:   "some-sender-id",
			Name: "some-sender",
		}

		response := senders.NewSenderResponse(sender)
		Expect(response).To(Equal(senders.SenderResponse{
			ID:   "some-sender-id",
			Name: "some-sender",
			Links: senders.SenderResponseLinks{
				Self:          senders.Link{"/senders/some-sender-id"},
				CampaignTypes: senders.Link{"/senders/some-sender-id/campaign_types"},
				Campaigns:     senders.Link{"/senders/some-sender-id/campaigns"},
			},
		}))
	})

	It("can marshal into JSON", func() {
		sender := collections.Sender{
			ID:   "some-sender-id",
			Name: "some-sender",
		}

		output, err := json.Marshal(senders.NewSenderResponse(sender))
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(MatchJSON(`{
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
		}`))
	})
})
