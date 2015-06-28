package strategies

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/postal"
)

type EnqueuerInterface interface {
	Enqueue(models.ConnectionInterface, []User, postal.Options, cf.CloudControllerSpace, cf.CloudControllerOrganization, string, string, string, time.Time) []Response
}

type Enqueuer struct {
	queue         gobble.QueueInterface
	guidGenerator postal.GUIDGenerationFunc
	messagesRepo  MessagesRepoInterface
}

type MessagesRepoInterface interface {
	Upsert(models.ConnectionInterface, models.Message) (models.Message, error)
}

func NewEnqueuer(queue gobble.QueueInterface, guidGenerator postal.GUIDGenerationFunc, messagesRepo MessagesRepoInterface) Enqueuer {
	return Enqueuer{
		queue:         queue,
		guidGenerator: guidGenerator,
		messagesRepo:  messagesRepo,
	}
}

func (enqueuer Enqueuer) Enqueue(conn models.ConnectionInterface, users []User,
	options postal.Options, space cf.CloudControllerSpace,
	organization cf.CloudControllerOrganization, clientID, scope, vcapRequestID string, reqReceived time.Time) []Response {

	responses := []Response{}
	jobsByMessageID := map[string]gobble.Job{}
	for _, user := range users {
		guid, err := enqueuer.guidGenerator()
		if err != nil {
			panic(err)
		}
		messageID := guid.String()

		jobsByMessageID[messageID] = gobble.NewJob(postal.Delivery{
			Options:         options,
			UserGUID:        user.GUID,
			Email:           user.Email,
			Space:           space,
			Organization:    organization,
			ClientID:        clientID,
			MessageID:       messageID,
			Scope:           scope,
			VCAPRequestID:   vcapRequestID,
			RequestReceived: reqReceived,
		})

		recipient := user.Email
		if recipient == "" {
			recipient = user.GUID
		}

		responses = append(responses, Response{
			Status:         postal.StatusQueued,
			NotificationID: messageID,
			Recipient:      recipient,
			VCAPRequestID:  vcapRequestID,
		})
	}

	transaction := conn.Transaction()
	transaction.Begin()
	for messageID := range jobsByMessageID {
		_, err := enqueuer.messagesRepo.Upsert(transaction, models.Message{
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

	for _, job := range jobsByMessageID {
		_, err := enqueuer.queue.Enqueue(job)
		if err != nil {
			return []Response{}
		}
	}

	return responses
}
