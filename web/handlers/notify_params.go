package handlers

import (
    "bytes"
    "encoding/json"
    "io"
    "strings"

    "github.com/PuerkitoBio/goquery"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type NotifyParams struct {
    ReplyTo           string `json:"reply_to"`
    Subject           string `json:"subject"`
    KindDescription   string `json:"kind_description"`
    SourceDescription string `json:"source_description"`
    Text              string `json:"text"`
    HTML              string `json:"html"`
    Kind              string `json:"kind"`
    Errors            []string
}

func NewNotifyParams(body io.Reader) (NotifyParams, error) {
    params := NotifyParams{}
    params.parseRequestBody(body)
    err := params.extractHTML()
    if err != nil {
        return params, err
    }

    return params, nil
}

func (params *NotifyParams) Validate() bool {
    params.Errors = []string{}

    if params.Kind == "" {
        params.Errors = append(params.Errors, `"kind" is a required field`)
    }

    if params.Text == "" && params.HTML == "" {
        params.Errors = append(params.Errors, `"text" or "html" fields must be supplied`)
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

func (params *NotifyParams) ToOptions() postal.Options {

    return postal.Options{
        ReplyTo:           params.ReplyTo,
        Subject:           params.Subject,
        KindDescription:   params.KindDescription,
        SourceDescription: params.SourceDescription,
        Text:              params.Text,
        HTML:              params.HTML,
        Kind:              params.Kind,
    }
}

func (params *NotifyParams) extractHTML() error {

    reader := strings.NewReader(params.HTML)
    document, err := goquery.NewDocumentFromReader(reader)
    if err != nil {
        return err
    }

    params.HTML, err = document.Find("body").Html()
    if err != nil {
        return err
    }

    return nil
}
