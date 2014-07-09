package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/dgrijalva/jwt-go"
)

type NotifyUserParams struct {
    req               *http.Request
    Subject           string `json:"subject"`
    KindDescription   string `json:"kind_description"`
    SourceDescription string `json:"source_description"`
    Text              string `json:"text"`
    Kind              string `json:"kind"`
    UserID            string
    ClientID          string
    From              string
    To                string
    Errors            []string
}

func NewNotifyUserParams(req *http.Request) NotifyUserParams {
    return NotifyUserParams{
        req: req,
    }
}

func (params *NotifyUserParams) ValidateRequestBody() bool {
    params.Errors = []string{}

    if params.Kind == "" {
        params.Errors = append(params.Errors, `"kind" is a required field`)
    }

    if params.Text == "" {
        params.Errors = append(params.Errors, `"text" is a required field`)
    }

    return len(params.Errors) == 0
}

func (params *NotifyUserParams) ParseRequestBody() {
    buffer := bytes.NewBuffer([]byte{})
    buffer.ReadFrom(params.req.Body)
    if buffer.Len() > 0 {
        err := json.Unmarshal(buffer.Bytes(), &params)
        if err != nil {
            panic(err)
        }
    }
}

func (params *NotifyUserParams) ParseAuthorizationToken() {
    token, _ := params.parseAuthorizationHeader()

    if clientID, ok := token.Claims["client_id"]; ok {
        params.ClientID = clientID.(string)
    }
}

func (params *NotifyUserParams) ValidateAuthorizationToken() bool {
    params.Errors = []string{}
    token, err := params.parseAuthorizationHeader()
    if err == nil {
        if _, ok := token.Claims["client_id"]; !ok {
            params.Errors = append(params.Errors, `Authorization header is invalid: missing "client_id" field`)
        }
    } else {
        if strings.Contains(err.Error(), "Token is expired") {
            params.Errors = append(params.Errors, "Authorization header is invalid: expired")
        } else {
            params.Errors = append(params.Errors, "Authorization header is invalid: missing")
        }
    }

    return len(params.Errors) == 0
}

func (params *NotifyUserParams) ConfirmPermissions() bool {
    params.Errors = []string{}
    token, err := params.parseAuthorizationHeader()
    if err != nil {
        params.Errors = append(params.Errors, err.Error())
        return false
    }

    if scopes, ok := token.Claims["scope"]; ok {
        hasNotificationsWrite := false
        for _, scope := range scopes.([]interface{}) {
            if scope.(string) == "notifications.write" {
                hasNotificationsWrite = true
                break
            }
        }
        if !hasNotificationsWrite {
            params.Errors = append(params.Errors, "You are not authorized to perform the requested action")
        }
    } else {
        params.Errors = append(params.Errors, "You are not authorized to perform the requested action")
    }

    return len(params.Errors) == 0
}

func (params *NotifyUserParams) parseAuthorizationHeader() (*jwt.Token, error) {
    authHeader := params.req.Header.Get("Authorization")
    rawToken := strings.TrimPrefix(authHeader, "Bearer ")

    token, err := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })
    if err != nil {
        return &jwt.Token{}, err
    }

    return token, nil
}

func (params *NotifyUserParams) ParseRequestPath() {
    params.UserID = strings.TrimPrefix(params.req.URL.Path, "/users/")
}

func (params *NotifyUserParams) ParseEnvironmentVariables() {
    env := config.NewEnvironment()
    params.From = env.Sender
}
