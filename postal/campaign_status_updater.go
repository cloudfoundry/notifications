package postal

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/v2/models"
)

type messageStatusCounter interface {
	CountByStatus(conn models.ConnectionInterface, campaignID string) (models.MessageCounts, error)
}

type campaignStatusManager interface {
	ListSendingCampaigns(conn models.ConnectionInterface) ([]models.Campaign, error)
	Update(conn models.ConnectionInterface, campaign models.Campaign) (models.Campaign, error)
}

type CampaignStatusUpdater struct {
	connection    models.ConnectionInterface
	campaignsRepo campaignStatusManager
	messagesRepo  messageStatusCounter
	pollInterval  time.Duration
}

func NewCampaignStatusUpdater(conn models.ConnectionInterface, messagesRepo messageStatusCounter, campaignsRepo campaignStatusManager, pollInterval time.Duration) CampaignStatusUpdater {
	return CampaignStatusUpdater{
		connection:    conn,
		campaignsRepo: campaignsRepo,
		messagesRepo:  messagesRepo,
		pollInterval:  pollInterval,
	}
}

func (u CampaignStatusUpdater) Run() {
	go func() {
		for {
			u.Update()
			time.Sleep(u.pollInterval)
		}
	}()
}

func (u CampaignStatusUpdater) Update() {
	campaignList, err := u.campaignsRepo.ListSendingCampaigns(u.connection)
	if err != nil {
		panic(err)
	}

	for _, campaign := range campaignList {
		messageCounts, err := u.messagesRepo.CountByStatus(u.connection, campaign.ID)
		if err != nil {
			panic(err)
		}

		campaign.Status = "completed"
		campaign.SentMessages = messageCounts.Delivered
		campaign.FailedMessages = messageCounts.Failed
		campaign.TotalMessages = messageCounts.Delivered + messageCounts.Failed

		u.campaignsRepo.Update(u.connection, campaign)
	}
}
