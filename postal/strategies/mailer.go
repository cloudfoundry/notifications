package strategies

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"
)

type MailerInterface interface {
	Deliver(models.ConnectionInterface, postal.Templates, map[string]uaa.User, postal.Options, cf.CloudControllerSpace, cf.CloudControllerOrganization, string, string) []Response
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

func (mailer Mailer) Deliver(conn models.ConnectionInterface, templates postal.Templates, users map[string]uaa.User,
	options postal.Options, space cf.CloudControllerSpace, organization cf.CloudControllerOrganization, clientID, scope string) []Response {

	responses := []Response{}
	jobsByMessageID := map[string]gobble.Job{}
	for userGUID, user := range users {
		guid, err := mailer.guidGenerator()
		if err != nil {
			panic(err)
		}
		messageID := guid.String()

		jobsByMessageID[messageID] = gobble.NewJob(postal.Delivery{
			User:         user,
			Options:      options,
			UserGUID:     userGUID,
			Space:        space,
			Organization: organization,
			ClientID:     clientID,
			Templates:    templates,
			MessageID:    messageID,
			Scope:        scope,
		})

		responses = append(responses, Response{
			Status:         postal.StatusQueued,
			NotificationID: messageID,
			Recipient:      userGUID,
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
			Status: "queued",
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
