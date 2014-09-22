package postal

import (
    "github.com/cloudfoundry-incubator/notifications/gobble"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type MailerInterface interface {
    Deliver(models.ConnectionInterface, Templates, map[string]uaa.User, Options, string, string, string) []Response
}

type Mailer struct {
    queue            gobble.QueueInterface
    guidGenerator    GUIDGenerationFunc
    unsubscribesRepo models.UnsubscribesRepoInterface
    kindsRepo        models.KindsRepoInterface
}

func NewMailer(queue gobble.QueueInterface, guidGenerator GUIDGenerationFunc, unsubscribesRepo models.UnsubscribesRepoInterface, kindsRepo models.KindsRepoInterface) Mailer {
    return Mailer{
        queue:            queue,
        guidGenerator:    guidGenerator,
        unsubscribesRepo: unsubscribesRepo,
        kindsRepo:        kindsRepo,
    }
}

func (mailer Mailer) Deliver(conn models.ConnectionInterface, templates Templates, users map[string]uaa.User,
    options Options, space, organization, clientID string) []Response {
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
        if (err == models.ErrRecordNotFound{}) || mailer.isCritical(conn, mailer.kindsRepo, options.KindID, clientID) {
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

func (mailer Mailer) isCritical(conn models.ConnectionInterface, kindsRepo models.KindsRepoInterface, kindID, clientID string) bool {

    kind, err := kindsRepo.Find(conn, kindID, clientID)
    if (err == models.ErrRecordNotFound{}) {
        return false
    }

    return kind.Critical
}
