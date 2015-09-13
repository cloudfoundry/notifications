package models

import (
	"database/sql"
	"errors"
	"fmt"
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

func NewKindsRepo() KindsRepo {
	return KindsRepo{}
}

func (repo KindsRepo) create(conn ConnectionInterface, kind Kind) (Kind, error) {
	err := conn.Insert(&kind)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			err = DuplicateError{errors.New("duplicate record")}
		}
		return kind, err
	}
	return kind, nil
}

func (repo KindsRepo) Find(conn ConnectionInterface, id, clientID string) (Kind, error) {
	kind := Kind{}
	err := conn.SelectOne(&kind, "SELECT * FROM `kinds` WHERE `id` = ? AND `client_id` = ?", id, clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = NotFoundError{fmt.Errorf("Notification with ID %q belonging to client %q could not be found", id, clientID)}
		}
		return kind, err
	}
	return kind, nil
}

func (repo KindsRepo) FindAll(conn ConnectionInterface) ([]Kind, error) {
	kinds := []Kind{}
	_, err := conn.Select(&kinds, `SELECT * FROM kinds`)

	if err != nil {
		return []Kind{}, err
	}

	return kinds, nil
}

func (repo KindsRepo) Update(conn ConnectionInterface, kind Kind) (Kind, error) {
	existingKind, err := repo.Find(conn, kind.ID, kind.ClientID)
	if err != nil {
		return kind, err
	}

	kind.Primary = existingKind.Primary
	kind.CreatedAt = existingKind.CreatedAt
	kind.UpdatedAt = time.Now().Truncate(1 * time.Second).UTC()
	if kind.TemplateID == DoNotSetTemplateID {
		kind.TemplateID = existingKind.TemplateID
	}

	_, err = conn.Update(&kind)
	if err != nil {
		return kind, err
	}

	return repo.Find(conn, kind.ID, kind.ClientID)
}

func (repo KindsRepo) Upsert(conn ConnectionInterface, kind Kind) (Kind, error) {
	existingKind, err := repo.Find(conn, kind.ID, kind.ClientID)
	kind.Primary = existingKind.Primary

	switch err.(type) {
	case NotFoundError:
		kind, err := repo.create(conn, kind)
		if _, ok := err.(DuplicateError); ok {
			return repo.Update(conn, kind)
		}

		return kind, err
	case nil:
		return repo.Update(conn, kind)
	default:
		return kind, err
	}
}

func (repo KindsRepo) Trim(conn ConnectionInterface, clientID string, kindIDs []string) (int, error) {
	kinds, err := repo.findAllByClientID(conn, clientID)
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

	count, err := conn.Delete(kindsToDelete...)
	return int(count), err
}

func (repo KindsRepo) findAllByClientID(conn ConnectionInterface, clientID string) ([]Kind, error) {
	var kinds []Kind
	results, err := conn.Select(Kind{}, "SELECT * FROM `kinds` WHERE `client_id` = ?", clientID)
	if err != nil {
		return kinds, err
	}
	for _, result := range results {
		kinds = append(kinds, *result.(*Kind))
	}
	return kinds, nil
}

func (repo KindsRepo) FindAllByTemplateID(conn ConnectionInterface, templateID string) ([]Kind, error) {
	kinds := []Kind{}
	_, err := conn.Select(&kinds, "SELECT * FROM `kinds` WHERE `template_id` = ?", templateID)
	if err != nil {
		return kinds, err
	}
	return kinds, nil
}
