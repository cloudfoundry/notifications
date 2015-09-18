package mocks

import "github.com/cloudfoundry-incubator/notifications/v2/models"

type MessagesRepository struct {
	InsertCall       messagesRepositoryInsertCall
	InsertCalls      messagesRepositoryInsertCalls
	InsertCallsCount int

	CountByStatusCall struct {
		Receives struct {
			CampaignIDList []string
			Connection     models.ConnectionInterface
		}

		Returns struct {
			MessageCounts models.MessageCounts
			Error         error
		}
	}

	MostRecentlyUpdatedByCampaignIDCall struct {
		Receives struct {
			CampaignID string
			Connection models.ConnectionInterface
		}

		Returns struct {
			Message models.Message
			Error   error
		}
	}

	UpdateCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			Message    models.Message
		}
		Returns struct {
			Message models.Message
			Error   error
		}
	}
}

type messagesRepositoryInsertCall struct {
	Receives struct {
		Connection models.ConnectionInterface
		Message    models.Message
	}
	Returns struct {
		Message models.Message
		Error   error
	}
}

type messagesRepositoryInsertCalls []messagesRepositoryInsertCall

func (c messagesRepositoryInsertCalls) WithMessages(messages []models.Message) messagesRepositoryInsertCalls {
	var calls messagesRepositoryInsertCalls
	for _, message := range messages {
		call := messagesRepositoryInsertCall{}
		call.Returns.Message = message
		calls = append(calls, call)
	}

	return calls
}

func NewMessagesRepository() *MessagesRepository {
	return &MessagesRepository{}
}

func (mr *MessagesRepository) CountByStatus(conn models.ConnectionInterface, campaignID string) (models.MessageCounts, error) {
	mr.CountByStatusCall.Receives.Connection = conn
	mr.CountByStatusCall.Receives.CampaignIDList = append(mr.CountByStatusCall.Receives.CampaignIDList, campaignID)

	return mr.CountByStatusCall.Returns.MessageCounts, mr.CountByStatusCall.Returns.Error
}

func (mr *MessagesRepository) MostRecentlyUpdatedByCampaignID(conn models.ConnectionInterface, campaignID string) (models.Message, error) {
	mr.MostRecentlyUpdatedByCampaignIDCall.Receives.Connection = conn
	mr.MostRecentlyUpdatedByCampaignIDCall.Receives.CampaignID = campaignID

	return mr.MostRecentlyUpdatedByCampaignIDCall.Returns.Message, mr.MostRecentlyUpdatedByCampaignIDCall.Returns.Error
}

func (mr *MessagesRepository) Insert(conn models.ConnectionInterface, message models.Message) (models.Message, error) {
	if len(mr.InsertCalls) <= mr.InsertCallsCount {
		mr.InsertCalls = append(mr.InsertCalls, messagesRepositoryInsertCall{})
	}
	mr.InsertCall = mr.InsertCalls[mr.InsertCallsCount]

	mr.InsertCall.Receives.Connection = conn
	mr.InsertCall.Receives.Message = message
	mr.InsertCalls[mr.InsertCallsCount] = mr.InsertCall

	mr.InsertCallsCount++

	return mr.InsertCall.Returns.Message, mr.InsertCall.Returns.Error
}

func (mr *MessagesRepository) Update(conn models.ConnectionInterface, message models.Message) (models.Message, error) {
	mr.UpdateCall.Receives.Connection = conn
	mr.UpdateCall.Receives.Message = message

	return mr.UpdateCall.Returns.Message, mr.UpdateCall.Returns.Error
}
