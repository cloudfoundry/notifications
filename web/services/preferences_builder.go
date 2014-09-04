package services

import "github.com/cloudfoundry-incubator/notifications/models"

type PreferencesBuilder map[string]map[string]map[string]interface{}

func NewPreferencesBuilder() PreferencesBuilder {
    return map[string]map[string]map[string]interface{}{}
}

func (pref PreferencesBuilder) Add(preference models.Preference) {
    if preference.KindDescription == "" {
        preference.KindDescription = preference.KindID
    }

    if preference.SourceDescription == "" {
        preference.SourceDescription = preference.ClientID
    }

    data := map[string]interface{}{
        "email":              preference.Email,
        "kind_description":   preference.KindDescription,
        "source_description": preference.SourceDescription,
    }

    if clientMap, ok := pref[preference.ClientID]; ok {
        clientMap[preference.KindID] = data
    } else {
        pref[preference.ClientID] = map[string]map[string]interface{}{
            preference.KindID: data,
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
                Email:    email["email"].(bool),
            })
        }
    }

    return preferences
}
