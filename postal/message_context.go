package postal

import (
	"html"

	"github.com/pivotal-golang/conceal"
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
	SpaceGUID         string
	Organization      string
	OrganizationGUID  string
	UnsubscribeID     string
	Scope             string
}

func NewMessageContext(delivery Delivery, sender string, cloak conceal.CloakInterface) MessageContext {
	options := delivery.Options
	templates := delivery.Templates

	var kindDescription string
	if options.KindDescription == "" {
		kindDescription = options.KindID
	} else {
		kindDescription = options.KindDescription
	}

	var sourceDescription string
	if options.SourceDescription == "" {
		sourceDescription = delivery.ClientID
	} else {
		sourceDescription = options.SourceDescription
	}

	messageContext := MessageContext{
		From:              sender,
		ReplyTo:           options.ReplyTo,
		To:                delivery.User.Emails[0],
		Subject:           options.Subject,
		Text:              options.Text,
		HTML:              options.HTML.BodyContent,
		HTMLComponents:    options.HTML,
		TextTemplate:      templates.Text,
		HTMLTemplate:      templates.HTML,
		SubjectTemplate:   templates.Subject,
		KindDescription:   kindDescription,
		SourceDescription: sourceDescription,
		UserGUID:          delivery.UserGUID,
		ClientID:          delivery.ClientID,
		MessageID:         delivery.MessageID,
		Space:             delivery.Space.Name,
		SpaceGUID:         delivery.Space.GUID,
		Organization:      delivery.Organization.Name,
		OrganizationGUID:  delivery.Organization.GUID,
		Scope:             delivery.Scope,
	}

	if messageContext.Subject == "" {
		messageContext.Subject = "[no subject]"
	}

	unsubscribeID, err := cloak.Veil([]byte(delivery.UserGUID + "|" + delivery.ClientID + "|" + options.KindID))
	if err != nil {
		panic(err)
	}

	messageContext.UnsubscribeID = string(unsubscribeID)
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
