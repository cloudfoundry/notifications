package models

type Preference struct {
    ClientID          string `db:"client_id"`
    Count             int    `db:"count"`
    KindID            string `db:"kind_id"`
    Email             bool
    KindDescription   string `db:"kind_description"`
    SourceDescription string `db:"source_description"`
}
