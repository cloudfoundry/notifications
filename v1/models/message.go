package models

import (
	"time"

	"gopkg.in/gorp.v1"
)

type Message struct {
	ID         string    `db:"id"`
	Status     string    `db:"status"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (m *Message) PreInsert(s gorp.SqlExecutor) error {
	m.UpdatedAt = time.Now().Truncate(1 * time.Second).UTC()

	return nil
}

func (m *Message) PreUpdate(s gorp.SqlExecutor) error {
	m.UpdatedAt = time.Now().Truncate(1 * time.Second).UTC()

	return nil
}
