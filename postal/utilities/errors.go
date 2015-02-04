package utilities

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/cf"
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
