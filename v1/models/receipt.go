package models

import (
	"time"

	"gopkg.in/gorp.v1"
)

type Receipt struct {
	Primary   int       `db:"primary"`
	UserGUID  string    `db:"user_guid"`
	ClientID  string    `db:"client_id"`
	KindID    string    `db:"kind_id"`
	Count     int       `db:"count"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *Receipt) PreInsert(s gorp.SqlExecutor) error {
	r.CreatedAt = time.Now().Truncate(1 * time.Second).UTC()

	if r.Count == 0 {
		r.Count = 1
	}

	return nil
}
