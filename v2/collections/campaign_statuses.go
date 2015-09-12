package collections

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/go-sql-driver/mysql"
)

const (
	CampaignStatusSending   = "sending"
	CampaignStatusCompleted = "completed"
)

type campaignsGetter interface {
	Get(conn models.ConnectionInterface, campaignID string) (models.Campaign, error)
}

type messageCountGetter interface {
	CountByStatus(conn models.ConnectionInterface, campaignID string) (models.MessageCounts, error)
	MostRecentlyUpdatedByCampaignID(conn models.ConnectionInterface, campaignID string) (models.Message, error)
}

type CampaignStatus struct {
	CampaignID     string
	Status         string
	TotalMessages  int
	SentMessages   int
	QueuedMessages int
	RetryMessages  int
	FailedMessages int
	StartTime      time.Time
	CompletedTime  mysql.NullTime
}

type CampaignStatusesCollection struct {
	campaignsRepository campaignsGetter
	messages            messageCountGetter
}

func NewCampaignStatusesCollection(campaignsRepository campaignsGetter, messages messageCountGetter) CampaignStatusesCollection {
	return CampaignStatusesCollection{
		campaignsRepository: campaignsRepository,
		messages:            messages,
	}
}

func (csc CampaignStatusesCollection) Get(conn ConnectionInterface, campaignID string) (CampaignStatus, error) {
	campaign, err := csc.campaignsRepository.Get(conn, campaignID)
	if err != nil {
		panic(err)
	}

	counts, err := csc.messages.CountByStatus(conn, campaign.ID)
	if err != nil {
		panic(err)
	}

	status := CampaignStatusSending
	var completedTime mysql.NullTime

	if counts.Total > 0 && (counts.Failed+counts.Delivered) == counts.Total {
		status = CampaignStatusCompleted

		mostRecentlyUpdatedMessage, err := csc.messages.MostRecentlyUpdatedByCampaignID(conn, campaign.ID)
		if err != nil {
			panic(err)
		}

		completedTime.Time = mostRecentlyUpdatedMessage.UpdatedAt
		completedTime.Valid = true
	}

	return CampaignStatus{
		CampaignID:     campaign.ID,
		Status:         status,
		TotalMessages:  counts.Total,
		SentMessages:   counts.Delivered,
		FailedMessages: counts.Failed,
		RetryMessages:  counts.Retry,
		QueuedMessages: counts.Queued,
		StartTime:      campaign.StartTime,
		CompletedTime:  completedTime,
	}, nil
}
