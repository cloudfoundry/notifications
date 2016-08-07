package models

import (
	"time"

	"gopkg.in/gorp.v1"
)

type Unsubscribe struct {
	Primary   int       `db:"primary"`
	UserID    string    `db:"user_id"`
	ClientID  string    `db:"client_id"`
	KindID    string    `db:"kind_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (u *Unsubscribe) PreInsert(s gorp.SqlExecutor) error {
	u.CreatedAt = time.Now().Truncate(1 * time.Second).UTC()

	return nil
}

type Unsubscribes []Unsubscribe

func (unsubscribes Unsubscribes) Contains(clientID, kindID string) bool {
	for _, unsubscribe := range unsubscribes {
		if unsubscribe.ClientID == clientID && unsubscribe.KindID == kindID {
			return true
		}
	}
	return false
}
