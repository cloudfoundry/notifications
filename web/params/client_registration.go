package params

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type ClientRegistration struct {
	SourceName    string                           `json:"source_name"`
	Notifications map[string](*NotificationStruct) `json:"notifications"`
}

type NotificationStruct struct {
	ID          string
	Description string `json:"description"`
	Critical    bool   `json:"critical"`
}

func NewClientRegistration(body io.Reader) (ClientRegistration, error) {
	var clientRegistration ClientRegistration

	decodeReader := bytes.NewBuffer([]byte{})
	validateReader := bytes.NewBuffer([]byte{})
	io.Copy(io.MultiWriter(decodeReader, validateReader), body)

	err := json.NewDecoder(decodeReader).Decode(&clientRegistration)
	if err != nil {
		return clientRegistration, ParseError{}
	}

	err = strictValidateJSON(validateReader.Bytes())
	if err != nil {
		return clientRegistration, err
	}

	for id, _ := range clientRegistration.Notifications {
		clientRegistration.Notifications[id].ID = id
	}

	return clientRegistration, nil
}

func strictValidateJSON(bytes []byte) error {
	var untypedClientRegistration map[string]interface{}
	err := json.Unmarshal(bytes, &untypedClientRegistration)
	if err != nil {
		return err
	}

	for key, _ := range untypedClientRegistration {
		if key == "source_name" {
			continue
		} else if key == "notifications" {
			if untypedClientRegistration[key] == nil {
				return SchemaError(fmt.Sprintf(`only include "notifications" key when adding a notification"`))
			}
			notifications := untypedClientRegistration[key].(map[string]interface{})
			for _, notificationData := range notifications {
				if notificationData == nil {
					return SchemaError(fmt.Sprintf(`notification must not be null`))
				}
				notificationMap := notificationData.(map[string]interface{})
				for propertyName, _ := range notificationMap {
					if propertyName == "description" || propertyName == "critical" {
						continue
					} else {
						return SchemaError(fmt.Sprintf(`"%+v" is not a valid property`, propertyName))
					}
				}
			}
		} else {
			return SchemaError(fmt.Sprintf(`"%+v" is not a valid property`, key))
		}
	}
	return nil
}

func (clientRegistration ClientRegistration) Validate() error {
	errors := ValidationError{}
	if clientRegistration.SourceName == "" {
		errors = append(errors, `"source_name" is a required field`)
	}

	for id, value := range clientRegistration.Notifications {
		if value == nil {
			errors = append(errors, fmt.Sprintf(`notification "%+v" is empty`, id))
		}
		if value.ID == "" {
			errors = append(errors, fmt.Sprintf(`notification "%+v" is missing required field "ID"`, id))
		}
		if value.Description == "" {
			errors = append(errors, fmt.Sprintf(`notification "%+v" is missing required field "Description"`, id))
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}
