package handlers

import (
    "log"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/mail"
    uuid "github.com/nu7hatch/gouuid"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type GUIDGenerationFunc func() (*uuid.UUID, error)

type NotifyUser struct {
    logger        *log.Logger
    mailClient    mail.ClientInterface
    uaaClient     uaa.UAAInterface
    guidGenerator GUIDGenerationFunc
    helper        NotifyHelper
}

func NewNotifyUser(logger *log.Logger, mailClient mail.ClientInterface, uaaClient uaa.UAAInterface, guidGenerator GUIDGenerationFunc) NotifyUser {
    return NotifyUser{
        logger:        logger,
        mailClient:    mailClient,
        uaaClient:     uaaClient,
        guidGenerator: guidGenerator,
        helper:        NewNotifyHelper(cf.CloudController{}, logger, uaaClient, guidGenerator, mailClient),
    }
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    userGUID := strings.TrimPrefix(req.URL.Path, "/users/")

    loadUsers := func(userGuid, accessToken string) ([]cf.CloudControllerUser, error) {
        ccUsers := []cf.CloudControllerUser{}

        user := cf.CloudControllerUser{
            Guid: userGuid,
        }

        ccUsers = append(ccUsers, user)
        return ccUsers, nil
    }

    loadSpaceAndOrganization := false
    handler.helper.NotifyServeHTTP(w, req, userGUID, loadUsers, loadSpaceAndOrganization)
}
