package params

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
)

const InvalidEmail = "<>InvalidEmail<>"

var validOrganizationRoles = []string{"OrgManager", "OrgAuditor", "BillingManager"}

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
	RawHTML           string `json:"html"`
	ParsedHTML        postal.HTML
	KindID            string `json:"kind_id"`
	KindDescription   string
	SourceDescription string
	Errors            []string
	To                string `json:"to"`
	Role              string `json:"role"`
}

func NewNotify(body io.Reader) (Notify, error) {
	notify := Notify{}

	err := notify.parseRequestBody(body)
	if err != nil {
		return notify, err
	}

	notify.formatEmail()

	err = notify.extractHTML()
	if err != nil {
		return notify, err
	}

	return notify, nil
}

func (notify *Notify) formatEmail() {
	if notify.To == "" {
		return
	}
	regex := regexp.MustCompile("[^<]*<([^@]*@[^@]*)>|([^<][^@]*@[^@]*)")
	email := regex.FindStringSubmatch(notify.To)
	if len(email) == 0 {
		notify.To = InvalidEmail
		return
	}

	if email[1] != "" {
		notify.To = email[1]
	} else {
		notify.To = email[2]
	}
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
		HTML:              notify.ParsedHTML,
		KindID:            notify.KindID,
		To:                notify.To,
		Role:              notify.Role,
	}
}

func (notify *Notify) extractHTML() error {
	reader := strings.NewReader(notify.RawHTML)
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return err
	}

	notify.ParsedHTML.Doctype, err = notify.extractDoctype(notify.RawHTML)
	if err != nil {
		return err
	}

	notify.ParsedHTML.Head, err = notify.extractHead(document)
	if err != nil {
		return err
	}

	bodyAttributes := ""
	for _, attribute := range document.Find("body").Nodes[0].Attr {
		bodyAttributes += " " + attribute.Key + `="` + attribute.Val + `"`
	}
	bodyAttributes = strings.TrimPrefix(bodyAttributes, " ")

	bodyContent, err := document.Find("body").Html()
	if err != nil {
		return err
	}

	if bodyContent != "" {
		notify.ParsedHTML.BodyAttributes = bodyAttributes
		notify.ParsedHTML.BodyContent = bodyContent
	}

	return nil
}

func (notify *Notify) extractDoctype(rawHTML string) (string, error) {
	r, err := regexp.Compile("<!DOCTYPE[^>]*>")
	if err != nil {
		return "", err
	}
	return r.FindString(rawHTML), nil

}

func (notify *Notify) extractHead(document *goquery.Document) (string, error) {
	htmlHead, err := document.Find("head").Html()
	if err != nil {
		return "", err
	}

	if htmlHead == "" {
		return "", nil
	}
	return htmlHead, nil
}
