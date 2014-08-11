package models

import (
    "database/sql"
    "strings"
    "time"
)

type IDSet []string

func (set IDSet) Contains(id string) bool {
    for _, item := range set {
        if id == item {
            return true
        }
    }
    return false
}

type KindsRepo struct{}

type KindsRepoInterface interface {
    Create(Kind) (Kind, error)
    Find(string) (Kind, error)
    Update(Kind) (Kind, error)
    Upsert(Kind) (Kind, error)
    Trim(string, []string) (int, error)
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

func (repo KindsRepo) Trim(clientID string, kindIDs []string) (int, error) {
    kinds, err := repo.findAllByClientID(clientID)
    if err != nil {
        return 0, err
    }

    ids := IDSet(kindIDs)
    var kindsToDelete []interface{}
    for _, k := range kinds {
        var kind = k
        if !ids.Contains(kind.ID) {
            kindsToDelete = append(kindsToDelete, &kind)
        }
    }

    count, err := Database().Connection.Delete(kindsToDelete...)
    return int(count), err
}

func (repo KindsRepo) findAllByClientID(clientID string) ([]Kind, error) {
    var kinds []Kind
    results, err := Database().Connection.Select(Kind{}, "SELECT * FROM `kinds` WHERE `client_id` = ?", clientID)
    if err != nil {
        return kinds, err
    }
    for _, result := range results {
        kinds = append(kinds, *result.(*Kind))
    }
    return kinds, nil
}
