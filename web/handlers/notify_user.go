package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/smtp"
    "strings"
    "text/template"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/dgrijalva/jwt-go"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const emailBody = `The following "{{.KindDescription}}" notification was sent to you directly by the "{{.SourceDescription}}" component of Cloud Foundry:

{{.Text}}`

type NotifyUserParams struct {
    UserID            string
    ClientID          string
    From              string
    To                string
    Subject           string `json:"subject"`
    KindDescription   string `json:"kind_description"`
    SourceDescription string `json:"source_description"`
    Text              string `json:"text"`
    Kind              string
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

type NotifyUser struct {
    logger *log.Logger
}

func NewNotifyUser(logger *log.Logger) NotifyUser {
    return NotifyUser{
        logger: logger,
    }
}

func (handler NotifyUser) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    params := handler.parseParams(req)

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
    user := handler.retrieveUser(params.UserID, token)

    status := "delivered"
    if len(user.Emails) > 0 {
        params.To = user.Emails[0]
        handler.logger.Printf("Sending email to %s", params.To)
        err := handler.sendEmailTo(params)
        if err != nil {
            status = "failed"
        }
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

func (handler NotifyUser) parseParams(req *http.Request) NotifyUserParams {
    var err error

    env := config.NewEnvironment()
    params := NotifyUserParams{
        UserID: strings.TrimPrefix(req.URL.Path, "/users/"),
        From:   env.Sender,
    }

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

    return params
}

func (handler NotifyUser) retrieveUser(userID, token string) uaa.User {
    env := config.NewEnvironment()
    uaaConfig := uaa.NewUAA("", env.UAAHost, env.UAAClientID, env.UAAClientSecret, token)
    user, err := uaa.UserByID(uaaConfig, userID)
    if err != nil {
        panic(err)
    }
    return user
}

func (handler NotifyUser) sendEmailTo(context NotifyUserParams) error {
    source, err := template.New("emailBody").Parse(emailBody)
    buffer := bytes.NewBuffer([]byte{})
    source.Execute(buffer, context)

    msg := mail.Message{
        From:    context.From,
        To:      context.To,
        Subject: fmt.Sprintf("CF Notification: %s", context.Subject),
        Body:    buffer.String(),
        Headers: []string{
            fmt.Sprintf("X-CF-Client-ID: %s", context.ClientID),
        },
    }

    handler.logger.Print(msg.Data())

    env := config.NewEnvironment()
    auth := smtp.PlainAuth("", env.SMTPUser, env.SMTPPass, env.SMTPHost)
    err = smtp.SendMail(fmt.Sprintf("%s:%s", env.SMTPHost, env.SMTPPort), auth, msg.From, []string{msg.To}, []byte(msg.Data()))
    return err
}
