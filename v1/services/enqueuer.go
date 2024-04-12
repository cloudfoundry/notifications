package services

import (
	"time"

	"gopkg.in/gorp.v1"

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

type queueInterface interface {
	Enqueue(job *gobble.Job, transaction gobble.ConnectionInterface) (*gobble.Job, error)
}

type gobbleInitializer interface {
	InitializeDBMap(*gorp.DbMap)
}

type Enqueuer struct {
	queue             queueInterface
	messagesRepo      messagesRepoUpserter
	gobbleInitializer gobbleInitializer
}

func NewEnqueuer(queue queueInterface, messagesRepo messagesRepoUpserter, gobbleInitializer gobbleInitializer) Enqueuer {
	return Enqueuer{
		queue:             queue,
		messagesRepo:      messagesRepo,
		gobbleInitializer: gobbleInitializer,
	}
}

func (enqueuer Enqueuer) Enqueue(
	conn ConnectionInterface,
	users []User,
	options Options,
	space cf.CloudControllerSpace,
	organization cf.CloudControllerOrganization,
	clientID,
	uaaHost,
	scope,
	vcapRequestID string,
	reqReceived time.Time) ([]Response, error) {

	var responses []Response

	transaction := conn.Transaction()
	enqueuer.gobbleInitializer.InitializeDBMap(transaction.GetDbMap())

	if err := transaction.Begin(); err != nil {
		return []Response{}, err
	}

	for _, user := range users {
		message, err := enqueuer.messagesRepo.Upsert(transaction, models.Message{
			Status: StatusQueued,
		})
		if err != nil {
			transaction.Rollback()
			return []Response{}, err
		}

		job := gobble.NewJob(Delivery{
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
		})

		_, err = enqueuer.queue.Enqueue(job, transaction)
		if err != nil {
			transaction.Rollback()
			return []Response{}, err
		}

		//TODO: don't append to responses if job returned is nil?

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

	if err := transaction.Commit(); err != nil {
		return []Response{}, err
	}

	return responses, nil
}
