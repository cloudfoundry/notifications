package postal

import (
    "log"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type UserLoader struct {
    uaaClient       UAAInterface
    logger          *log.Logger
    cloudController cf.CloudControllerInterface
}

func NewUserLoader(uaaClient UAAInterface, logger *log.Logger, cloudController cf.CloudControllerInterface) UserLoader {
    return UserLoader{
        uaaClient:       uaaClient,
        logger:          logger,
        cloudController: cloudController,
    }
}

func (loader UserLoader) Load(notificationType NotificationType, guid, token string) (map[string]uaa.User, error) {
    users := make(map[string]uaa.User)

    var guids []string
    var ccUsers []cf.CloudControllerUser
    var err error

    if notificationType == IsSpace {
        ccUsers, err = loader.cloudController.GetUsersBySpaceGuid(guid, token)
        if err != nil {
            return users, CCDownError("Cloud Controller is unavailable")
        }
    } else {
        ccUsers = []cf.CloudControllerUser{{Guid: guid}}
    }

    for _, ccUser := range ccUsers {
        loader.logger.Println("CloudController user guid: " + ccUser.Guid)
        guids = append(guids, ccUser.Guid)
    }

    usersByIDs, err := loader.uaaClient.UsersByIDs(guids...)
    if err != nil {
        err = UAAErrorFor(err)
        return users, err
    }

    for _, user := range usersByIDs {
        users[user.ID] = user
    }

    for _, guid := range guids {
        if _, ok := users[guid]; !ok {
            users[guid] = uaa.User{}
        }
    }

    return users, nil
}
