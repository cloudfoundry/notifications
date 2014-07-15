package handlers

import (
    "bytes"
    "encoding/json"
    "io"
)

type NotifyUserParams struct {
    Subject           string `json:"subject"`
    KindDescription   string `json:"kind_description"`
    SourceDescription string `json:"source_description"`
    Text              string `json:"text"`
    Kind              string `json:"kind"`
    Errors            []string
}

func NewNotifyUserParams(body io.Reader) NotifyUserParams {
    params := NotifyUserParams{}
    params.parseRequestBody(body)
    return params
}

func (params *NotifyUserParams) Validate() bool {
    params.Errors = []string{}

    if params.Kind == "" {
        params.Errors = append(params.Errors, `"kind" is a required field`)
    }

    if params.Text == "" {
        params.Errors = append(params.Errors, `"text" is a required field`)
    }

    return len(params.Errors) == 0
}

func (params *NotifyUserParams) parseRequestBody(body io.Reader) {
    buffer := bytes.NewBuffer([]byte{})
    buffer.ReadFrom(body)
    if buffer.Len() > 0 {
        err := json.Unmarshal(buffer.Bytes(), &params)
        if err != nil {
            panic(err)
        }
    }
}
