package models

import "time"

const (
	DefaultTemplateID  = "default"
	DoNotSetTemplateID = ""
)

type Template struct {
	Primary    int       `db:"primary"`
	ID         string    `db:"id"`
	Name       string    `db:"name"`
	Subject    string    `db:"subject"`
	Text       string    `db:"text"`
	HTML       string    `db:"html"`
	Metadata   string    `db:"metadata"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Overridden bool      `db:"overridden"`
}
