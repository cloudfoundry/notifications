package utilities

import "github.com/cloudfoundry-incubator/notifications/uaa"

type UAAInterface interface {
	UsersGUIDsByScope(string) ([]string, error)
	AllUsers() ([]uaa.User, error)
}
