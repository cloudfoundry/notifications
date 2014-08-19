package models

import (
    "database/sql"
    "strings"
    "time"
)

type ReceiptsRepo struct{}

type ReceiptsRepoInterface interface {
    CreateReceipts(ConnectionInterface, []string, string, string) error
}

func NewReceiptsRepo() ReceiptsRepo {
    return ReceiptsRepo{}
}

func (repo ReceiptsRepo) Create(conn ConnectionInterface, receipt Receipt) (Receipt, error) {
    receipt.CreatedAt = time.Now().Truncate(1 * time.Second).UTC()
    receipt.Count = 1
    err := conn.Insert(&receipt)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate entry") {
            err = ErrDuplicateRecord{}
        }
        return Receipt{}, err
    }

    return receipt, nil
}

func (repo ReceiptsRepo) Find(conn ConnectionInterface, userGUID, clientID, kindID string) (Receipt, error) {
    receipt := Receipt{}
    err := conn.SelectOne(&receipt, "SELECT * FROM  `receipts` WHERE `user_guid` = ? AND `client_id` = ? AND `kind_id` = ?", userGUID, clientID, kindID)
    if err != nil {
        if err == sql.ErrNoRows {
            err = ErrRecordNotFound{}
        }
        return Receipt{}, err
    }
    return receipt, nil
}

func (repo ReceiptsRepo) Update(conn ConnectionInterface, receipt Receipt) (Receipt, error) {
    _, err := conn.Update(&receipt)
    if err != nil {
        return receipt, err
    }

    return repo.Find(conn, receipt.UserGUID, receipt.ClientID, receipt.KindID)
}

func (repo ReceiptsRepo) Upsert(conn ConnectionInterface, receipt Receipt) (Receipt, error) {
    existingReceipt, err := repo.Find(conn, receipt.UserGUID, receipt.ClientID, receipt.KindID)
    if err != nil {
        if (err == ErrRecordNotFound{}) {
            return repo.Create(conn, receipt)
        } else {
            return receipt, err
        }
    }

    receipt.Primary = existingReceipt.Primary
    receipt.CreatedAt = existingReceipt.CreatedAt
    receipt.Count = existingReceipt.Count + 1

    return repo.Update(conn, receipt)
}

func (repo ReceiptsRepo) CreateReceipts(conn ConnectionInterface, userGUIDs []string, clientID, kindID string) error {
    for _, guid := range userGUIDs {
        receipt := Receipt{
            UserGUID: guid,
            ClientID: clientID,
            KindID:   kindID,
        }
        _, err := repo.Upsert(conn, receipt)
        if err != nil {
            return err
        }
    }
    return nil
}
