package models

import "time"

type Client struct {
    ID          string    `db:"id"`
    Description string    `db:"description"`
    CreatedAt   time.Time `db:"created_at"`
}
