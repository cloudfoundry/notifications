package postal

import (
    "bytes"
    "fmt"
    "strings"
    "text/template"

    "github.com/cloudfoundry-incubator/notifications/mail"
)

type Packager struct{}

func NewPackager() Packager {
    return Packager{}
}

func (packager Packager) Pack(context MessageContext) (mail.Message, error) {
    body, err := packager.CompileBody(context)
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
        Body:    body,
        Headers: []string{
            fmt.Sprintf("X-CF-Client-ID: %s", context.ClientID),
            fmt.Sprintf("X-CF-Notification-ID: %s", context.MessageID),
        },
    }, nil
}

func (packager Packager) CompileBody(context MessageContext) (string, error) {
    var plainText string
    var err error

    headerPart := "\nThis is a multi-part message in MIME format...\n\n"
    plainTextPart := ""
    htmlPart := ""
    closingBoundary := "--our-content-boundary--"

    if context.Text != "" {
        plainText, err = packager.compileTemplate(context, context.TextTemplate, false)
        if err != nil {
            return "", err
        }

        plainTextPart = fmt.Sprintf("--our-content-boundary\nContent-type: text/plain\n\n%s\n", plainText)
    }

    var html string
    if context.HTML != "" {
        html, err = packager.compileTemplate(context, context.HTMLTemplate, true)
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

func (packager Packager) compileTemplate(context MessageContext, theTemplate string, escapeContext bool) (string, error) {
    buffer := bytes.NewBuffer([]byte{})

    source, err := template.New("compileTemplate").Parse(theTemplate)
    if err != nil {
        return "", err
    }

    if escapeContext {
        source.Execute(buffer, context.Escape())
    } else {
        source.Execute(buffer, context)
    }

    compiledTemplate := strings.TrimSuffix(buffer.String(), "\n")

    return compiledTemplate, nil
}
