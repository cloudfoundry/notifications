package models

import (
	"time"

	"gopkg.in/gorp.v1"
)

type Kind struct {
	Primary     int       `db:"primary"`
	ID          string    `db:"id"`
	Description string    `db:"description"`
	Critical    bool      `db:"critical"`
	ClientID    string    `db:"client_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	TemplateID  string    `db:"template_id"`
}

func (k Kind) TemplateToUse() string {
	if k.TemplateID != "" {
		return k.TemplateID
	}

	return DefaultTemplateID
}

func (k *Kind) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now().Truncate(1 * time.Second).UTC()
	k.CreatedAt = now
	k.UpdatedAt = now

	if k.TemplateID == "" {
		k.TemplateID = DefaultTemplateID
	}

	return nil
}
