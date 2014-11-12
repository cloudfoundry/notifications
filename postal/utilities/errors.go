package utilities

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"
)

func CCErrorFor(err error) error {
	if failure, ok := err.(cf.Failure); ok {
		if failure.Code == http.StatusNotFound {
			return CCNotFoundError(err.Error())
		}
		return CCDownError(err.Error())
	}
	return err
}

type CCNotFoundError string

func (err CCNotFoundError) Error() string {
	return "CloudController Error: " + string(err)
}

type CCDownError string

func (err CCDownError) Error() string {
	return string(err)
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
