package queue

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
)

const StatusQueued = "queued"

type GUIDGenerationFunc func() (string, error)

type User struct {
	GUID  string
	Email string
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

type JobEnqueuer struct {
	queue         Enqueuer
	guidGenerator GUIDGenerationFunc
	messagesRepo  messagesRepoInserter
}

func NewJobEnqueuer(queue Enqueuer, guidGenerator GUIDGenerationFunc, messagesRepo messagesRepoInserter) JobEnqueuer {
	return JobEnqueuer{
		queue:         queue,
		guidGenerator: guidGenerator,
		messagesRepo:  messagesRepo,
	}
}

func (enqueuer JobEnqueuer) Enqueue(conn ConnectionInterface, users []User, options Options, space cf.CloudControllerSpace, organization cf.CloudControllerOrganization, clientID, uaaHost, scope, vcapRequestID string, reqReceived time.Time, campaignID string) []Response {
	responses := []Response{}
	jobsByMessageID := map[string]gobble.Job{}
	for _, user := range users {
		messageID, err := enqueuer.guidGenerator()
		if err != nil {
			panic(err)
		}

		jobsByMessageID[messageID] = gobble.NewJob(Delivery{
			JobType:         "v2",
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
			CampaignID:      campaignID,
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
		_, err := enqueuer.messagesRepo.Insert(transaction, models.Message{
			ID:         messageID,
			Status:     StatusQueued,
			CampaignID: campaignID,
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
