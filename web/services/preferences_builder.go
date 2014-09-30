package services

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/models"
)

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
        "count":              preference.Count,
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

func (pref PreferencesBuilder) ToPreferences() ([]models.Preference, error) {
    preferences := []models.Preference{}
    for clientID, kinds := range pref {
        if len(kinds) == 0 {
            return preferences, errors.New("Missing kinds")
        }

        for kindID, kind := range kinds {

            email, ok := kind["email"]
            if !ok {
                return preferences, errors.New("Missing the email field")
            }

            shouldEmail, ok := email.(bool)
            if !ok {
                return preferences, errors.New("Email field cannot be coerced to bool")
            }

            preferences = append(preferences, models.Preference{
                ClientID: clientID,
                KindID:   kindID,
                Email:    shouldEmail,
            })
        }
    }

    return preferences, nil
}
