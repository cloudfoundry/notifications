package postal

import (
    "log"
    "net/http"
    "net/url"
    "strings"

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
        return loader.errorFor(err)
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

func (loader UserLoader) errorFor(err error) (map[string]uaa.User, error) {
    users := make(map[string]uaa.User)

    switch err.(type) {
    case *url.Error:
        return users, UAADownError("UAA is unavailable")
    case uaa.Failure:
        uaaFailure := err.(uaa.Failure)
        loader.logger.Printf("error:  %v", err)

        if uaaFailure.Code() == http.StatusNotFound {
            if strings.Contains(uaaFailure.Message(), "Requested route") {
                return users, UAADownError("UAA is unavailable")
            } else {
                return users, UAAGenericError("UAA Unknown 404 error message: " + uaaFailure.Message())
            }
        }

        return users, UAADownError("UAA is unavailable")
    default:
        return users, UAAGenericError("UAA Unknown Error: " + err.Error())
    }
}
