package handlers

import (
    "log"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
)

type NotifySpace struct {
    logger          *log.Logger
    cloudController cf.CloudControllerInterface
}

func NewNotifySpace(logger *log.Logger, cloudController cf.CloudControllerInterface) NotifySpace {
    return NotifySpace{
        logger:          logger,
        cloudController: cloudController,
    }
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    spaceGuid := strings.TrimPrefix(req.URL.Path, "/spaces/")
    token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")

    users, err := handler.cloudController.GetUsersBySpaceGuid(spaceGuid, token)
    if err != nil {
        panic(err)
    }

    for _, user := range users {
        handler.logger.Println(user.Guid)
    }
}
