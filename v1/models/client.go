package models

import (
	"time"

	"gopkg.in/gorp.v1"
)

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

func (c *Client) PreInsert(s gorp.SqlExecutor) error {
	c.CreatedAt = time.Now().Truncate(1 * time.Second).UTC()

	if c.TemplateID == "" {
		c.TemplateID = DefaultTemplateID
	}

	return nil
}
