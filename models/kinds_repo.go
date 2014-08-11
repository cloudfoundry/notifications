package models

import (
    "database/sql"
    "strings"
    "time"
)

type KindsRepo struct{}

type KindsRepoInterface interface {
    Create(Kind) (Kind, error)
    Find(string) (Kind, error)
    Update(Kind) (Kind, error)
    Upsert(Kind) (Kind, error)
}

func NewKindsRepo() KindsRepo {
    return KindsRepo{}
}

func (repo KindsRepo) Create(kind Kind) (Kind, error) {
    kind.CreatedAt = time.Now().Truncate(1 * time.Second).UTC()
    err := Database().Connection.Insert(&kind)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            err = ErrDuplicateRecord{}
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
            err = ErrRecordNotFound{}
        }
        return kind, err
    }
    return kind, nil
}

func (repo KindsRepo) Update(kind Kind) (Kind, error) {
    _, err := Database().Connection.Update(&kind)
    if err != nil {
        return kind, err
    }

    return repo.Find(kind.ID)
}

func (repo KindsRepo) Upsert(kind Kind) (Kind, error) {
    existingKind, err := repo.Find(kind.ID)
    kind.CreatedAt = existingKind.CreatedAt

    if err != nil {
        if (err == ErrRecordNotFound{}) {
            return repo.Create(kind)
        } else {
            return kind, err
        }
    }
    return repo.Update(kind)
}
