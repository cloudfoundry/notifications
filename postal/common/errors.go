package common

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/uaa"
)

type UAAUserNotFoundError struct {
	Err error
}

func (e UAAUserNotFoundError) Error() string {
	return e.Err.Error()
}

type UAADownError struct {
	Err error
}

func (e UAADownError) Error() string {
	return e.Err.Error()
}

type UAAGenericError struct {
	Err error
}

func (e UAAGenericError) Error() string {
	return e.Err.Error()
}

func UAAErrorFor(err error) error {
	switch err.(type) {
	case *url.Error:
		return UAADownError{errors.New("UAA is unavailable")}
	case uaa.Failure:
		failure := err.(uaa.Failure)

		if failure.Code() == http.StatusNotFound {
			if strings.Contains(failure.Message(), "Requested route") {
				return UAADownError{errors.New("UAA is unavailable")}
			} else {
				return UAAGenericError{errors.New("UAA Unknown 404 error message: " + failure.Message())}
			}
		}

		return UAADownError{errors.New(failure.Message())}
	default:
		return UAAGenericError{errors.New("UAA Unknown Error: " + err.Error())}
	}
}
