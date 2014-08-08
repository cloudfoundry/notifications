package models

import (
    "database/sql"
    "strings"
    "time"
)

type KindsRepo struct{}

func NewKindsRepo() KindsRepo {
    return KindsRepo{}
}

func (repo KindsRepo) Create(kind Kind) (Kind, error) {
    kind.CreatedAt = time.Now().Truncate(1 * time.Second).UTC()
    err := Database().Connection.Insert(&kind)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            err = ErrDuplicateRecord
        }
        return kind, err
    }
    return kind, nil
}

func (repo KindsRepo) Find(id string) (Kind, error) {
    kind := Kind{}
    err := Database().Connection.SelectOne(&kind, "SELECT * FROM `kinds` WHERE `id` = ?", id)
    if err != nil {
        if err == sql.ErrNoRows {
            err = ErrRecordNotFound
        }
        return kind, err
    }
    return kind, nil
}
