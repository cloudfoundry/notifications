package mail

import (
    "bytes"
    "text/template"
)

const emailTemplate = `{{range .Headers}}{{.}}
{{end}}From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="our-content-boundary"

{{.Body}}`

type Message struct {
    From    string
    To      string
    Subject string
    Body    string
    Headers []string
}

func (msg Message) Data() string {
    buf := bytes.NewBuffer([]byte{})

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
