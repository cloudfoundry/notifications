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
    authHeader := params.req.Header.Get("Authorization")
    rawToken := strings.TrimPrefix(authHeader, "Bearer ")

    token, _ := jwt.Parse(rawToken, func(t *jwt.Token) ([]byte, error) {
        return []byte(config.UAAPublicKey), nil
    })

    if clientID, ok := token.Claims["client_id"]; ok {
        params.ClientID = clientID.(string)
    }
}

func (params *NotifyUserParams) ParseRequestPath() {
    params.UserID = strings.TrimPrefix(params.req.URL.Path, "/users/")
}

func (params *NotifyUserParams) ParseEnvironmentVariables() {
    env := config.NewEnvironment()
    params.From = env.Sender
}
