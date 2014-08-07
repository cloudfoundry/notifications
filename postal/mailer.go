package postal

import "github.com/pivotal-cf/uaa-sso-golang/uaa"

type Mailer struct {
    queue *DeliveryQueue
}

func NewMailer(queue *DeliveryQueue) Mailer {
    return Mailer{
        queue: queue,
    }
}

func (mailer Mailer) Deliver(templates Templates, users map[string]uaa.User, options Options, space, organization, clientID string) []Response {
    responses := []Response{}
    deliveryResponses := make(chan Response)

    for userGUID, user := range users {
        delivery := Delivery{
            User:         user,
            Options:      options,
            UserGUID:     userGUID,
            Space:        space,
            Organization: organization,
            ClientID:     clientID,
            Templates:    templates,
            Response:     deliveryResponses,
        }
        mailer.queue.Enqueue(delivery)
    }

    for _, _ = range users {
        responses = append(responses, <-deliveryResponses)
    }

    return responses
}
