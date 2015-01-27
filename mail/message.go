package mail

import (
	"bytes"
	"io/ioutil"
	"mime"
	"strings"
	"text/template"

	"gopkg.in/gomail.v1"
)

const emailTemplate = `{{range .Headers}}{{.}}
{{end}}Date: {{.Date}}
Mime-Version: {{.MimeVersion}}
Content-Type: {{.ContentType}}
{{if .ContentTransferEncoding}}Content-Transfer-Encoding: {{.ContentTransferEncoding}}
{{end}}From: {{.From}}{{if .ReplyTo}}
Reply-To: {{.ReplyTo}}{{end}}
To: {{.To}}
Subject: {{.Subject}}

{{.CompiledBody}}`

type Message struct {
	Date                    string
	MimeVersion             string
	ContentType             string
	ContentTransferEncoding string
	From                    string
	ReplyTo                 string
	To                      string
	Subject                 string
	Body                    []Part
	Headers                 []string
	CompiledBody            string
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
		message.AddAlternative(part.ContentType, part.Content)
	}

	m := message.Export()
	body, err := ioutil.ReadAll(m.Body)
	if err != nil {
		panic(err)
	}

	msg.CompiledBody = strings.Replace(string(body), "\r\n", "\n", -1)
	msg.Date = m.Header.Get("Date")
	msg.MimeVersion = m.Header.Get("Mime-Version")
	msg.ContentType = m.Header.Get("Content-Type")
	msg.ContentTransferEncoding = m.Header.Get("Content-Transfer-Encoding")

	return nil
}

func (msg Message) Boundary() string {
	_, params, err := mime.ParseMediaType(msg.ContentType)
	if err != nil {
		panic(err)
	}

	return params["boundary"]
}
