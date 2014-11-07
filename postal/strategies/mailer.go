package strategies

import (
    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/gobble"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type MailerInterface interface {
    Deliver(models.ConnectionInterface, postal.Templates, map[string]uaa.User, postal.Options, cf.CloudControllerSpace, cf.CloudControllerOrganization, string) []Response
}

type Mailer struct {
    queue         gobble.QueueInterface
    guidGenerator postal.GUIDGenerationFunc
}

func NewMailer(queue gobble.QueueInterface, guidGenerator postal.GUIDGenerationFunc) Mailer {
    return Mailer{
        queue:         queue,
        guidGenerator: guidGenerator,
    }
}

func (mailer Mailer) Deliver(conn models.ConnectionInterface, templates postal.Templates, users map[string]uaa.User,
    options postal.Options, space cf.CloudControllerSpace, organization cf.CloudControllerOrganization, clientID string) []Response {

    responses := []Response{}
    transaction := conn.Transaction()

    transaction.Begin()
    for userGUID, user := range users {
        guid, err := mailer.guidGenerator()
        if err != nil {
            panic(err)
        }
        messageID := guid.String()

        job := gobble.NewJob(postal.Delivery{
            User:         user,
            Options:      options,
            UserGUID:     userGUID,
            Space:        space,
            Organization: organization,
            ClientID:     clientID,
            Templates:    templates,
            MessageID:    messageID,
        })

        _, err = mailer.queue.Enqueue(job)
        if err != nil {
            transaction.Rollback()
            return []Response{}
        }

        emailAddress := ""
        if len(user.Emails) > 0 {
            emailAddress = user.Emails[0]
        }

        responses = append(responses, Response{
            Status:         postal.StatusQueued,
            Recipient:      userGUID,
            NotificationID: messageID,
            Email:          emailAddress,
        })
    }

    transaction.Commit()
    return responses
}
