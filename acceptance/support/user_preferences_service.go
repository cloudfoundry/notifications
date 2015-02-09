package support

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type UserPreferencesService struct {
	client   *Client
	userGUID string
}

func (s UserPreferencesService) Get(token string) (int, Preferences, error) {
	var response PreferenceDocument

	status, responseBody, err := s.client.makeRequest("GET", s.client.SpecificUserPreferencesPath(s.userGUID), nil, token)
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

func (s UserPreferencesService) Unsubscribe(token, clientID, notificationID string) (int, error) {
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

	status, _, err := s.client.makeRequest("PATCH", s.client.SpecificUserPreferencesPath(s.userGUID), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (s UserPreferencesService) GlobalUnsubscribe(token string) (int, error) {
	body, err := json.Marshal(PreferenceDocument{
		GlobalUnsubscribe: true,
	})
	if err != nil {
		return 0, err
	}

	status, _, err := s.client.makeRequest("PATCH", s.client.SpecificUserPreferencesPath(s.userGUID), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}

func (s UserPreferencesService) GlobalSubscribe(token string) (int, error) {
	body, err := json.Marshal(PreferenceDocument{
		GlobalUnsubscribe: false,
	})
	if err != nil {
		return 0, err
	}

	status, _, err := s.client.makeRequest("PATCH", s.client.SpecificUserPreferencesPath(s.userGUID), bytes.NewBuffer(body), token)
	if err != nil {
		return 0, err
	}

	return status, nil
}
