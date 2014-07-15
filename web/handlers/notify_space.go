package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "net/url"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type NotifySpace struct {
    logger          *log.Logger
    cloudController cf.CloudControllerInterface
    uaaClient       uaa.UAAInterface
    mailClient      mail.ClientInterface
}

func NewNotifySpace(logger *log.Logger, cloudController cf.CloudControllerInterface, uaaClient uaa.UAAInterface, mailClient mail.ClientInterface) NotifySpace {
    return NotifySpace{
        logger:          logger,
        cloudController: cloudController,
        uaaClient:       uaaClient,
        mailClient:      mailClient,
    }
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    spaceGuid := strings.TrimPrefix(req.URL.Path, "/spaces/")

    params := NewNotifySpaceParams(req.Body)
    if !params.Validate() {
        handler.Error(w, 422, params.Errors)
        return
    }

    token, err := handler.uaaClient.GetClientToken()
    if err != nil {
        panic(err)
    }

    ccUsers, err := handler.cloudController.GetUsersBySpaceGuid(spaceGuid, token.Access)
    if err != nil {
        handler.Error(w, http.StatusBadGateway, []string{"Cloud Controller is unavailable"})
        return
    }

    env := config.NewEnvironment()

    for _, ccUser := range ccUsers {
        handler.logger.Println(ccUser.Guid)
        handler.uaaClient.SetToken(token.Access)
        user, err := handler.uaaClient.UserByID(ccUser.Guid)
        if err != nil {
            switch err.(type) {
            case *url.Error:
                w.WriteHeader(http.StatusBadGateway)
            case uaa.Failure:
                w.WriteHeader(http.StatusGone)
            default:
                w.WriteHeader(http.StatusInternalServerError)
            }
            return
        }

        if len(user.Emails) > 0 {
            context := MessageContext{
                From:              env.Sender,
                To:                user.Emails[0],
                Subject:           "",
                Text:              params.Text,
                Template:          "{{.Text}}",
                KindDescription:   "",
                SourceDescription: "",
                ClientID:          "",
                MessageID:         "",
            }

            status, err := SendMail(handler.mailClient, context)
            if err != nil {
                panic(err)
            }

            handler.logger.Println(status)
        }
    }
}

func (handler NotifySpace) Error(w http.ResponseWriter, code int, errors []string) {
    response, err := json.Marshal(map[string][]string{
        "errors": errors,
    })
    if err != nil {
        panic(err)
    }

    w.WriteHeader(code)
    w.Write(response)
}
