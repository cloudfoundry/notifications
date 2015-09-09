package mocks

import "github.com/cloudfoundry-incubator/notifications/v1/models"

type MessagesRepository struct {
	UpsertCall  messagesRepositoryUpsertCall
	UpsertCalls []messagesRepositoryUpsertCall
}

type messagesRepositoryUpsertCall struct {
	Receives struct {
		Connection models.ConnectionInterface
		Message    models.Message
	}
	Returns struct {
		Message models.Message
		Error   error
	}
}

func NewMessagesRepository() *MessagesRepository {
	return &MessagesRepository{}
}

func (mr *MessagesRepository) Upsert(conn models.ConnectionInterface, message models.Message) (models.Message, error) {
	mr.UpsertCall.Receives.Connection = conn
	mr.UpsertCall.Receives.Message = message

	mr.UpsertCalls = append(mr.UpsertCalls, mr.UpsertCall)

	return mr.UpsertCall.Returns.Message, mr.UpsertCall.Returns.Error
}
