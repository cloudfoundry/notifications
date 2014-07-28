package handlers

import (
    "bytes"
    "fmt"
    "html"
    "strings"
    "text/template"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
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

    message, err := sender.CompileMessage(compiledBody)
    if err != nil {
        return "", message, err
    }

    return sender.Deliver(message), message, nil
}

type MessageContext struct {
    From                   string
    ReplyTo                string
    To                     string
    Subject                string
    Text                   string
    HTML                   string
    PlainTextEmailTemplate string
    HTMLEmailTemplate      string
    SubjectEmailTemplate   string
    KindDescription        string
    SourceDescription      string
    ClientID               string
    MessageID              string
    Space                  string
    Organization           string
}

func NewMessageContext(user uaa.User, params NotifyParams,
    env config.Environment, space, organization, clientID string, guidGenerator GUIDGenerationFunc,
    plainTextEmailTemplate, htmlEmailTemplate, subjectEmailTemplate string) MessageContext {

    guid, err := guidGenerator()
    if err != nil {
        panic(err)
    }

    var kindDescription string
    if params.KindDescription == "" {
        kindDescription = params.Kind
    } else {
        kindDescription = params.KindDescription
    }

    var sourceDescription string
    if params.SourceDescription == "" {
        sourceDescription = clientID
    } else {
        sourceDescription = params.SourceDescription
    }

    return MessageContext{
        From:    env.Sender,
        ReplyTo: params.ReplyTo,
        To:      user.Emails[0],
        Subject: params.Subject,
        Text:    params.Text,
        HTML:    params.HTML,
        PlainTextEmailTemplate: plainTextEmailTemplate,
        HTMLEmailTemplate:      htmlEmailTemplate,
        SubjectEmailTemplate:   subjectEmailTemplate,
        KindDescription:        kindDescription,
        SourceDescription:      sourceDescription,
        ClientID:               clientID,
        MessageID:              guid.String(),
        Space:                  space,
        Organization:           organization,
    }
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
    var plainText string
    var err error

    headerPart := "\nThis is a multi-part message in MIME format...\n\n"

    plainTextPart := ""
    htmlPart := ""
    closingBoundary := "--our-content-boundary--"

    if sender.context.Text != "" {
        plainText, err = sender.compileTemplate(sender.context.PlainTextEmailTemplate, false)
        if err != nil {
            return "", err
        }

        plainTextPart = fmt.Sprintf("--our-content-boundary\nContent-type: text/plain\n\n%s\n", plainText)
    }

    var html string
    if sender.context.HTML != "" {
        html, err = sender.compileTemplate(sender.context.HTMLEmailTemplate, true)
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
`, html)
    }

    return headerPart + plainTextPart + htmlPart + closingBoundary, nil
}

func (sender MailSender) compileTemplate(theTemplate string, escapeContext bool) (string, error) {
    buffer := bytes.NewBuffer([]byte{})

    source, err := template.New("compileTemplate").Parse(theTemplate)
    if err != nil {
        return "", err
    }

    if escapeContext {
        source.Execute(buffer, sender.escapeContext(sender.context))
    } else {
        source.Execute(buffer, sender.context)
    }

    compiledTemplate := strings.TrimSuffix(buffer.String(), "\n")

    return compiledTemplate, nil
}

func (sender MailSender) escapeContext(context MessageContext) MessageContext {
    return MessageContext{
        From:    html.EscapeString(context.From),
        To:      html.EscapeString(context.To),
        Subject: html.EscapeString(context.Subject),
        Text:    html.EscapeString(context.Text),
        HTML:    context.HTML,
        PlainTextEmailTemplate: context.PlainTextEmailTemplate,
        HTMLEmailTemplate:      context.HTMLEmailTemplate,
        SubjectEmailTemplate:   context.SubjectEmailTemplate,
        KindDescription:        html.EscapeString(context.KindDescription),
        SourceDescription:      html.EscapeString(context.SourceDescription),
        ClientID:               html.EscapeString(context.ClientID),
        MessageID:              html.EscapeString(context.MessageID),
        Space:                  html.EscapeString(context.Space),
        Organization:           html.EscapeString(context.Organization),
    }
}

func (sender MailSender) CompileMessage(body string) (mail.Message, error) {
    compiledSubject, err := sender.compileTemplate(sender.context.SubjectEmailTemplate, false)
    if err != nil {
        return mail.Message{}, err
    }

    return mail.Message{
        From:    sender.context.From,
        ReplyTo: sender.context.ReplyTo,
        To:      sender.context.To,
        Subject: compiledSubject,
        Body:    body,
        Headers: []string{
            fmt.Sprintf("X-CF-Client-ID: %s", sender.context.ClientID),
            fmt.Sprintf("X-CF-Notification-ID: %s", sender.context.MessageID),
        },
    }, nil
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
