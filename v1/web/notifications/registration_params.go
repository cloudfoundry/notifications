package notifications

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
)

var kindIDFormat = regexp.MustCompile(`^[0-9a-zA-Z_\-.]+$`)

type RegistrationParams struct {
	SourceDescription string        `json:"source_description"`
	Kinds             []models.Kind `json:"kinds"`
	IncludesKinds     bool
}

func NewRegistrationParams(body io.ReadCloser) (RegistrationParams, error) {
	defer body.Close()

	var registration RegistrationParams
	var hashParams map[string]interface{}

	hashReader := bytes.NewBuffer([]byte{})
	structReader := bytes.NewBuffer([]byte{})
	io.Copy(io.MultiWriter(hashReader, structReader), body)

	err := json.NewDecoder(hashReader).Decode(&hashParams)
	if err != nil {
		return registration, webutil.ParseError{}
	}

	err = json.NewDecoder(structReader).Decode(&registration)
	if err != nil {
		return registration, webutil.ParseError{}
	}

	if _, ok := hashParams["kinds"]; ok {
		registration.IncludesKinds = true
	}

	return registration, nil
}

func (registration RegistrationParams) Validate() error {
	var errs []string
	if registration.SourceDescription == "" {
		errs = append(errs, `"source_description" is a required field`)
	}

	var kindErrors []string
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

	errs = append(errs, kindErrors...)

	if len(errs) > 0 {
		return webutil.ValidationError{Err: errors.New(strings.Join(errs, ", "))}
	}
	return nil
}
