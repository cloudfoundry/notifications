package queue

import (
	"time"

	"gopkg.in/gorp.v1"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
)

const StatusQueued = "queued"

type User struct {
	GUID        string
	Email       string
	Endorsement string
}

type Response struct {
	Status         string `json:"status"`
	Recipient      string `json:"recipient"`
	NotificationID string `json:"notification_id"`
	VCAPRequestID  string `json:"vcap_request_id"`
}

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

type HTML struct {
	BodyContent    string
	BodyAttributes string
	Head           string
	Doctype        string
}

type Delivery struct {
	JobType         string
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
	CampaignID      string
}

type messagesRepoInserter interface {
	Insert(models.ConnectionInterface, models.Message) (models.Message, error)
}

type gobbleInitializer interface {
	InitializeDBMap(*gorp.DbMap)
}

type JobEnqueuer struct {
	queue             enqueuer
	messagesRepo      messagesRepoInserter
	gobbleInitializer gobbleInitializer
}

func NewJobEnqueuer(queue enqueuer, messagesRepo messagesRepoInserter, gobbleInitializer gobbleInitializer) JobEnqueuer {
	return JobEnqueuer{
		queue:             queue,
		messagesRepo:      messagesRepo,
		gobbleInitializer: gobbleInitializer,
	}
}

func (enqueuer JobEnqueuer) Enqueue(conn ConnectionInterface, users []User, options Options, space cf.CloudControllerSpace, organization cf.CloudControllerOrganization, clientID, uaaHost, scope, vcapRequestID string, reqReceived time.Time, campaignID string) {
	transaction := conn.Transaction()
	enqueuer.gobbleInitializer.InitializeDBMap(transaction.GetDbMap())

	transaction.Begin()

	for _, user := range users {
		message, err := enqueuer.messagesRepo.Insert(transaction, models.Message{
			Status:     StatusQueued,
			CampaignID: campaignID,
		})
		if err != nil {
			transaction.Rollback()
			return
		}

		options.Endorsement = user.Endorsement

		job := gobble.NewJob(Delivery{
			JobType:         "v2",
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
			CampaignID:      campaignID,
		})

		_, err = enqueuer.queue.Enqueue(job, transaction)
		if err != nil {
			transaction.Rollback()
			return
		}

		recipient := user.Email
		if recipient == "" {
			recipient = user.GUID
		}
	}

	transaction.Commit()
}
