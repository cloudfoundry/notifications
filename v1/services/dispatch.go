package services

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type Dispatch struct {
	GUID       string
	Role       string
	Connection models.ConnectionInterface
	UAAHost    string

	VCAPRequest DispatchVCAPRequest
	Message     DispatchMessage
	Kind        DispatchKind
	Client      DispatchClient
}

type HTML struct {
	BodyContent    string
	BodyAttributes string
	Head           string
	Doctype        string
}

type DispatchVCAPRequest struct {
	ID          string
	ReceiptTime time.Time
}

type DispatchMessage struct {
	To      string
	ReplyTo string
	Subject string
	Text    string
	HTML    HTML
}

type DispatchClient struct {
	ID          string
	Description string
}

type DispatchKind struct {
	ID          string
	Description string
}
