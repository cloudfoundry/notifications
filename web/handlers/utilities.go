package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "net/url"

    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

func Error(w http.ResponseWriter, code int, errors []string) {
    response, err := json.Marshal(map[string][]string{
        "errors": errors,
    })
    if err != nil {
        panic(err)
    }

    w.WriteHeader(code)
    w.Write(response)
}

func loadUser(w http.ResponseWriter, guid string, uaaClient uaa.UAAInterface) (uaa.User, bool) {
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

func sendMailToUser(context MessageContext, logger *log.Logger, mailClient mail.ClientInterface) string {
    logger.Printf("Sending email to %s", context.To)
    status, message, err := SendMail(mailClient, context)
    if err != nil {
        panic(err)
    }

    logger.Print(message.Data())
    return status
}
