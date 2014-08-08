package models

import "time"

type Kind struct {
    ID          string    `db:"id"`
    Description string    `db:"description"`
    Critical    bool      `db:"critical"`
    ClientID    string    `db:"client_id"`
    CreatedAt   time.Time `db:"created_at"`
}
