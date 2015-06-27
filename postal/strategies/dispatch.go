package strategies

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type Dispatch struct {
	GUID        string
	VCAPRequest VCAPRequest
	Message     Message
	Connection  models.ConnectionInterface
	Kind        Kind
	Client      Client
	Role        string
}

type VCAPRequest struct {
	ID          string
	ReceiptTime time.Time
}

type Message struct {
	To      string
	ReplyTo string
	Subject string
	Text    string
	HTML    HTML
}

type HTML struct {
	BodyContent    string
	BodyAttributes string
	Head           string
	Doctype        string
}

type Client struct {
	ID          string
	Description string
}

type Kind struct {
	ID          string
	Description string
}
