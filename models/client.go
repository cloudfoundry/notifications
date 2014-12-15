package models

import "time"

type Client struct {
	Primary     int       `db:"primary"`
	ID          string    `db:"id"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	Template    string    `db:"-"`
}
