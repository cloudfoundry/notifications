package models

type Template struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	HTML     string `db:"html"`
	Text     string `db:"text"`
	Subject  string `db:"subject"`
	Metadata string `db:"metadata"`
	ClientID string `db:"client_id"`
}
