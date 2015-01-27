package mail

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"text/template"

	"gopkg.in/gomail.v1"
)

const emailTemplate = "From: {{.From}}\r{{if .ReplyTo}}\nReply-To: {{.ReplyTo}}\r{{end}}\nTo: {{.To}}\r\nSubject: {{.Subject}}\r\n{{range .Headers}}{{.}}\r\n{{end}}{{.CompiledBody}}"

type Message struct {
	From         string
	ReplyTo      string
	To           string
	Subject      string
	Body         []Part
	Headers      []string
	CompiledBody string
}

type Part struct {
	ContentType string
	Content     string
}

func (msg *Message) Data() string {
	buf := bytes.NewBuffer([]byte{})

	err := msg.CompileBody()
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("test").Parse(emailTemplate)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(buf, msg)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func (msg *Message) CompileBody() error {
	message := gomail.NewMessage()
	for _, part := range msg.Body {
		message.AddAlternative(part.ContentType, "\r\n"+part.Content)
	}

	m := message.Export()
	body, err := ioutil.ReadAll(m.Body)
	if err != nil {
		panic(err)
	}

	msg.CompiledBody = string(body)

	var headers []string
	for key, _ := range m.Header {
		headers = append(headers, fmt.Sprintf("%s: %s", key, m.Header.Get(key)))
	}

	sort.Sort(sort.StringSlice(headers))

	msg.Headers = append(msg.Headers, headers...)

	return nil
}

func (msg Message) Boundary() string {
	var boundary string

	for _, header := range msg.Headers {
		matches := regexp.MustCompile(`boundary=(.*)`).FindStringSubmatch(header)
		if matches != nil {
			return matches[1]
		}
	}

	return boundary
}
