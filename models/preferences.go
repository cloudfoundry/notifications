package models

import "database/sql"

type Preference struct {
    ClientID          string        `db:"client_id"`
    Count             sql.NullInt64 `db:"count"`
    KindID            string        `db:"kind_id"`
    Email             bool
    KindDescription   string `db:"kind_description"`
    SourceDescription string `db:"source_description"`
}
