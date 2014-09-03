package postal

import (
    "github.com/cloudfoundry-incubator/notifications/gobble"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type Mailer struct {
    queue            gobble.QueueInterface
    guidGenerator    GUIDGenerationFunc
    unsubscribesRepo models.UnsubscribesRepoInterface
}

func NewMailer(queue gobble.QueueInterface, guidGenerator GUIDGenerationFunc, unsubscribesRepo models.UnsubscribesRepoInterface) Mailer {
    return Mailer{
        queue:            queue,
        guidGenerator:    guidGenerator,
        unsubscribesRepo: unsubscribesRepo,
    }
}

func (mailer Mailer) Deliver(conn models.ConnectionInterface, templates Templates, users map[string]uaa.User, options Options, space, organization, clientID string) []Response {
    responses := []Response{}
    transaction := conn.Transaction()

    transaction.Begin()
    for userGUID, user := range users {
        guid, err := mailer.guidGenerator()
        if err != nil {
            panic(err)
        }
        messageID := guid.String()

        _, err = mailer.unsubscribesRepo.Find(conn, clientID, options.KindID, userGUID)
        if (err == models.ErrRecordNotFound{}) {
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
        }

        responses = append(responses, Response{
            Status:         StatusQueued,
            Recipient:      userGUID,
            NotificationID: messageID,
        })
    }

    transaction.Commit()
    return responses
}
