package models

import "time"

type Message struct {
	ID        string    `db:"id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}
