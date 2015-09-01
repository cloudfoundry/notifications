package models

type Campaign struct {
	ID             string `db:"id"`
	SendTo         string `db:"send_to"`
	CampaignTypeID string `db:"campaign_type_id"`
	Text           string `db:"text"`
	HTML           string `db:"html"`
	Subject        string `db:"subject"`
	TemplateID     string `db:"template_id"`
	ReplyTo        string `db:"reply_to"`
	SenderID       string `db:"sender_id"`
}
