package common

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/cloudfoundry-incubator/notifications/mail"
	"github.com/pivotal-golang/conceal"
)

const HTMLWrapperTemplate = `{{.HTMLComponents.Doctype}}
<head>{{.HTMLComponents.Head}}</head>
<html>
	<body {{.HTMLComponents.BodyAttributes}}>
		{{.HTMLComponents.BodyContent}}
	</body>
</html>`

type templatesLoader interface {
	LoadTemplates(clientID, kindID, templateID string) (Templates, error)
}

type Packager struct {
	templates templatesLoader
	cloak     conceal.CloakInterface
}

func NewPackager(templates templatesLoader, cloak conceal.CloakInterface) Packager {
	return Packager{
		templates: templates,
		cloak:     cloak,
	}
}

func (packager Packager) PrepareContext(delivery Delivery, sender, domain string) (MessageContext, error) {
	templates, err := packager.templates.LoadTemplates(delivery.ClientID, delivery.Options.KindID, delivery.Options.TemplateID)
	if err != nil {
		return MessageContext{}, err
	}

	return NewMessageContext(delivery, sender, domain, packager.cloak, templates), nil
}

func (packager Packager) Pack(context MessageContext) (mail.Message, error) {
	parts, err := packager.CompileParts(context)
	if err != nil {
		return mail.Message{}, err
	}

	compiledSubject, err := packager.compileTemplate(context, context.SubjectTemplate, false)
	if err != nil {
		return mail.Message{}, err
	}

	return mail.Message{
		From:    context.From,
		ReplyTo: context.ReplyTo,
		To:      context.To,
		Subject: compiledSubject,
		Body:    parts,
		Headers: []string{
			fmt.Sprintf("X-CF-Client-ID: %s", context.ClientID),
			fmt.Sprintf("X-CF-Notification-ID: %s", context.MessageID),
			fmt.Sprintf("X-CF-Notification-Timestamp: %s", time.Now().Format(time.RFC3339Nano)),
			fmt.Sprintf("X-CF-Notification-Request-Received: %s", context.RequestReceived.Format(time.RFC3339Nano)),
		},
	}, nil
}

func (packager Packager) CompileParts(context MessageContext) ([]mail.Part, error) {
	var parts []mail.Part
	var err error

	context.Endorsement, err = packager.compileTemplate(context, context.Endorsement, false)
	if err != nil {
		return parts, err
	}

	if context.Text != "" {
		plainText, err := packager.compileTemplate(context, context.TextTemplate, false)
		if err != nil {
			return parts, err
		}

		parts = append(parts, mail.Part{
			ContentType: "text/plain",
			Content:     plainText,
		})

	}

	if context.HTML != "" {
		var err error

		context.HTMLComponents.BodyContent, err = packager.compileTemplate(context, context.HTMLTemplate, true)
		if err != nil {
			return parts, err
		}

		htmlPart, err := packager.compileTemplate(context, HTMLWrapperTemplate, true)
		if err != nil {
			return parts, err
		}

		parts = append(parts, mail.Part{
			ContentType: "text/html",
			Content:     htmlPart,
		})
	}

	return parts, nil
}

func (packager Packager) compileTemplate(context MessageContext, theTemplate string, escapeContext bool) (string, error) {
	buffer := bytes.NewBuffer([]byte{})

	source, err := template.New("compileTemplate").Parse(theTemplate)
	if err != nil {
		return "", err
	}

	if escapeContext {
		context.Escape()
	}

	source.Execute(buffer, context)
	compiledTemplate := strings.TrimSuffix(buffer.String(), "\n")

	return compiledTemplate, nil
}
