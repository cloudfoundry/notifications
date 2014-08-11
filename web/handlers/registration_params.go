package handlers

import (
    "encoding/json"
    "io"
    "io/ioutil"
    "regexp"

    "github.com/cloudfoundry-incubator/notifications/models"
)

var kindIDFormat = regexp.MustCompile(`^[0-9a-zA-z_\-.]+$`)

type RegistrationParams struct {
    SourceDescription string        `json:"source_description"`
    Kinds             []models.Kind `json:"kinds"`
}

func NewRegistrationParams(body io.Reader) (RegistrationParams, error) {
    var params RegistrationParams

    bytes, err := ioutil.ReadAll(body)
    if err != nil {
        return params, ParamsParseError{}
    }

    err = json.Unmarshal(bytes, &params)
    if err != nil {
        return params, ParamsParseError{}
    }

    return params, nil
}

func (params RegistrationParams) Validate() error {
    errors := ParamsValidationError{}
    if params.SourceDescription == "" {
        errors = append(errors, `"source_description" is a required field`)
    }

    kindErrors := ParamsValidationError{}
    for _, kind := range params.Kinds {
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
