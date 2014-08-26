package postal

import (
    "github.com/cloudfoundry-incubator/notifications/gobble"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

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

func (mailer Mailer) Deliver(templates Templates, users map[string]uaa.User, options Options, space, organization, clientID string) []Response {
    responses := []Response{}

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

        mailer.queue.Enqueue(job)

        responses = append(responses, Response{
            Status:         StatusQueued,
            Recipient:      userGUID,
            NotificationID: messageID,
        })
    }

    return responses
}
