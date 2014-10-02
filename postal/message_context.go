package postal

import (
    "html"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/cryptography"
)

type MessageContext struct {
    From              string
    ReplyTo           string
    To                string
    Subject           string
    Text              string
    HTML              string
    HTMLComponents    HTML
    TextTemplate      string
    HTMLTemplate      string
    SubjectTemplate   string
    KindDescription   string
    SourceDescription string
    UserGUID          string
    ClientID          string
    MessageID         string
    Space             string
    Organization      string
    UnsubscribeID     string
}

func NewMessageContext(email string, options Options, env config.Environment, space, organization,
    clientID string, messageID string, userGUID string, templates Templates, cryptoClient cryptography.CryptoInterface) MessageContext {

    var kindDescription string
    if options.KindDescription == "" {
        kindDescription = options.KindID
    } else {
        kindDescription = options.KindDescription
    }

    var sourceDescription string
    if options.SourceDescription == "" {
        sourceDescription = clientID
    } else {
        sourceDescription = options.SourceDescription
    }

    messageContext := MessageContext{
        From:              env.Sender,
        ReplyTo:           options.ReplyTo,
        To:                email,
        Subject:           options.Subject,
        Text:              options.Text,
        HTML:              options.HTML.BodyContent,
        HTMLComponents:    options.HTML,
        TextTemplate:      templates.Text,
        HTMLTemplate:      templates.HTML,
        SubjectTemplate:   templates.Subject,
        KindDescription:   kindDescription,
        SourceDescription: sourceDescription,
        UserGUID:          userGUID,
        ClientID:          clientID,
        MessageID:         messageID,
        Space:             space,
        Organization:      organization,
    }

    var err error
    messageContext.UnsubscribeID, err = cryptoClient.Encrypt(userGUID + "|" + clientID + "|" + options.KindID)
    if err != nil {
        panic(err)
    }
    return messageContext
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
