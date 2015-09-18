package services

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

const StatusQueued = "queued"

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
	TemplateID        string
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

type messagesRepoUpserter interface {
	Upsert(models.ConnectionInterface, models.Message) (models.Message, error)
}

type Enqueuer struct {
	queue        gobble.QueueInterface
	messagesRepo messagesRepoUpserter
}

func NewEnqueuer(queue gobble.QueueInterface, messagesRepo messagesRepoUpserter) Enqueuer {
	return Enqueuer{
		queue:        queue,
		messagesRepo: messagesRepo,
	}
}

func (enqueuer Enqueuer) Enqueue(conn ConnectionInterface, users []User, options Options, space cf.CloudControllerSpace, organization cf.CloudControllerOrganization, clientID, uaaHost, scope, vcapRequestID string, reqReceived time.Time) []Response {
	var (
		responses []Response
		jobs      []gobble.Job
	)

	transaction := conn.Transaction()
	transaction.Begin()
	for _, user := range users {
		message, err := enqueuer.messagesRepo.Upsert(transaction, models.Message{
			Status: StatusQueued,
		})
		if err != nil {
			transaction.Rollback()
			return []Response{}
		}

		jobs = append(jobs, gobble.NewJob(Delivery{
			Options:         options,
			UserGUID:        user.GUID,
			Email:           user.Email,
			Space:           space,
			Organization:    organization,
			ClientID:        clientID,
			MessageID:       message.ID,
			UAAHost:         uaaHost,
			Scope:           scope,
			VCAPRequestID:   vcapRequestID,
			RequestReceived: reqReceived,
		}))

		recipient := user.Email
		if recipient == "" {
			recipient = user.GUID
		}

		responses = append(responses, Response{
			Status:         message.Status,
			NotificationID: message.ID,
			Recipient:      recipient,
			VCAPRequestID:  vcapRequestID,
		})
	}

	err := transaction.Commit()
	if err != nil {
		return []Response{}
	}

	for _, job := range jobs {
		_, err := enqueuer.queue.Enqueue(job)
		if err != nil {
			return []Response{}
		}
	}

	return responses
}
