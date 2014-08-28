package models

import "time"

type Unsubscribe struct {
    Primary   int       `db:"primary"`
    UserID    string    `db:"user_id"`
    ClientID  string    `db:"client_id"`
    KindID    string    `db:"kind_id"`
    CreatedAt time.Time `db:"created_at"`
}
