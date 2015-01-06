package models

import "time"

type Client struct {
	Primary     int       `db:"primary"`
	ID          string    `db:"id"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	TemplateID  string    `db:"template_id"`
}

func (c Client) TemplateToUse() string {
	if c.TemplateID != "" {
		return c.TemplateID
	}

	return DefaultTemplateID
}
