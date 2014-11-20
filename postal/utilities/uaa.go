package utilities

import "github.com/pivotal-cf/uaa-sso-golang/uaa"

type UAAInterface interface {
	uaa.GetClientTokenInterface
	uaa.SetTokenInterface
	uaa.UsersEmailsByIDsInterface
	uaa.UsersGUIDsByScopeInterface
	uaa.AllUsersInterface
}
