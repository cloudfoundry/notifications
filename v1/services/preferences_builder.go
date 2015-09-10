package services

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type Kind struct {
	Email             *bool  `json:"email"`
	KindDescription   string `json:"kind_description"`
	SourceDescription string `json:"source_description"`
}

type ClientMap map[string]Kind
type ClientsMap map[string]ClientMap

type PreferencesBuilder struct {
	GlobalUnsubscribe bool       `json:"global_unsubscribe"`
	Clients           ClientsMap `json:"clients"`
}

func NewPreferencesBuilder() PreferencesBuilder {
	return PreferencesBuilder{
		Clients: ClientsMap{},
	}
}

func (pref PreferencesBuilder) Add(preference models.Preference) {
	if preference.KindDescription == "" {
		preference.KindDescription = preference.KindID
	}

	if preference.SourceDescription == "" {
		preference.SourceDescription = preference.ClientID
	}

	data := Kind{
		Email:             &preference.Email,
		KindDescription:   preference.KindDescription,
		SourceDescription: preference.SourceDescription,
	}

	if clientMap, ok := pref.Clients[preference.ClientID]; ok {
		clientMap[preference.KindID] = data
	} else {
		pref.Clients[preference.ClientID] = ClientMap{
			preference.KindID: data,
		}
	}
}

func (pref PreferencesBuilder) ToPreferences() ([]models.Preference, error) {
	preferences := []models.Preference{}
	for clientID, kinds := range pref.Clients {
		if len(kinds) == 0 {
			return preferences, errors.New("Missing kinds")
		}

		for kindID, kind := range kinds {

			if kind.Email == nil {
				return preferences, errors.New("Missing the email field")
			}

			preferences = append(preferences, models.Preference{
				ClientID: clientID,
				KindID:   kindID,
				Email:    *kind.Email,
			})
		}
	}

	return preferences, nil
}
