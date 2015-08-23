package mocks

import "github.com/cloudfoundry-incubator/notifications/v2/models"

type SendersRepository struct {
	InsertCall struct {
		Receives struct {
			Conn   models.ConnectionInterface
			Sender models.Sender
		}
		Returns struct {
			Sender models.Sender
			Err    error
		}
	}

	ListCall struct {
		Receives struct {
			Conn     models.ConnectionInterface
			ClientID string
		}
		Returns struct {
			Senders []models.Sender
			Err     error
		}
	}

	GetCall struct {
		Receives struct {
			Conn     models.ConnectionInterface
			SenderID string
		}
		Returns struct {
			Sender models.Sender
			Err    error
		}
	}

	GetByClientIDAndNameCall struct {
		Receives struct {
			Conn     models.ConnectionInterface
			ClientID string
			Name     string
		}
		Returns struct {
			Sender models.Sender
			Err    error
		}
	}
}

func NewSendersRepository() *SendersRepository {
	return &SendersRepository{}
}

func (s *SendersRepository) Insert(conn models.ConnectionInterface, sender models.Sender) (models.Sender, error) {
	s.InsertCall.Receives.Conn = conn
	s.InsertCall.Receives.Sender = sender

	return s.InsertCall.Returns.Sender, s.InsertCall.Returns.Err
}

func (s *SendersRepository) List(conn models.ConnectionInterface, clientID string) ([]models.Sender, error) {
	s.ListCall.Receives.Conn = conn
	s.ListCall.Receives.ClientID = clientID

	return s.ListCall.Returns.Senders, s.ListCall.Returns.Err
}

func (s *SendersRepository) Get(conn models.ConnectionInterface, senderID string) (models.Sender, error) {
	s.GetCall.Receives.Conn = conn
	s.GetCall.Receives.SenderID = senderID

	return s.GetCall.Returns.Sender, s.GetCall.Returns.Err
}

func (s *SendersRepository) GetByClientIDAndName(conn models.ConnectionInterface, clientID, name string) (models.Sender, error) {
	s.GetByClientIDAndNameCall.Receives.Conn = conn
	s.GetByClientIDAndNameCall.Receives.ClientID = clientID
	s.GetByClientIDAndNameCall.Receives.Name = name

	return s.GetByClientIDAndNameCall.Returns.Sender, s.GetByClientIDAndNameCall.Returns.Err
}
