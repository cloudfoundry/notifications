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

type FindsUserGUIDs struct {
	cc  cloudController
	uaa uaaUsersGUIDsByScope
}

func NewFindsUserGUIDs(cc cloudController, uaa uaaUsersGUIDsByScope) FindsUserGUIDs {
	return FindsUserGUIDs{
		cc:  cc,
		uaa: uaa,
	}
}

func (finder FindsUserGUIDs) UserGUIDsBelongingToSpace(spaceGUID, token string) ([]string, error) {
	var userGUIDs []string

	users, err := finder.cc.GetUsersBySpaceGuid(spaceGUID, token)
	if err != nil {
		return userGUIDs, err
	}

	for _, user := range users {
		userGUIDs = append(userGUIDs, user.GUID)
	}

	return userGUIDs, nil
}

func (finder FindsUserGUIDs) UserGUIDsBelongingToOrganization(orgGUID, role, token string) ([]string, error) {
	var (
		userGUIDs []string
		users     []cf.CloudControllerUser
		err       error
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
		return userGUIDs, err
	}

	for _, user := range users {
		userGUIDs = append(userGUIDs, user.GUID)
	}

	return userGUIDs, nil
}

func (finder FindsUserGUIDs) UserGUIDsBelongingToScope(token, scope string) ([]string, error) {
	var userGUIDs []string

	userGUIDs, err := finder.uaa.UsersGUIDsByScope(token, scope)
	if err != nil {
		return userGUIDs, err
	}

	return userGUIDs, nil
}
