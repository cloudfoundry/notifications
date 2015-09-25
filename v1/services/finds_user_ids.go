package services

import "github.com/cloudfoundry-incubator/notifications/cf"

type uaaUsersGUIDsByScope interface {
	UsersGUIDsByScope(token, scope string) ([]string, error)
}

type cloudController interface {
	GetManagersByOrgGuid(orgGUID, token string) ([]cf.CloudControllerUser, error)
	GetAuditorsByOrgGuid(orgGUID, token string) ([]cf.CloudControllerUser, error)
	GetBillingManagersByOrgGuid(orgGUID, token string) ([]cf.CloudControllerUser, error)
	GetUsersByOrgGuid(orgGUID, token string) ([]cf.CloudControllerUser, error)
	GetUsersBySpaceGuid(spaceGUID, token string) ([]cf.CloudControllerUser, error)
	LoadSpace(spaceGUID, token string) (cf.CloudControllerSpace, error)
	LoadOrganization(orgGUID, token string) (cf.CloudControllerOrganization, error)
}

type FindsUserIDs struct {
	cc  cloudController
	uaa uaaUsersGUIDsByScope
}

func NewFindsUserIDs(cc cloudController, uaa uaaUsersGUIDsByScope) FindsUserIDs {
	return FindsUserIDs{
		cc:  cc,
		uaa: uaa,
	}
}

func (finder FindsUserIDs) UserIDsBelongingToSpace(spaceGUID, token string) ([]string, error) {
	var userIDs []string

	users, err := finder.cc.GetUsersBySpaceGuid(spaceGUID, token)
	if err != nil {
		return userIDs, err
	}

	for _, user := range users {
		userIDs = append(userIDs, user.GUID)
	}

	return userIDs, nil
}

func (finder FindsUserIDs) UserIDsBelongingToOrganization(orgGUID, role, token string) ([]string, error) {
	var (
		userIDs []string
		users   []cf.CloudControllerUser
		err     error
	)

	switch role {
	case "OrgManager":
		users, err = finder.cc.GetManagersByOrgGuid(orgGUID, token)
	case "OrgAuditor":
		users, err = finder.cc.GetAuditorsByOrgGuid(orgGUID, token)
	case "BillingManager":
		users, err = finder.cc.GetBillingManagersByOrgGuid(orgGUID, token)
	default:
		users, err = finder.cc.GetUsersByOrgGuid(orgGUID, token)
	}

	if err != nil {
		return userIDs, err
	}

	for _, user := range users {
		userIDs = append(userIDs, user.GUID)
	}

	return userIDs, nil
}

func (finder FindsUserIDs) UserIDsBelongingToScope(token, scope string) ([]string, error) {
	return finder.uaa.UsersGUIDsByScope(token, scope)
}
