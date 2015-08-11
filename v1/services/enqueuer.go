package services

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/nu7hatch/gouuid"
)

const StatusQueued = "queued"

type GUIDGenerationFunc func() (*uuid.UUID, error)

type Options struct {
	ReplyTo           string
	Subject           string
	KindDescription   string
	SourceDescription string
	Text              string
	HTML              HTML
	KindID            string
	To                string
	Role              string
	Endorsement       string
}

type Delivery struct {
	MessageID       string
	Options         Options
	UserGUID        string
	Email           string
	Space           cf.CloudControllerSpace
	Organization    cf.CloudControllerOrganization
	ClientID        string
	UAAHost         string
	Scope           string
	VCAPRequestID   string
	RequestReceived time.Time
}

type EnqueuerInterface interface {
	Enqueue(db.ConnectionInterface, []User, Options, cf.CloudControllerSpace, cf.CloudControllerOrganization, string, string, string, string, time.Time) []Response
}

type Enqueuer struct {
	queue         gobble.QueueInterface
	guidGenerator GUIDGenerationFunc
	messagesRepo  MessagesRepoInterface
}

func NewEnqueuer(queue gobble.QueueInterface, guidGenerator GUIDGenerationFunc, messagesRepo MessagesRepoInterface) Enqueuer {
	return Enqueuer{
		queue:         queue,
		guidGenerator: guidGenerator,
		messagesRepo:  messagesRepo,
	}
}

func (enqueuer Enqueuer) Enqueue(conn db.ConnectionInterface, users []User, options Options, space cf.CloudControllerSpace, organization cf.CloudControllerOrganization, clientID, uaaHost, scope, vcapRequestID string, reqReceived time.Time) []Response {

	responses := []Response{}
	jobsByMessageID := map[string]gobble.Job{}
	for _, user := range users {
		guid, err := enqueuer.guidGenerator()
		if err != nil {
			panic(err)
		}
		messageID := guid.String()

		jobsByMessageID[messageID] = gobble.NewJob(Delivery{
			Options:         options,
			UserGUID:        user.GUID,
			Email:           user.Email,
			Space:           space,
			Organization:    organization,
			ClientID:        clientID,
			MessageID:       messageID,
			UAAHost:         uaaHost,
			Scope:           scope,
			VCAPRequestID:   vcapRequestID,
			RequestReceived: reqReceived,
		})

		recipient := user.Email
		if recipient == "" {
			recipient = user.GUID
		}

		responses = append(responses, Response{
			Status:         StatusQueued,
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
			Status: StatusQueued,
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
