package postal

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type UAA string

type UAAScopesError string

func (err UAAScopesError) Error() string {
	return string(err)
}

type UAAUserNotFoundError string

func (err UAAUserNotFoundError) Error() string {
	return string(err)
}

type TemplateLoadError string

func (err TemplateLoadError) Error() string {
	return string(err)
}

type CriticalNotificationError struct {
	kindID string
}

func NewCriticalNotificationError(kindID string) CriticalNotificationError {
	return CriticalNotificationError{kindID: kindID}
}

func (err CriticalNotificationError) Error() string {
	return "Insufficient privileges to send notification " + err.kindID
}

func UAAErrorFor(err error) error {
	switch err.(type) {
	case *url.Error:
		return UAADownError("UAA is unavailable")
	case uaa.Failure:
		failure := err.(uaa.Failure)

		if failure.Code() == http.StatusNotFound {
			if strings.Contains(failure.Message(), "Requested route") {
				return UAADownError("UAA is unavailable")
			} else {
				return UAAGenericError("UAA Unknown 404 error message: " + failure.Message())
			}
		}

		return UAADownError(failure.Message())
	default:
		return UAAGenericError("UAA Unknown Error: " + err.Error())
	}
}

type UAADownError string

func (err UAADownError) Error() string {
	return string(err)
}

type UAAGenericError string

func (err UAAGenericError) Error() string {
	return string(err)
}
