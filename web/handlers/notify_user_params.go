package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "strings"
    "time"

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
    token := params.parseAuthorizationHeader()
    if clientID, ok := token["client_id"]; ok {
        params.ClientID = clientID.(string)
    }
}

func (params *NotifyUserParams) parseAuthorizationHeader() map[string]interface{} {
    token := make(map[string]interface{})
    if authHeader := params.req.Header.Get("Authorization"); authHeader != "" {
        parts := strings.SplitN(authHeader, " ", 2)
        parts = strings.Split(parts[1], ".")
        decoded, err := jwt.DecodeSegment(parts[1])
        if err != nil {
            panic(err)
        }
        err = json.Unmarshal(decoded, &token)
        if err != nil {
            panic(err)
        }
    }
    return token
}

func (params *NotifyUserParams) ValidateAuthorizationToken() bool {
    token := params.parseAuthorizationHeader()
    if len(token) > 0 {
        if _, ok := token["client_id"]; !ok {
            params.Errors = append(params.Errors, `Authorization header is invalid: missing "client_id" field`)
        }

        if exp, ok := token["exp"]; ok {
            expirationTime := time.Unix(int64(exp.(float64)), 0)
            if expirationTime.Before(time.Now()) {
                params.Errors = append(params.Errors, "Authorization header is invalid: expired")
            }
        } else {
            params.Errors = append(params.Errors, `Authorization header is invalid: missing "exp" field`)
        }
    } else {
        params.Errors = append(params.Errors, "Authorization header is invalid: missing")
    }

    return len(params.Errors) == 0
}

func (params *NotifyUserParams) ParseRequestPath() {
    params.UserID = strings.TrimPrefix(params.req.URL.Path, "/users/")
}

func (params *NotifyUserParams) ParseEnvironmentVariables() {
    env := config.NewEnvironment()
    params.From = env.Sender
}
