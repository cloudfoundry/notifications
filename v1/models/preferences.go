package models

type Preference struct {
	ClientID          string `db:"client_id"`
	KindID            string `db:"kind_id"`
	KindDescription   string `db:"kind_description"`
	SourceDescription string `db:"source_description"`
	Email             bool
}
