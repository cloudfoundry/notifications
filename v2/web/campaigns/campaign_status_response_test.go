package campaigns_test

import (
	"encoding/json"
	"time"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/campaigns"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CampaignStatusResponse", func() {
	var startTime, completedTime time.Time
	BeforeEach(func() {
		var err error
		startTime, err = time.Parse(time.RFC3339, "2009-12-11T10:21:45Z")
		Expect(err).NotTo(HaveOccurred())

		completedTime, err = time.Parse(time.RFC3339, "2009-12-11T10:21:59Z")
		Expect(err).NotTo(HaveOccurred())
	})

	It("provides a JSON representation of a campaign status resource", func() {
		campaignStatus := collections.CampaignStatus{
			CampaignID:            "some-campaign-id",
			Status:                "sending",
			TotalMessages:         5,
			SentMessages:          1,
			RetryMessages:         1,
			FailedMessages:        1,
			QueuedMessages:        1,
			UndeliverableMessages: 1,
			StartTime:             startTime,
			CompletedTime:         nil,
		}

		response := campaigns.NewCampaignStatusResponse(campaignStatus)
		Expect(response).To(Equal(campaigns.CampaignStatusResponse{
			CampaignID:            "some-campaign-id",
			Status:                "sending",
			TotalMessages:         5,
			SentMessages:          1,
			RetryMessages:         1,
			FailedMessages:        1,
			QueuedMessages:        1,
			UndeliverableMessages: 1,
			StartTime:             startTime,
			CompletedTime:         nil,
			Links: campaigns.CampaignStatusResponseLinks{
				Self:     campaigns.Link{"/campaigns/some-campaign-id/status"},
				Campaign: campaigns.Link{"/campaigns/some-campaign-id"},
			},
		}))
	})

	It("can marshal into JSON", func() {
		campaignStatus := collections.CampaignStatus{
			CampaignID:            "some-campaign-id",
			Status:                "completed",
			TotalMessages:         5,
			SentMessages:          2,
			RetryMessages:         0,
			FailedMessages:        1,
			QueuedMessages:        0,
			UndeliverableMessages: 2,
			StartTime:             startTime,
			CompletedTime:         &completedTime,
		}

		output, err := json.Marshal(campaigns.NewCampaignStatusResponse(campaignStatus))
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(MatchJSON(`{
			"id": "some-campaign-id",
			"status": "completed",
			"total_messages": 5,
			"sent_messages": 2,
			"retry_messages": 0,
			"failed_messages": 1,
			"queued_messages": 0,
			"undeliverable_messages": 2,
			"start_time": "2009-12-11T10:21:45Z",
			"completed_time": "2009-12-11T10:21:59Z",
			"_links": {
				"self": {
					"href": "/campaigns/some-campaign-id/status"
				},
				"campaign": {
					"href": "/campaigns/some-campaign-id"
				}
			}
		}`))
	})
})
