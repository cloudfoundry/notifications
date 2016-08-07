package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Campaign struct {
	ID             string         `db:"id"`
	SendTo         string         `db:"send_to"`
	CampaignTypeID string         `db:"campaign_type_id"`
	Text           string         `db:"text"`
	HTML           string         `db:"html"`
	Subject        string         `db:"subject"`
	TemplateID     string         `db:"template_id"`
	ReplyTo        string         `db:"reply_to"`
	SenderID       string         `db:"sender_id"`
	Status         string         `db:"status"`
	TotalMessages  int            `db:"total_messages"`
	SentMessages   int            `db:"sent_messages"`
	RetryMessages  int            `db:"retry_messages"`
	FailedMessages int            `db:"failed_messages"`
	StartTime      time.Time      `db:"start_time"`
	CompletedTime  mysql.NullTime `db:"completed_time"`
}

type CampaignsRepository struct {
	guidGenerator guidGeneratorFunc
	clock         clock
}

func NewCampaignsRepository(guidGenerator guidGeneratorFunc, clock clock) CampaignsRepository {
	return CampaignsRepository{
		guidGenerator: guidGenerator,
		clock:         clock,
	}
}

func (r CampaignsRepository) Insert(conn ConnectionInterface, campaign Campaign) (Campaign, error) {
	var err error
	campaign.ID, err = r.guidGenerator()
	if err != nil {
		return campaign, err
	}

	if (campaign.StartTime == time.Time{}) {
		campaign.StartTime = r.clock.Now()
	}

	err = conn.Insert(&campaign)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r CampaignsRepository) Get(conn ConnectionInterface, campaignID string) (Campaign, error) {
	campaign := Campaign{}

	err := conn.SelectOne(&campaign, "SELECT * FROM `campaigns` WHERE `id` = ?", campaignID)
	if err != nil {
		if err == sql.ErrNoRows {
			return campaign, RecordNotFoundError{fmt.Errorf("Campaign with id %q could not be found", campaignID)}
		}

		return campaign, err
	}

	return campaign, nil
}

func (r CampaignsRepository) ListSendingCampaigns(conn ConnectionInterface) ([]Campaign, error) {
	campaignList := []Campaign{}

	_, err := conn.Select(&campaignList, "SELECT * FROM `campaigns` WHERE `status` != \"completed\"")

	return campaignList, err
}
