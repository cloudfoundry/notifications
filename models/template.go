package models

import "time"

type Template struct {
    Primary    int       `db:"primary"`
    Name       string    `db:"name"`
    Text       string    `db:"text"`
    HTML       string    `db:"html"`
    Overridden bool      `db:"overridden"`
    CreatedAt  time.Time `db:"created_at"`
}
