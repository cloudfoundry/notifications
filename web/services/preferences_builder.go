package services

import "github.com/cloudfoundry-incubator/notifications/models"

type PreferencesBuilder map[string]map[string]map[string]bool

func NewPreferencesBuilder() PreferencesBuilder {
    return map[string]map[string]map[string]bool{}
}

func (pref PreferencesBuilder) Add(client string, kind string, emails bool) {
    if clientMap, ok := pref[client]; ok {
        clientMap[kind] = map[string]bool{
            "email": emails,
        }
    } else {
        pref[client] = map[string]map[string]bool{
            kind: map[string]bool{
                "email": emails,
            },
        }

    }
}

func (pref PreferencesBuilder) ToPreferences() []models.Preference {
    preferences := []models.Preference{}
    for clientID, kind := range pref {
        for kindID, email := range kind {
            preferences = append(preferences, models.Preference{
                ClientID: clientID,
                KindID:   kindID,
                Email:    email["email"],
            })
        }
    }

    return preferences
}
