package models

import (
	"time"

	"github.com/coopernurse/gorp"
	"github.com/nu7hatch/gouuid"
)

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

func (t *Template) PreInsert(s gorp.SqlExecutor) error {
	if t.ID == "" {
		guid, err := uuid.NewV4()
		if err != nil {
			return err
		}

		t.ID = guid.String()
	}

	if (t.CreatedAt == time.Time{}) {
		t.CreatedAt = time.Now().Truncate(1 * time.Second).UTC()
	}
	t.UpdatedAt = t.CreatedAt

	return nil
}
