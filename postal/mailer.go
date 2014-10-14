package postal

import (
    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/gobble"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type MailerInterface interface {
    Deliver(models.ConnectionInterface, Templates, map[string]uaa.User, Options, cf.CloudControllerSpace, cf.CloudControllerOrganization, string) []Response
}

type Mailer struct {
    queue         gobble.QueueInterface
    guidGenerator GUIDGenerationFunc
}

func NewMailer(queue gobble.QueueInterface, guidGenerator GUIDGenerationFunc) Mailer {
    return Mailer{
        queue:         queue,
        guidGenerator: guidGenerator,
    }
}

func (mailer Mailer) Deliver(conn models.ConnectionInterface, templates Templates, users map[string]uaa.User,
    options Options, space cf.CloudControllerSpace, organization cf.CloudControllerOrganization, clientID string) []Response {

    responses := []Response{}
    transaction := conn.Transaction()

    transaction.Begin()
    for userGUID, user := range users {
        guid, err := mailer.guidGenerator()
        if err != nil {
            panic(err)
        }
        messageID := guid.String()

        job := gobble.NewJob(Delivery{
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
            Status:         StatusQueued,
            Recipient:      userGUID,
            NotificationID: messageID,
            Email:          emailAddress,
        })
    }

    transaction.Commit()
    return responses
}
