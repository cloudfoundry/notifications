package models

type PreferencesRepo struct {
    unsubscribesRepo UnsubscribesRepo
}

type PreferencesRepoInterface interface {
    FindNonCriticalPreferences(ConnectionInterface, string) ([]Preference, error)
}

func NewPreferencesRepo() PreferencesRepo {
    return PreferencesRepo{}
}

func (repo PreferencesRepo) FindNonCriticalPreferences(conn ConnectionInterface, userGUID string) ([]Preference, error) {
    preferences := []Preference{}
    sql := `SELECT kinds.id as kind_id, kinds.client_id as client_id, kinds.description as kind_description, clients.description as source_description
            FROM kinds
            LEFT OUTER JOIN receipts ON kinds.id = receipts.kind_id
            LEFT OUTER JOIN clients ON clients.id = kinds.client_id
            WHERE kinds.critical = "false"
            AND kinds.client_id IN
                (SELECT DISTINCT kinds.client_id
                 FROM kinds
                 JOIN receipts ON kinds.client_id = receipts.client_id
                 WHERE receipts.user_guid = ?)`

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
