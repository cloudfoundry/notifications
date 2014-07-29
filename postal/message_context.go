package postal

import (
    "html"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
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

func NewMessageContext(user uaa.User, options Options, env config.Environment, space, organization,
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
        To:                user.Emails[0],
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

func (context MessageContext) Escape() MessageContext {
    return MessageContext{
        From:              html.EscapeString(context.From),
        To:                html.EscapeString(context.To),
        Subject:           html.EscapeString(context.Subject),
        Text:              html.EscapeString(context.Text),
        HTML:              context.HTML,
        TextTemplate:      context.TextTemplate,
        HTMLTemplate:      context.HTMLTemplate,
        SubjectTemplate:   context.SubjectTemplate,
        KindDescription:   html.EscapeString(context.KindDescription),
        SourceDescription: html.EscapeString(context.SourceDescription),
        ClientID:          html.EscapeString(context.ClientID),
        MessageID:         html.EscapeString(context.MessageID),
        Space:             html.EscapeString(context.Space),
        Organization:      html.EscapeString(context.Organization),
    }
}
