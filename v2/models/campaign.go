package models

import (
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
