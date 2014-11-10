package utilities

import "github.com/cloudfoundry-incubator/notifications/cf"

type FindsUserGUIDs struct {
    cloudController cf.CloudControllerInterface
}

type FindsUserGUIDsInterface interface {
    UserGUIDsBelongingToSpace(string, string) ([]string, error)
    UserGUIDsBelongingToOrganization(string, string) ([]string, error)
}

func NewFindsUserGUIDs(cloudController cf.CloudControllerInterface) FindsUserGUIDs {
    return FindsUserGUIDs{
        cloudController: cloudController,
    }
}

func (finder FindsUserGUIDs) UserGUIDsBelongingToSpace(spaceGUID, token string) ([]string, error) {
    var userGUIDs []string

    users, err := finder.cloudController.GetUsersBySpaceGuid(spaceGUID, token)
    if err != nil {
        return userGUIDs, err
    }

    for _, user := range users {
        userGUIDs = append(userGUIDs, user.GUID)
    }

    return userGUIDs, nil
}

func (finder FindsUserGUIDs) UserGUIDsBelongingToOrganization(orgGUID, token string) ([]string, error) {
    var userGUIDs []string

    users, err := finder.cloudController.GetUsersByOrgGuid(orgGUID, token)
    if err != nil {
        return userGUIDs, err
    }

    for _, user := range users {
        userGUIDs = append(userGUIDs, user.GUID)
    }

    return userGUIDs, nil
}
