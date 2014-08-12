package postal

import (
    "html"

    "github.com/cloudfoundry-incubator/notifications/config"
)

type MessageContext struct {
    From              string
    ReplyTo           string
    To                string
    Subject           string
    Text              string
    HTML              string
    TextTemplate      string
    HTMLTemplate      string
    SubjectTemplate   string
    KindDescription   string
    SourceDescription string
    ClientID          string
    MessageID         string
    Space             string
    Organization      string
}

func NewMessageContext(email string, options Options, env config.Environment, space, organization,
    clientID string, guidGenerator GUIDGenerationFunc, templates Templates) MessageContext {

    guid, err := guidGenerator()
    if err != nil {
        panic(err)
    }

    var kindDescription string
    if options.KindDescription == "" {
        kindDescription = options.Kind
    } else {
        kindDescription = options.KindDescription
    }

    var sourceDescription string
    if options.SourceDescription == "" {
        sourceDescription = clientID
    } else {
        sourceDescription = options.SourceDescription
    }

    return MessageContext{
        From:              env.Sender,
        ReplyTo:           options.ReplyTo,
        To:                email,
        Subject:           options.Subject,
        Text:              options.Text,
        HTML:              options.HTML,
        TextTemplate:      templates.Text,
        HTMLTemplate:      templates.HTML,
        SubjectTemplate:   templates.Subject,
        KindDescription:   kindDescription,
        SourceDescription: sourceDescription,
        ClientID:          clientID,
        MessageID:         guid.String(),
        Space:             space,
        Organization:      organization,
    }
}

func (context *MessageContext) Escape() {
    context.From = html.EscapeString(context.From)
    context.To = html.EscapeString(context.To)
    context.ReplyTo = html.EscapeString(context.ReplyTo)
    context.Subject = html.EscapeString(context.Subject)
    context.Text = html.EscapeString(context.Text)
    context.KindDescription = html.EscapeString(context.KindDescription)
    context.SourceDescription = html.EscapeString(context.SourceDescription)
    context.ClientID = html.EscapeString(context.ClientID)
    context.MessageID = html.EscapeString(context.MessageID)
    context.Space = html.EscapeString(context.Space)
    context.Organization = html.EscapeString(context.Organization)
}
