package notify

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
)

const InvalidEmail = "<>InvalidEmail<>"

var (
	validOrganizationRoles = []string{"OrgManager", "OrgAuditor", "BillingManager"}
	emailRegexp            = regexp.MustCompile("[^<]*<([^@]*@[^@]*)>|([^<][^@]*@[^@]*)")
)

type NotifyParams struct {
	ReplyTo string `json:"reply_to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	RawHTML string `json:"html"`
	KindID  string `json:"kind_id"`
	To      string `json:"to"`
	Role    string `json:"role"`

	ParsedHTML        HTML
	KindDescription   string
	SourceDescription string
	Errors            []string
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

	err = notify.FormatEmailAndExtractHTML()
	if err != nil {
		return notify, err
	}

	return notify, nil
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

func (notify *NotifyParams) FormatEmailAndExtractHTML() error {
	notify.To = EmailFormatter{}.Format(notify.To)

	doctype, head, bodyContent, bodyAttributes, err := HTMLExtractor{}.Extract(notify.RawHTML)
	if err != nil {
		return err
	}

	notify.ParsedHTML.Doctype = doctype
	notify.ParsedHTML.Head = head
	notify.ParsedHTML.BodyContent = bodyContent
	notify.ParsedHTML.BodyAttributes = bodyAttributes

	return nil
}

type EmailFormatter struct{}

func (EmailFormatter) Format(email string) string {
	if email == "" {
		return email
	}

	matches := emailRegexp.FindStringSubmatch(email)

	if len(matches) == 0 {
		return InvalidEmail
	}

	if matches[1] != "" {
		return matches[1]
	}

	return matches[2]
}

type HTMLExtractor struct{}

func (HTMLExtractor) Extract(rawHTML string) (string, string, string, string, error) {
	reader := strings.NewReader(rawHTML)
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", "", "", "", err
	}

	doctype, err := extractDoctype(rawHTML)
	if err != nil {
		return "", "", "", "", err
	}

	head, err := extractHead(document)
	if err != nil {
		return "", "", "", "", err
	}

	bodyAttributes := ""
	for _, attribute := range document.Find("body").Nodes[0].Attr {
		bodyAttributes += " " + attribute.Key + `="` + attribute.Val + `"`
	}
	bodyAttributes = strings.TrimPrefix(bodyAttributes, " ")

	bodyContent, err := document.Find("body").Html()
	if err != nil {
		return "", "", "", "", err
	}

	return doctype, head, bodyContent, bodyAttributes, nil
}

func extractDoctype(rawHTML string) (string, error) {
	r, err := regexp.Compile("<!DOCTYPE[^>]*>")
	if err != nil {
		return "", err
	}
	return r.FindString(rawHTML), nil

}

func extractHead(document *goquery.Document) (string, error) {
	htmlHead, err := document.Find("head").Html()
	if err != nil {
		return "", err
	}

	if htmlHead == "" {
		return "", nil
	}
	return htmlHead, nil
}
