package models

type Sender struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	ClientID string `db:"client_id"`
}
