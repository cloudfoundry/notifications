package models

import "time"

const DefaultTemplateID = "default"

type Client struct {
	Primary     int       `db:"primary"`
	ID          string    `db:"id"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	Template    string    `db:"template_id"`
}

func (c *Client) TemplateToUse() string {
	if c.Template != "" {
		return c.Template
	}

	return DefaultTemplateID
}
