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
	if receipt.Count == 0 {
		receipt.Count = 1
	}
	err := conn.Insert(&receipt)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			err = DuplicateRecordError{}
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
			err = NewRecordNotFoundError("Receipt for user %q of client %q and notification %q could not be found", userGUID, clientID, kindID)
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

func (repo ReceiptsRepo) Upsert(conn ConnectionInterface, receipt Receipt) error {
    query := "INSERT INTO `receipts` (`user_guid`, `client_id`, `kind_id`, `count`, `created_at`) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `count`=`count`+1"
    _, err := conn.Exec(query, receipt.UserGUID, receipt.ClientID, receipt.KindID, 1, time.Now().Truncate(1*time.Second).UTC())
    if err != nil {
        return err
    }

    return nil
}

func (repo ReceiptsRepo) CreateReceipts(conn ConnectionInterface, userGUIDs []string, clientID, kindID string) error {
    for _, guid := range userGUIDs {
        receipt := Receipt{
            UserGUID: guid,
            ClientID: clientID,
            KindID:   kindID,
        }
        err := repo.Upsert(conn, receipt)
        if err != nil {
            return err
        }
    }
    return nil
}
