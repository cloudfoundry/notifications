package models

type PreferencesRepo struct{}

type PreferencesRepoInterface interface {
    FindNonCriticalPreferences(ConnectionInterface, string) ([]Preference, error)
}

func NewPreferencesRepo() PreferencesRepo {
    return PreferencesRepo{}
}

func (repo PreferencesRepo) FindNonCriticalPreferences(conn ConnectionInterface, userGUID string) ([]Preference, error) {
    preferences := []Preference{}

    sql := `SELECT receipts.kind_id as kind_id, receipts.client_id as client_id
            FROM receipts LEFT JOIN kinds ON receipts.kind_id = kinds.id
            WHERE receipts.user_guid = ? and kinds.critical = false`

    _, err := conn.Select(&preferences, sql, userGUID)
    if err != nil {
        return preferences, err
    }

    for index, _ := range preferences {
        preferences[index].Email = "true"
    }

    return preferences, nil
}
