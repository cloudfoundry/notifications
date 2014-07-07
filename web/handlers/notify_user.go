package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "strings"
    "text/template"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/dgrijalva/jwt-go"
    uuid "github.com/nu7hatch/gouuid"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const emailBody = `The following "{{.KindDescription}}" notification was sent to you directly by the "{{.SourceDescription}}" component of Cloud Foundry:

{{.Text}}`

type NotifyUserParams struct {
    Subject           string `json:"subject"`
    KindDescription   string `json:"kind_description"`
    SourceDescription string `json:"source_description"`
    Text              string `json:"text"`
    Kind              string
    UserID            string
    ClientID          string
    From              string
    To                string
    Errors            []string
}

func (params *NotifyUserParams) Invalid() bool {
    if params.Kind == "" {
        params.Errors = append(params.Errors, `"kind" is a required field`)
    }
    if params.Text == "" {
        params.Errors = append(params.Errors, `"text" is a required field`)
    }
    return len(params.Errors) > 0
}

func (params *NotifyUserParams) Parse(req *http.Request) {
    var err error

    env := config.NewEnvironment()
    params.UserID = strings.TrimPrefix(req.URL.Path, "/users/")
    params.From = env.Sender

    if authHeader := req.Header.Get("Authorization"); authHeader != "" {
        parts := strings.SplitN(authHeader, " ", 2)
        parts = strings.Split(parts[1], ".")
        decoded, err := jwt.DecodeSegment(parts[1])
        if err != nil {
            panic(err)
        }
        token := map[string]interface{}{}
        err = json.Unmarshal(decoded, &token)
        if err != nil {
            panic(err)
        }
        if clientID, ok := token["client_id"]; ok {
            params.ClientID = clientID.(string)
        }
    }

    buffer := bytes.NewBuffer([]byte{})
    buffer.ReadFrom(req.Body)
    if buffer.Len() > 0 {
        err = json.Unmarshal(buffer.Bytes(), &params)
        if err != nil {
            panic(err)
        }
    }
}

type GUIDGenerationFunc func() (*uuid.UUID, error)

type NotifyUser struct {
    logger        *log.Logger
    mailClient    mail.ClientInterface
    uaaClient     uaa.UAAInterface
    guidGenerator GUIDGenerationFunc
}

func NewNotifyUser(logger *log.Logger, mailClient mail.ClientInterface, uaaClient uaa.UAAInterface, guidGenerator GUIDGenerationFunc) NotifyUser {
    return NotifyUser{
        logger:        logger,
        mailClient:    mailClient,
        uaaClient:     uaaClient,
        guidGenerator: guidGenerator,
    }
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    params := NotifyUserParams{}
    params.Parse(req)

    if params.Invalid() {
        w.WriteHeader(422)
        response, err := json.Marshal(map[string][]string{
            "errors": params.Errors,
        })
        if err != nil {
            panic(err)
        }
        w.Write(response)
        return
    }

    token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
    handler.uaaClient.SetToken(token)
    user, err := handler.uaaClient.UserByID(params.UserID)
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

    var status string
    if len(user.Emails) > 0 {
        params.To = user.Emails[0]
        handler.logger.Printf("Sending email to %s", params.To)
        status = handler.sendEmail(params)
    }

    results := []map[string]string{
        map[string]string{
            "status": status,
        },
    }
    response, err := json.Marshal(results)
    if err != nil {
        panic(err)
    }

    w.WriteHeader(http.StatusOK)
    w.Write(response)
}

func (handler NotifyUser) sendEmail(context NotifyUserParams) string {
    source, err := template.New("emailBody").Parse(emailBody)
    buffer := bytes.NewBuffer([]byte{})
    source.Execute(buffer, context)

    guid, err := handler.guidGenerator()
    if err != nil {
        panic(err)
    }

    msg := mail.Message{
        From:    context.From,
        To:      context.To,
        Subject: fmt.Sprintf("CF Notification: %s", context.Subject),
        Body:    buffer.String(),
        Headers: []string{
            fmt.Sprintf("X-CF-Client-ID: %s", context.ClientID),
            fmt.Sprintf("X-CF-Notification-ID: %s", guid.String()),
        },
    }

    handler.logger.Print(msg.Data())

    err = handler.mailClient.Connect()
    if err != nil {
        return "unavailable"
    }

    err = handler.mailClient.Send(msg)
    if err != nil {
        return "failed"
    }

    return "delivered"
}
