package campaigns

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type CampaignStatusResponseLinks struct {
	Self     Link `json:"self"`
	Campaign Link `json:"campaign"`
}

type CampaignStatusResponse struct {
	CampaignID            string                      `json:"id"`
	Status                string                      `json:"status"`
	TotalMessages         int                         `json:"total_messages"`
	SentMessages          int                         `json:"sent_messages"`
	RetryMessages         int                         `json:"retry_messages"`
	FailedMessages        int                         `json:"failed_messages"`
	QueuedMessages        int                         `json:"queued_messages"`
	UndeliverableMessages int                         `json:"undeliverable_messages"`
	StartTime             time.Time                   `json:"start_time"`
	CompletedTime         *time.Time                  `json:"completed_time"`
	Links                 CampaignStatusResponseLinks `json:"_links"`
}

func NewCampaignStatusResponse(status collections.CampaignStatus) CampaignStatusResponse {
	return CampaignStatusResponse{
		CampaignID:            status.CampaignID,
		Status:                status.Status,
		TotalMessages:         status.TotalMessages,
		SentMessages:          status.SentMessages,
		RetryMessages:         status.RetryMessages,
		FailedMessages:        status.FailedMessages,
		QueuedMessages:        status.QueuedMessages,
		UndeliverableMessages: status.UndeliverableMessages,
		StartTime:             status.StartTime,
		CompletedTime:         status.CompletedTime,
		Links: CampaignStatusResponseLinks{
			Self:     Link{fmt.Sprintf("/campaigns/%s/status", status.CampaignID)},
			Campaign: Link{fmt.Sprintf("/campaigns/%s", status.CampaignID)},
		},
	}
}
