package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "net/url"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type NotifyHelper struct{}

func (helper NotifyHelper) Error(w http.ResponseWriter, code int, errors []string) {
    response, err := json.Marshal(map[string][]string{
        "errors": errors,
    })
    if err != nil {
        panic(err)
    }

    w.WriteHeader(code)
    w.Write(response)
}

func (helper NotifyHelper) LoadUser(w http.ResponseWriter, guid string, uaaClient uaa.UAAInterface) (uaa.User, bool) {
    user, err := uaaClient.UserByID(guid)
    if err != nil {
        switch err.(type) {
        case *url.Error:
            w.WriteHeader(http.StatusBadGateway)
        case uaa.Failure:
            w.WriteHeader(http.StatusGone)
        default:
            w.WriteHeader(http.StatusInternalServerError)
        }
        return uaa.User{}, false
    }
    return user, true
}

func (helper NotifyHelper) SendMailToUser(context MessageContext, logger *log.Logger, mailClient mail.ClientInterface) string {
    logger.Printf("Sending email to %s", context.To)
    status, message, err := SendMail(mailClient, context)
    if err != nil {
        panic(err)
    }

    logger.Print(message.Data())
    return status
}

func (helper NotifyHelper) BuildSpaceContext(user uaa.User, params NotifyParams, env config.Environment, space, organization, clientID string, guidGenerator GUIDGenerationFunc, plainTextEmailTemplate, htmlEmailTemplate string) MessageContext {
    return helper.buildContext(user, params, env, space, organization, clientID, guidGenerator, plainTextEmailTemplate, htmlEmailTemplate)
}

func (helper NotifyHelper) BuildUserContext(user uaa.User, params NotifyParams, env config.Environment, clientID string, guidGenerator GUIDGenerationFunc, plainTextEmailTemplate, htmlEmailTexmplate string) MessageContext {
    return helper.buildContext(user, params, env, "", "", clientID, guidGenerator, plainTextEmailTemplate, htmlEmailTemplate)
}

func (handler NotifyHelper) buildContext(user uaa.User, params NotifyParams, env config.Environment, space, organization, clientID string, guidGenerator GUIDGenerationFunc, plainTextEmailTemplate, htmlEmailTemplate string) MessageContext {
    guid, err := guidGenerator()
    if err != nil {
        panic(err)
    }

    var kindDescription string
    if params.KindDescription == "" {
        kindDescription = params.Kind
    } else {
        kindDescription = params.KindDescription
    }

    var sourceDescription string
    if params.SourceDescription == "" {
        sourceDescription = clientID
    } else {
        sourceDescription = params.SourceDescription
    }

    return MessageContext{
        From:    env.Sender,
        To:      user.Emails[0],
        Subject: params.Subject,
        Text:    params.Text,
        HTML:    params.HTML,
        PlainTextEmailTemplate: plainTextEmailTemplate,
        HTMLEmailTemplate:      htmlEmailTemplate,
        KindDescription:        kindDescription,
        SourceDescription:      sourceDescription,
        ClientID:               clientID,
        MessageID:              guid.String(),
        Space:                  space,
        Organization:           organization,
    }
}
