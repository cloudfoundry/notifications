package handlers

import (
    "bytes"
    "fmt"
    "text/template"

    "github.com/cloudfoundry-incubator/notifications/mail"
)

const (
    MailServerUnavailable  = "unavailable"
    MailDeliveryFailed     = "failed"
    MailDeliverySuccessful = "delivered"
)

func SendMail(client mail.ClientInterface, context MessageContext) (string, mail.Message, error) {
    sender := NewMailSender(client, context)

    compiledBody, err := sender.CompileBody()
    if err != nil {
        return "", mail.Message{}, err
    }

    message := sender.CompileMessage(compiledBody)
    return sender.Deliver(message), message, nil
}

type MessageContext struct {
    From                   string
    To                     string
    Subject                string
    Text                   string
    HTML                   string
    PlainTextEmailTemplate string
    HTMLEmailTemplate      string
    KindDescription        string
    SourceDescription      string
    ClientID               string
    MessageID              string
    Space                  string
    Organization           string
}

type MailSender struct {
    client  mail.ClientInterface
    context MessageContext
}

func NewMailSender(client mail.ClientInterface, context MessageContext) MailSender {
    return MailSender{
        client:  client,
        context: context,
    }
}

func (sender MailSender) CompileBody() (string, error) {
    var plainTextBuffer *bytes.Buffer
    var err error

    headerPart := "\nThis is a multi-part message in MIME format...\n\n"

    plainTextPart := ""
    htmlPart := ""
    closingBoundary := "--our-content-boundary--"

    if sender.context.Text != "" {
        plainTextBuffer, err = sender.compileTemplate(sender.context.PlainTextEmailTemplate)
        if err != nil {
            return "", err
        }

        plainTextPart = fmt.Sprintf("--our-content-boundary\nContent-type: text/plain\n\n%s\n", plainTextBuffer.String())
    }

    var htmlBuffer *bytes.Buffer
    if sender.context.HTML != "" {
        htmlBuffer, err = sender.compileTemplate(sender.context.HTMLEmailTemplate)
        if err != nil {
            return "", err
        }
        htmlPart = fmt.Sprintf(`--our-content-boundary
Content-Type: text/html
Content-Disposition: inline
Content-Transfer-Encoding: quoted-printable

<html>
    <body>
        %s
    </body>
</html>
`, htmlBuffer.String())
    }

    return headerPart + plainTextPart + htmlPart + closingBoundary, nil
}

func (sender MailSender) compileTemplate(theTemplate string) (*bytes.Buffer, error) {
    buffer := bytes.NewBuffer([]byte{})

    source, err := template.New("compileTemplate").Parse(theTemplate)
    if err != nil {
        return buffer, err
    }

    source.Execute(buffer, sender.context)
    return buffer, nil
}

func (sender MailSender) CompileMessage(body string) mail.Message {
    return mail.Message{
        From:    sender.context.From,
        To:      sender.context.To,
        Subject: fmt.Sprintf("CF Notification: %s", sender.context.Subject),
        Body:    body,
        Headers: []string{
            fmt.Sprintf("X-CF-Client-ID: %s", sender.context.ClientID),
            fmt.Sprintf("X-CF-Notification-ID: %s", sender.context.MessageID),
        },
    }
}

func (sender MailSender) Deliver(msg mail.Message) string {
    err := sender.client.Connect()
    if err != nil {
        return MailServerUnavailable
    }

    err = sender.client.Send(msg)
    if err != nil {
        return MailDeliveryFailed
    }

    return MailDeliverySuccessful
}
