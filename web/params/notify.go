package params

import (
    "bytes"
    "encoding/json"
    "io"
    "strings"

    "github.com/PuerkitoBio/goquery"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type ParseError struct{}

func (err ParseError) Error() string {
    return "Request body could not be parsed"
}

type ValidationError []string

func (err ValidationError) Error() string {
    return strings.Join(err, ", ")
}

func (err ValidationError) Errors() []string {
    return []string(err)
}

type Notify struct {
    ReplyTo           string `json:"reply_to"`
    Subject           string `json:"subject"`
    Text              string `json:"text"`
    HTML              string `json:"html"`
    KindID            string `json:"kind_id"`
    KindDescription   string
    SourceDescription string
    Errors            []string
}

func NewNotify(body io.Reader) (Notify, error) {
    notify := Notify{}
    err := notify.parseRequestBody(body)
    if err != nil {
        return notify, err
    }

    err = notify.extractHTML()
    if err != nil {
        return notify, err
    }

    return notify, nil
}

func (notify *Notify) Validate() bool {
    notify.Errors = []string{}

    if notify.KindID == "" {
        notify.Errors = append(notify.Errors, `"kind_id" is a required field`)
    } else {
        if !kindIDFormat.MatchString(notify.KindID) {
            notify.Errors = append(notify.Errors, `"kind_id" is improperly formatted`)
        }
    }

    if notify.Text == "" && notify.HTML == "" {
        notify.Errors = append(notify.Errors, `"text" or "html" fields must be supplied`)
    }

    return len(notify.Errors) == 0
}

func (notify *Notify) parseRequestBody(body io.Reader) error {
    buffer := bytes.NewBuffer([]byte{})
    buffer.ReadFrom(body)
    if buffer.Len() > 0 {
        err := json.Unmarshal(buffer.Bytes(), &notify)
        if err != nil {
            return ParseError{}
        }
    }
    return nil
}

func (notify *Notify) ToOptions(client models.Client, kind models.Kind) postal.Options {
    return postal.Options{
        ReplyTo:           notify.ReplyTo,
        Subject:           notify.Subject,
        KindDescription:   kind.Description,
        SourceDescription: client.Description,
        Text:              notify.Text,
        HTML:              notify.HTML,
        KindID:            notify.KindID,
    }
}

func (notify *Notify) extractHTML() error {

    reader := strings.NewReader(notify.HTML)
    document, err := goquery.NewDocumentFromReader(reader)
    if err != nil {
        return err
    }

    notify.HTML, err = document.Find("body").Html()
    if err != nil {
        return err
    }

    return nil
}
