package postal

import (
    "log"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const (
    StatusUnavailable = "unavailable"
    StatusFailed      = "failed"
    StatusDelivered   = "delivered"
    StatusNotFound    = "notfound"
    StatusNoAddress   = "noaddress"
)

type Mailer struct {
    guidGenerator GUIDGenerationFunc
    logger        *log.Logger
    mailClient    mail.ClientInterface
}

func NewMailer(guidGenerator GUIDGenerationFunc, logger *log.Logger, mailClient mail.ClientInterface) Mailer {
    return Mailer{
        guidGenerator: guidGenerator,
        logger:        logger,
        mailClient:    mailClient,
    }
}

func (mailer Mailer) Deliver(templates Templates, users map[string]uaa.User, options Options, space, organization, clientID string) []Response {
    env := config.NewEnvironment()
    messages := []Response{}

    for userGUID, user := range users {
        var status, notificationID string
        if len(user.Emails) > 0 && strings.Contains(user.Emails[0], "@") {
            context := NewMessageContext(user.Emails[0], options, env, space, organization, clientID, mailer.guidGenerator, templates)
            status = mailer.SendMailToUser(context, mailer.logger, mailer.mailClient)
            notificationID = context.MessageID
        } else {
            if user.ID == "" {
                status = StatusNotFound
            } else {
                status = StatusNoAddress
            }
        }

        messages = append(messages, Response{
            Status:         status,
            Recipient:      userGUID,
            NotificationID: notificationID,
        })
    }
    return messages
}

func (mailer Mailer) SendMailToUser(context MessageContext, logger *log.Logger, mailClient mail.ClientInterface) string {
    logger.Printf("Sending email to %s", context.To)
    packager := NewPackager()

    message, err := packager.Pack(context)
    if err != nil {
        panic(err)
    }

    status := mailer.SendMail(message)

    logger.Print(message.Data())
    return status
}

func (mailer Mailer) SendMail(msg mail.Message) string {
    err := mailer.mailClient.Connect()
    if err != nil {
        return StatusUnavailable
    }

    err = mailer.mailClient.Send(msg)
    if err != nil {
        return StatusFailed
    }

    return StatusDelivered
}
