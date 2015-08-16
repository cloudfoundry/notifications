package models

import "time"

type ReceiptsRepo struct{}

func NewReceiptsRepo() ReceiptsRepo {
	return ReceiptsRepo{}
}

func (repo ReceiptsRepo) upsert(conn ConnectionInterface, receipt Receipt) error {
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
		err := repo.upsert(conn, receipt)
		if err != nil {
			return err
		}
	}
	return nil
}
