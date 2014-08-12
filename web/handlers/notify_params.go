package handlers

import (
    "bytes"
    "encoding/json"
    "io"
    "strings"

    "github.com/PuerkitoBio/goquery"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type ParamsParseError struct{}

func (err ParamsParseError) Error() string {
    return "Request body could not be parsed"
}

type ParamsValidationError []string

func (err ParamsValidationError) Error() string {
    return strings.Join(err, ", ")
}

func (err ParamsValidationError) Errors() []string {
    return []string(err)
}

type NotifyParams struct {
    ReplyTo           string `json:"reply_to"`
    Subject           string `json:"subject"`
    Text              string `json:"text"`
    HTML              string `json:"html"`
    KindID            string `json:"kind_id"`
    KindDescription   string
    SourceDescription string
    Errors            []string
}

func NewNotifyParams(body io.Reader) (NotifyParams, error) {
    params := NotifyParams{}
    err := params.parseRequestBody(body)
    if err != nil {
        return params, err
    }

    err = params.extractHTML()
    if err != nil {
        return params, err
    }

    return params, nil
}

func (params *NotifyParams) Validate() bool {
    params.Errors = []string{}

    if params.KindID == "" {
        params.Errors = append(params.Errors, `"kind_id" is a required field`)
    } else {
        if !kindIDFormat.MatchString(params.KindID) {
            params.Errors = append(params.Errors, `"kind_id" is improperly formatted`)
        }
    }

    if params.Text == "" && params.HTML == "" {
        params.Errors = append(params.Errors, `"text" or "html" fields must be supplied`)
    }

    return len(params.Errors) == 0
}

func (params *NotifyParams) parseRequestBody(body io.Reader) error {
    buffer := bytes.NewBuffer([]byte{})
    buffer.ReadFrom(body)
    if buffer.Len() > 0 {
        err := json.Unmarshal(buffer.Bytes(), &params)
        if err != nil {
            return ParamsParseError{}
        }
    }
    return nil
}

func (params *NotifyParams) ToOptions(client models.Client, kind models.Kind) postal.Options {
    return postal.Options{
        ReplyTo:           params.ReplyTo,
        Subject:           params.Subject,
        KindDescription:   kind.Description,
        SourceDescription: client.Description,
        Text:              params.Text,
        HTML:              params.HTML,
        KindID:            params.KindID,
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
