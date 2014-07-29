package postal

import (
    "log"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const (
    MailServerUnavailable  = "unavailable"
    MailDeliveryFailed     = "failed"
    MailDeliverySuccessful = "delivered"
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

    for userGUID, uaaUser := range users {
        if len(uaaUser.Emails) > 0 {
            context := NewMessageContext(uaaUser, options, env, space, organization,
                clientID, mailer.guidGenerator, templates)

            emailStatus := mailer.SendMailToUser(context, mailer.logger, mailer.mailClient)
            mailer.logger.Println(emailStatus)

            mailInfo := Response{
                Status:         emailStatus,
                Recipient:      uaaUser.ID,
                NotificationID: context.MessageID,
            }

            messages = append(messages, mailInfo)
        } else {
            var status string
            if uaaUser.ID == "" {
                status = StatusNotFound
            } else {
                status = StatusNoAddress
            }
            mailInfo := Response{
                Status:         status,
                Recipient:      userGUID,
                NotificationID: "",
            }

            messages = append(messages, mailInfo)
        }
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
        return MailServerUnavailable
    }

    err = mailer.mailClient.Send(msg)
    if err != nil {
        return MailDeliveryFailed
    }

    return MailDeliverySuccessful
}
