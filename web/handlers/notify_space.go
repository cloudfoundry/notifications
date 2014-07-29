package handlers

import (
    "log"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/mail"
)

type NotifySpace struct {
    logger          *log.Logger
    cloudController cf.CloudControllerInterface
    uaaClient       UAAInterface
    mailClient      mail.ClientInterface
    guidGenerator   GUIDGenerationFunc
    helper          NotifyHelper
}

func NewNotifySpace(logger *log.Logger, cloudController cf.CloudControllerInterface,
    uaaClient UAAInterface, mailClient mail.ClientInterface, guidGenerator GUIDGenerationFunc) NotifySpace {
    return NotifySpace{
        logger:          logger,
        cloudController: cloudController,
        uaaClient:       uaaClient,
        mailClient:      mailClient,
        guidGenerator:   guidGenerator,
        helper:          NewNotifyHelper(cloudController, logger, uaaClient, guidGenerator, mailClient),
    }
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    spaceGUID := strings.TrimPrefix(req.URL.Path, "/spaces/")

    loadUsers := func(spaceGuid, accessToken string) ([]cf.CloudControllerUser, error) {
        return handler.cloudController.GetUsersBySpaceGuid(spaceGuid, accessToken)
    }

    isSpace := true
    handler.helper.NotifyServeHTTP(w, req, spaceGUID, loadUsers, isSpace)
}
