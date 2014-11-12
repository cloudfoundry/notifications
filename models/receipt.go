package models

import "time"

type Receipt struct {
	Primary   int       `db:"primary"`
	UserGUID  string    `db:"user_guid"`
	ClientID  string    `db:"client_id"`
	KindID    string    `db:"kind_id"`
	Count     int       `db:"count"`
	CreatedAt time.Time `db:"created_at"`
}
