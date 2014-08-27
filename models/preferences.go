package models

type Preference struct {
    ClientID string `db:"client_id"`
    KindID   string `db:"kind_id"`
    Email    bool
}
