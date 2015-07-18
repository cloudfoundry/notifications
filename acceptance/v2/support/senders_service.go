package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Sender struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SendersService struct {
	config Config
}

func NewSendersService(config Config) SendersService {
	return SendersService{
		config: config,
	}
}

func (s SendersService) Create(name, token string) (Sender, error) {
	var sender Sender

	content, err := json.Marshal(map[string]string{
		"name": name,
	})
	if err != nil {
		return sender, err
	}

	status, body, err := NewClient(s.config).makeRequest("POST", s.config.Host+"/senders", bytes.NewBuffer(content), token)
	if err != nil {
		return sender, err
	}

	if status != http.StatusCreated {
		return sender, UnexpectedStatusError{status, string(body)}
	}

	err = json.Unmarshal(body, &sender)
	if err != nil {
		return sender, err
	}

	return sender, nil
}

func (s SendersService) Get(id, token string) (Sender, error) {
	var sender Sender

	status, body, err := NewClient(s.config).makeRequest("GET", s.config.Host+"/senders/"+id, nil, token)
	if err != nil {
		return sender, err
	}

	if status != http.StatusOK {
		return sender, UnexpectedStatusError{status, string(body)}
	}

	err = json.Unmarshal(body, &sender)
	if err != nil {
		return sender, err
	}

	return sender, nil
}

type UnexpectedStatusError struct {
	Status int
	Body   string
}

func (e UnexpectedStatusError) Error() string {
	return fmt.Sprintf("Unexpected status %d: %s", e.Status, e.Body)
}
