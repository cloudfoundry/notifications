package notifications

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
)

type ClientRegistrationParams struct {
	SourceName    string                           `json:"source_name"`
	Notifications map[string](*NotificationStruct) `json:"notifications"`
}

type NotificationStruct struct {
	ID          string
	Description string `json:"description"`
	Critical    bool   `json:"critical"`
}

func NewClientRegistrationParams(body io.Reader) (ClientRegistrationParams, error) {
	var clientRegistration ClientRegistrationParams

	decodeReader := bytes.NewBuffer([]byte{})
	validateReader := bytes.NewBuffer([]byte{})
	io.Copy(io.MultiWriter(decodeReader, validateReader), body)

	err := json.NewDecoder(decodeReader).Decode(&clientRegistration)
	if err != nil {
		return clientRegistration, webutil.ParseError{}
	}

	err = strictValidateJSON(validateReader.Bytes())
	if err != nil {
		return clientRegistration, err
	}

	for id := range clientRegistration.Notifications {
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

	for key := range untypedClientRegistration {
		if key == "source_name" {
			continue
		} else if key == "notifications" {
			if untypedClientRegistration[key] == nil {
				return webutil.SchemaError{Err: errors.New("only include \"notifications\" key when adding a notification")}
			}
			notifications := untypedClientRegistration[key].(map[string]interface{})
			for _, notificationData := range notifications {
				if notificationData == nil {
					return webutil.SchemaError{Err: errors.New("notification must not be null")}
				}
				notificationMap := notificationData.(map[string]interface{})
				for propertyName := range notificationMap {
					if propertyName == "description" || propertyName == "critical" {
						continue
					} else {
						return webutil.SchemaError{Err: fmt.Errorf("%q is not a valid property", propertyName)}
					}
				}
			}
		} else {
			return webutil.SchemaError{Err: fmt.Errorf("%q is not a valid property", key)}
		}
	}
	return nil
}

func (clientRegistration ClientRegistrationParams) Validate() error {
	var errs []string
	if clientRegistration.SourceName == "" {
		errs = append(errs, `"source_name" is a required field`)
	}

	for id, value := range clientRegistration.Notifications {
		if value == nil {
			errs = append(errs, fmt.Sprintf(`notification "%+v" is empty`, id))
		}
		if value.ID == "" {
			errs = append(errs, fmt.Sprintf(`notification "%+v" is missing required field "ID"`, id))
		}
		if value.Description == "" {
			errs = append(errs, fmt.Sprintf(`notification "%+v" is missing required field "Description"`, id))
		}
	}

	if len(errs) > 0 {
		return webutil.ValidationError{Err: errors.New(strings.Join(errs, ", "))}
	}

	return nil
}
