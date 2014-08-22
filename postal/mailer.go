package postal

import "github.com/pivotal-cf/uaa-sso-golang/uaa"

type Mailer struct {
    queue         *DeliveryQueue
    guidGenerator GUIDGenerationFunc
}

func NewMailer(queue *DeliveryQueue, guidGenerator GUIDGenerationFunc) Mailer {
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

        delivery := Delivery{
            User:         user,
            Options:      options,
            UserGUID:     userGUID,
            Space:        space,
            Organization: organization,
            ClientID:     clientID,
            Templates:    templates,
            MessageID:    messageID,
        }
        mailer.queue.Enqueue(delivery)
        responses = append(responses, Response{
            Status:         StatusQueued,
            Recipient:      userGUID,
            NotificationID: messageID,
        })
    }

    return responses
}
