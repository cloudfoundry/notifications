package handlers

import (
    "bytes"
    "encoding/json"
    "io"
)

type NotifySpaceParams struct {
    Kind   string
    Text   string
    Errors []string
}

func NewNotifySpaceParams(requestBody io.Reader) NotifySpaceParams {
    params := NotifySpaceParams{}

    params.parseRequestBody(requestBody)

    return params
}

func (params *NotifySpaceParams) Validate() bool {
    params.Errors = make([]string, 0)

    if params.Kind == "" {
        params.Errors = append(params.Errors, `"kind" is a required field`)
    }

    if params.Text == "" {
        params.Errors = append(params.Errors, `"text" is a required field`)
    }

    return len(params.Errors) == 0
}

func (params *NotifySpaceParams) parseRequestBody(requestBody io.Reader) {
    buffer := bytes.NewBuffer([]byte{})
    buffer.ReadFrom(requestBody)
    if buffer.Len() > 0 {
        err := json.Unmarshal(buffer.Bytes(), &params)
        if err != nil {
            panic(err)
        }
    }
}
