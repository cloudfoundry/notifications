package mocks

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/v1/models"
)

type MessagesRepo struct {
	UpsertCall struct {
		CallCount int
		Receives  struct {
			Connection models.ConnectionInterface
			Messages   []models.Message
		}
		Returns struct {
			Messages []models.Message
			Error    error
		}
	}

	UpdateCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			Messages   []models.Message
		}
		Returns struct {
			Message models.Message
			Error   error
		}
	}

	FindByIDCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			MessageID  string
		}
		Returns struct {
			Message models.Message
			Error   error
		}
	}

	DeleteBeforeCall struct {
		InvocationTimes []time.Time
		CallCount       int
		Receives        struct {
			Connection    models.ConnectionInterface
			ThresholdTime time.Time
		}
		Returns struct {
			RowsAffected int
			Error        error
		}
	}
}

func NewMessagesRepo() *MessagesRepo {
	return &MessagesRepo{}
}

func (mr *MessagesRepo) Upsert(conn models.ConnectionInterface, message models.Message) (models.Message, error) {
	mr.UpsertCall.Receives.Connection = conn
	mr.UpsertCall.Receives.Messages = append(mr.UpsertCall.Receives.Messages, message)

	message = mr.UpsertCall.Returns.Messages[mr.UpsertCall.CallCount]
	mr.UpsertCall.CallCount++

	return message, mr.UpsertCall.Returns.Error
}

func (mr *MessagesRepo) Update(conn models.ConnectionInterface, message models.Message) (models.Message, error) {
	mr.UpdateCall.Receives.Connection = conn
	mr.UpdateCall.Receives.Messages = append(mr.UpdateCall.Receives.Messages, message)

	return mr.UpdateCall.Returns.Message, mr.UpdateCall.Returns.Error
}

func (mr *MessagesRepo) FindByID(conn models.ConnectionInterface, messageID string) (models.Message, error) {
	mr.FindByIDCall.Receives.Connection = conn
	mr.FindByIDCall.Receives.MessageID = messageID

	return mr.FindByIDCall.Returns.Message, mr.FindByIDCall.Returns.Error
}

func (mr *MessagesRepo) DeleteBefore(conn models.ConnectionInterface, thresholdTime time.Time) (int, error) {
	mr.DeleteBeforeCall.Receives.Connection = conn
	mr.DeleteBeforeCall.Receives.ThresholdTime = thresholdTime
	mr.DeleteBeforeCall.InvocationTimes = append(mr.DeleteBeforeCall.InvocationTimes, time.Now())
	mr.DeleteBeforeCall.CallCount++

	return mr.DeleteBeforeCall.Returns.RowsAffected, mr.DeleteBeforeCall.Returns.Error
}
