package params

import (
    "encoding/json"
    "io"
    "io/ioutil"
    "regexp"

    "github.com/cloudfoundry-incubator/notifications/models"
)

var kindIDFormat = regexp.MustCompile(`^[0-9a-zA-Z_\-.]+$`)

type Registration struct {
    SourceDescription string        `json:"source_description"`
    Kinds             []models.Kind `json:"kinds"`
    IncludesKinds     bool
}

func NewRegistration(body io.Reader) (Registration, error) {
    var registration Registration

    bytes, err := ioutil.ReadAll(body)
    if err != nil {
        return registration, ParseError{}
    }

    var hashParams map[string]interface{}
    err = json.Unmarshal(bytes, &hashParams)
    if err != nil {
        return registration, ParseError{}
    }

    if _, ok := hashParams["kinds"]; ok {
        registration.IncludesKinds = true
    }

    err = json.Unmarshal(bytes, &registration)
    if err != nil {
        return registration, ParseError{}
    }

    return registration, nil
}

func (registration Registration) Validate() error {
    errors := ValidationError{}
    if registration.SourceDescription == "" {
        errors = append(errors, `"source_description" is a required field`)
    }

    kindErrors := ValidationError{}
    for _, kind := range registration.Kinds {
        if kind.ID == "" {
            kindErrors = append(kindErrors, `"kind.id" is a required field`)
        } else if !kindIDFormat.MatchString(kind.ID) {
            kindErrors = append(kindErrors, `"kind.id" is improperly formatted`)
        }

        if kind.Description == "" {
            kindErrors = append(kindErrors, `"kind.description" is a required field`)
        }

        if len(kindErrors) > 0 {
            break
        }
    }

    errors = append(errors, kindErrors...)

    if len(errors) > 0 {
        return errors
    }
    return nil
}
