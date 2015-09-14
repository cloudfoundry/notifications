package collections

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/go-sql-driver/mysql"
)

const (
	CampaignStatusSending   = "sending"
	CampaignStatusCompleted = "completed"
)

type campaignGetter interface {
	Get(conn models.ConnectionInterface, campaignID string) (models.Campaign, error)
}

type senderGetter interface {
	Get(conn models.ConnectionInterface, senderID string) (models.Sender, error)
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
	campaignsRepository campaignGetter
	sendersRepository   senderGetter
	messages            messageCountGetter
}

func NewCampaignStatusesCollection(campaignsRepository campaignGetter, sendersRepository senderGetter, messages messageCountGetter) CampaignStatusesCollection {
	return CampaignStatusesCollection{
		campaignsRepository: campaignsRepository,
		sendersRepository:   sendersRepository,
		messages:            messages,
	}
}

func (csc CampaignStatusesCollection) Get(conn ConnectionInterface, campaignID, senderID string) (CampaignStatus, error) {
	campaign, err := csc.campaignsRepository.Get(conn, campaignID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return CampaignStatus{}, NotFoundError{err}
		default:
			return CampaignStatus{}, UnknownError{err}
		}
	}

	_, err = csc.sendersRepository.Get(conn, senderID)
	if err != nil {
		switch err.(type) {
		case models.RecordNotFoundError:
			return CampaignStatus{}, NotFoundError{err}
		default:
			return CampaignStatus{}, UnknownError{err}
		}
	}

	if campaign.SenderID != senderID {
		return CampaignStatus{}, NotFoundError{fmt.Errorf("Campaign with id %q could not be found", campaignID)}
	}

	counts, err := csc.messages.CountByStatus(conn, campaign.ID)
	if err != nil {
		return CampaignStatus{}, UnknownError{err}
	}

	status := CampaignStatusSending
	var completedTime mysql.NullTime

	if counts.Total > 0 && (counts.Failed+counts.Delivered) == counts.Total {
		status = CampaignStatusCompleted

		mostRecentlyUpdatedMessage, err := csc.messages.MostRecentlyUpdatedByCampaignID(conn, campaign.ID)
		if err != nil {
			return CampaignStatus{}, UnknownError{err}
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
