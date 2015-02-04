package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
)

type MailerInterface interface {
	Deliver(models.ConnectionInterface, []User, postal.Options, cf.CloudControllerSpace, cf.CloudControllerOrganization, string, string) []Response
}

type Mailer struct {
	queue         gobble.QueueInterface
	guidGenerator postal.GUIDGenerationFunc
	messagesRepo  MessagesRepoInterface
}

type MessagesRepoInterface interface {
	Upsert(models.ConnectionInterface, models.Message) (models.Message, error)
}

func NewMailer(queue gobble.QueueInterface, guidGenerator postal.GUIDGenerationFunc, messagesRepo MessagesRepoInterface) Mailer {
	return Mailer{
		queue:         queue,
		guidGenerator: guidGenerator,
		messagesRepo:  messagesRepo,
	}
}

func (mailer Mailer) Deliver(conn models.ConnectionInterface, users []User,
	options postal.Options, space cf.CloudControllerSpace,
	organization cf.CloudControllerOrganization, clientID, scope string) []Response {

	responses := []Response{}
	jobsByMessageID := map[string]gobble.Job{}
	for _, user := range users {
		guid, err := mailer.guidGenerator()
		if err != nil {
			panic(err)
		}
		messageID := guid.String()

		jobsByMessageID[messageID] = gobble.NewJob(postal.Delivery{
			Options:      options,
			UserGUID:     user.GUID,
			Email:        user.Email,
			Space:        space,
			Organization: organization,
			ClientID:     clientID,
			MessageID:    messageID,
			Scope:        scope,
		})

		recipient := user.Email
		if recipient == "" {
			recipient = user.GUID
		}

		responses = append(responses, Response{
			Status:         postal.StatusQueued,
			NotificationID: messageID,
			Recipient:      recipient,
		})
	}

	transaction := conn.Transaction()
	transaction.Begin()
	for messageID, job := range jobsByMessageID {
		_, err := mailer.queue.Enqueue(job)
		if err != nil {
			transaction.Rollback()
			return []Response{}
		}
		_, err = mailer.messagesRepo.Upsert(transaction, models.Message{
			ID:     messageID,
			Status: postal.StatusQueued,
		})
		if err != nil {
			transaction.Rollback()
			return []Response{}
		}
	}
	err := transaction.Commit()
	if err != nil {
		return []Response{}
	}

	return responses
}
