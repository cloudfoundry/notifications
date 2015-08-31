package mocks

import "github.com/cloudfoundry-incubator/notifications/v2/models"

type SendersRepository struct {
	InsertCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			Sender     models.Sender
		}
		Returns struct {
			Sender models.Sender
			Error  error
		}
	}

	UpdateCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			Sender     models.Sender
		}
		Returns struct {
			Sender models.Sender
			Error  error
		}
	}

	ListCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			ClientID   string
		}
		Returns struct {
			Senders []models.Sender
			Error   error
		}
	}

	GetCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			SenderID   string
		}
		Returns struct {
			Sender models.Sender
			Error  error
		}
	}

	GetByClientIDAndNameCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			ClientID   string
			Name       string
		}
		Returns struct {
			Sender models.Sender
			Error  error
		}
	}

	DeleteCall struct {
		Receives struct {
			Connection models.ConnectionInterface
			Sender     models.Sender
		}
		Returns struct {
			Error error
		}
	}
}

func NewSendersRepository() *SendersRepository {
	return &SendersRepository{}
}

func (s *SendersRepository) Insert(conn models.ConnectionInterface, sender models.Sender) (models.Sender, error) {
	s.InsertCall.Receives.Connection = conn
	s.InsertCall.Receives.Sender = sender

	return s.InsertCall.Returns.Sender, s.InsertCall.Returns.Error
}

func (s *SendersRepository) Update(conn models.ConnectionInterface, sender models.Sender) (models.Sender, error) {
	s.UpdateCall.Receives.Connection = conn
	s.UpdateCall.Receives.Sender = sender

	return s.UpdateCall.Returns.Sender, s.UpdateCall.Returns.Error
}

func (s *SendersRepository) List(conn models.ConnectionInterface, clientID string) ([]models.Sender, error) {
	s.ListCall.Receives.Connection = conn
	s.ListCall.Receives.ClientID = clientID

	return s.ListCall.Returns.Senders, s.ListCall.Returns.Error
}

func (s *SendersRepository) Get(conn models.ConnectionInterface, senderID string) (models.Sender, error) {
	s.GetCall.Receives.Connection = conn
	s.GetCall.Receives.SenderID = senderID

	return s.GetCall.Returns.Sender, s.GetCall.Returns.Error
}

func (s *SendersRepository) GetByClientIDAndName(conn models.ConnectionInterface, clientID, name string) (models.Sender, error) {
	s.GetByClientIDAndNameCall.Receives.Connection = conn
	s.GetByClientIDAndNameCall.Receives.ClientID = clientID
	s.GetByClientIDAndNameCall.Receives.Name = name

	return s.GetByClientIDAndNameCall.Returns.Sender, s.GetByClientIDAndNameCall.Returns.Error
}

func (s *SendersRepository) Delete(conn models.ConnectionInterface, sender models.Sender) error {
	s.DeleteCall.Receives.Connection = conn
	s.DeleteCall.Receives.Sender = sender

	return s.DeleteCall.Returns.Error
}
