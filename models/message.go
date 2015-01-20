package models

import "time"

type Message struct {
	ID        string    `db:"id"`
	Status    string    `db:"status"`
	UpdatedAt time.Time `db:"updated_at"`
}
