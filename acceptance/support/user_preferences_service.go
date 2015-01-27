package support

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Preference struct {
	ClientID                string
	NotificationID          string
	Count                   int
	Email                   bool
	NotificationDescription string
	SourceDescription       string
}

type NotificationPreference struct {
	Count                   int    `json:"count"`
	Email                   bool   `json:"email"`
	NotificationDescription string `json:"kind_description"`
	SourceDescription       string `json:"source_description"`
}

type Preferences struct {
	GlobalUnsubscribe       bool
	NotificationPreferences []Preference
}

type PreferenceDocument struct {
	GlobalUnsubscribe bool                         `json:"global_unsubscribe"`
	Clients           map[string]ClientPreferences `json:"clients,omitempty"`
}

type ClientPreferences map[string]NotificationPreference

func (response PreferenceDocument) Preferences() Preferences {
	var preferences []Preference

	for clientID, clientPreferences := range response.Clients {
		for notificationID, notificationPreference := range clientPreferences {
			preferences = append(preferences, Preference{
				ClientID:       clientID,
				NotificationID: notificationID,
				Count:          notificationPreference.Count,
				Email:          notificationPreference.Email,
				NotificationDescription: notificationPreference.NotificationDescription,
				SourceDescription:       notificationPreference.SourceDescription,
			})
		}
	}

	return Preferences{
		GlobalUnsubscribe:       response.GlobalUnsubscribe,
		NotificationPreferences: preferences,
	}
}

type UserPreferencesService struct {
	client   *Client
	userGUID string
}

func (service UserPreferencesService) Get(token string) (int, Preferences, error) {
	var response PreferenceDocument

	status, responseBody, err := service.client.makeRequest("GET", service.client.server.SpecificUserPreferencesPath(service.userGUID), nil, token)
	if err != nil {
		return 0, response.Preferences(), err
	}

	if status == http.StatusOK {
		err = json.NewDecoder(responseBody).Decode(&response)
		if err != nil {
			return 0, response.Preferences(), err
		}
	}

	return status, response.Preferences(), nil
}

func (service UserPreferencesService) Unsubscribe(token, clientID, notificationID string) (int, error) {
	body, err := json.Marshal(PreferenceDocument{
		Clients: map[string]ClientPreferences{
			clientID: {
				notificationID: {Email: false},
			},
		},
	})
	if err != nil {
		return 0, err
	}

	status, _, err := service.client.makeRequest("PATCH", service.client.server.SpecificUserPreferencesPath(service.userGUID), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (service UserPreferencesService) GlobalUnsubscribe(token string) (int, error) {
	body, err := json.Marshal(PreferenceDocument{
		GlobalUnsubscribe: true,
	})
	if err != nil {
		return 0, err
	}

	status, _, err := service.client.makeRequest("PATCH", service.client.server.SpecificUserPreferencesPath(service.userGUID), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (service UserPreferencesService) GlobalSubscribe(token string) (int, error) {
	body, err := json.Marshal(PreferenceDocument{
		GlobalUnsubscribe: false,
	})
	if err != nil {
		return 0, err
	}

	status, _, err := service.client.makeRequest("PATCH", service.client.server.SpecificUserPreferencesPath(service.userGUID), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}
