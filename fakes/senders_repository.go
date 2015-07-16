package fakes

import "github.com/cloudfoundry-incubator/notifications/models"

type SendersRepository struct {
	InsertCall struct {
		Conn         models.ConnectionInterface
		Sender       models.Sender
		ReturnSender models.Sender
		Err          error
	}

	GetByClientIDAndNameCall struct {
		ReturnSender models.Sender
		Err          error
		Conn         models.ConnectionInterface
		ClientID     string
		Name         string
	}
}

func NewSendersRepository() *SendersRepository {
	return &SendersRepository{}
}

func (s *SendersRepository) Insert(conn models.ConnectionInterface, sender models.Sender) (models.Sender, error) {
	s.InsertCall.Conn = conn
	s.InsertCall.Sender = sender
	return s.InsertCall.ReturnSender, s.InsertCall.Err
}

func (s *SendersRepository) GetByClientIDAndName(conn models.ConnectionInterface, clientID, name string) (models.Sender, error) {
	s.GetByClientIDAndNameCall.Conn = conn
	s.GetByClientIDAndNameCall.ClientID = clientID
	s.GetByClientIDAndNameCall.Name = name
	return s.GetByClientIDAndNameCall.ReturnSender, s.GetByClientIDAndNameCall.Err
}
