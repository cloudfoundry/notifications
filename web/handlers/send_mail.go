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

func SendMail(client mail.ClientInterface, context MessageContext) (string, error) {
    sender := NewMailSender(client, context)

    compiledTemplate, err := sender.CompileTemplate()
    if err != nil {
        return "", err
    }

    message := sender.CompileMessage(compiledTemplate)
    return sender.Deliver(message), nil
}

type MessageContext struct {
    From              string
    To                string
    Subject           string
    Text              string
    Template          string
    KindDescription   string
    SourceDescription string
    ClientID          string
    MessageID         string
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

func (sender MailSender) CompileTemplate() (*bytes.Buffer, error) {
    buffer := bytes.NewBuffer([]byte{})
    source, err := template.New("emailBody").Parse(sender.context.Template)
    if err != nil {
        return buffer, err
    }

    source.Execute(buffer, sender.context)
    return buffer, nil
}

func (sender MailSender) CompileMessage(buffer *bytes.Buffer) mail.Message {
    return mail.Message{
        From:    sender.context.From,
        To:      sender.context.To,
        Subject: fmt.Sprintf("CF Notification: %s", sender.context.Subject),
        Body:    buffer.String(),
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
