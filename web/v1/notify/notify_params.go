package notify

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cloudfoundry-incubator/notifications/web/webutil"
)

const InvalidEmail = "<>InvalidEmail<>"

var validOrganizationRoles = []string{"OrgManager", "OrgAuditor", "BillingManager"}

type NotifyParams struct {
	ReplyTo           string `json:"reply_to"`
	Subject           string `json:"subject"`
	Text              string `json:"text"`
	RawHTML           string `json:"html"`
	ParsedHTML        HTML
	KindID            string `json:"kind_id"`
	KindDescription   string
	SourceDescription string
	Errors            []string
	To                string `json:"to"`
	Role              string `json:"role"`
}

type HTML struct {
	BodyContent    string
	BodyAttributes string
	Head           string
	Doctype        string
}

func NewNotifyParams(body io.ReadCloser) (NotifyParams, error) {
	notify := NotifyParams{}

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

func (notify *NotifyParams) formatEmail() {
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

func (notify *NotifyParams) parseRequestBody(body io.ReadCloser) error {
	defer body.Close()

	buffer := bytes.NewBuffer([]byte{})
	buffer.ReadFrom(body)
	if buffer.Len() > 0 {
		err := json.Unmarshal(buffer.Bytes(), &notify)
		if err != nil {
			return webutil.ParseError{}
		}
	}
	return nil
}

func (notify *NotifyParams) extractHTML() error {
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

func (notify *NotifyParams) extractDoctype(rawHTML string) (string, error) {
	r, err := regexp.Compile("<!DOCTYPE[^>]*>")
	if err != nil {
		return "", err
	}
	return r.FindString(rawHTML), nil

}

func (notify *NotifyParams) extractHead(document *goquery.Document) (string, error) {
	htmlHead, err := document.Find("head").Html()
	if err != nil {
		return "", err
	}

	if htmlHead == "" {
		return "", nil
	}
	return htmlHead, nil
}
