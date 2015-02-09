package support

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type PreferencesService struct {
	client *Client
}

func (service PreferencesService) Get(token string) (int, Preferences, error) {
	var response PreferenceDocument

	status, responseBody, err := service.client.makeRequest("GET", service.client.server.UserPreferencesPath(), nil, token)
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

func (service PreferencesService) Unsubscribe(token, clientID, notificationID string) (int, error) {
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

	status, _, err := service.client.makeRequest("PATCH", service.client.server.UserPreferencesPath(), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (service PreferencesService) Subscribe(token, clientID, notificationID string) (int, error) {
	body, err := json.Marshal(PreferenceDocument{
		Clients: map[string]ClientPreferences{
			clientID: {
				notificationID: {Email: true},
			},
		},
	})
	if err != nil {
		return 0, err
	}

	status, _, err := service.client.makeRequest("PATCH", service.client.server.UserPreferencesPath(), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (service PreferencesService) GlobalUnsubscribe(token string) (int, error) {
	body, err := json.Marshal(PreferenceDocument{
		GlobalUnsubscribe: true,
	})
	if err != nil {
		return 0, err
	}

	status, _, err := service.client.makeRequest("PATCH", service.client.server.UserPreferencesPath(), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (service PreferencesService) GlobalSubscribe(token string) (int, error) {
	body, err := json.Marshal(PreferenceDocument{
		GlobalUnsubscribe: false,
	})
	if err != nil {
		return 0, err
	}

	status, _, err := service.client.makeRequest("PATCH", service.client.server.UserPreferencesPath(), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (service PreferencesService) User(userGUID string) UserPreferencesService {
	return UserPreferencesService{
		client:   service.client,
		userGUID: userGUID,
	}
}
