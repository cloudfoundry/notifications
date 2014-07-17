package handlers

import (
    "bytes"
    "encoding/json"
    "io"
)

type NotifyParams struct {
    Subject           string `json:"subject"`
    KindDescription   string `json:"kind_description"`
    SourceDescription string `json:"source_description"`
    Text              string `json:"text"`
    Kind              string `json:"kind"`
    Errors            []string
}

func NewNotifyParams(body io.Reader) NotifyParams {
    params := NotifyParams{}
    params.parseRequestBody(body)
    return params
}

func (params *NotifyParams) Validate() bool {
    params.Errors = []string{}

    if params.Kind == "" {
        params.Errors = append(params.Errors, `"kind" is a required field`)
    }

    if params.Text == "" {
        params.Errors = append(params.Errors, `"text" is a required field`)
    }

    return len(params.Errors) == 0
}

func (params *NotifyParams) parseRequestBody(body io.Reader) {
    buffer := bytes.NewBuffer([]byte{})
    buffer.ReadFrom(body)
    if buffer.Len() > 0 {
        err := json.Unmarshal(buffer.Bytes(), &params)
        if err != nil {
            panic(err)
        }
    }
}
