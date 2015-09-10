package models

type PreferencesRepo struct {
	unsubscribesRepo UnsubscribesRepo
}

func NewPreferencesRepo() PreferencesRepo {
	return PreferencesRepo{}
}

func (repo PreferencesRepo) FindNonCriticalPreferences(conn ConnectionInterface, userGUID string) ([]Preference, error) {
	preferences := []Preference{}
	sql := `SELECT DISTINCT kinds.id AS kind_id,
				clients.id AS client_id,
				kinds.description AS kind_description,
				clients.description AS source_description
			FROM kinds
			JOIN clients on kinds.client_id = clients.id
			WHERE kinds.client_id IN (
				SELECT client_id
				FROM receipts
				WHERE user_guid = ?
			)
			AND kinds.critical = false`

	_, err := conn.Select(&preferences, sql, userGUID)
	if err != nil {
		return preferences, err
	}

	unsubs, err := repo.unsubscribesRepo.FindAllByUserID(conn, userGUID)
	if err != nil {
		return preferences, err
	}

	unsubscribes := Unsubscribes(unsubs)
	for index, preference := range preferences {
		preferences[index].Email = !unsubscribes.Contains(preference.ClientID, preference.KindID)
	}

	return preferences, nil
}
