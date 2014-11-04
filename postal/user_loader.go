package postal

import (
    "log"
    "time"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/metrics"
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

func (loader UserLoader) Load(guid TypedGUID, token string) (map[string]uaa.User, error) {
    users := make(map[string]uaa.User)

    var guids []string
    var ccUsers []cf.CloudControllerUser
    var err error

    if guid.BelongsToSpace() {
        ccUsers, err = loader.cloudController.GetUsersBySpaceGuid(guid.String(), token)
    } else if guid.BelongsToOrganization() {
        ccUsers, err = loader.cloudController.GetUsersByOrgGuid(guid.String(), token)
    } else {
        ccUsers = []cf.CloudControllerUser{{Guid: guid.String()}}
    }

    if err != nil {
        return users, CCDownError("Cloud Controller is unavailable")
    }

    for _, ccUser := range ccUsers {
        loader.logger.Println("CloudController user guid: " + ccUser.Guid)
        guids = append(guids, ccUser.Guid)
    }

    usersByIDs, err := loader.fetchUsersByIDs(guids)
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

func (loader UserLoader) fetchUsersByIDs(guids []string) ([]uaa.User, error) {
    then := time.Now()

    usersByIDs, err := loader.uaaClient.UsersEmailsByIDs(guids...)

    duration := time.Now().Sub(then)

    metrics.NewMetric("histogram", map[string]interface{}{
        "name":  "notifications.external-requests.uaa.users-email",
        "value": duration.Seconds(),
    }).Log()

    return usersByIDs, err
}
