package models

import (
	"time"
)

type GlobalUnsubscribe struct {
	Primary   int       `db:"primary"`
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}
