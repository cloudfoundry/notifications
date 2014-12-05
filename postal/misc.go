package postal

import "github.com/nu7hatch/gouuid"

const (
	StatusUnavailable = "unavailable"
	StatusFailed      = "failed"
	StatusDelivered   = "delivered"
	StatusNotFound    = "notfound"
	StatusNoAddress   = "noaddress"
	StatusQueued      = "queued"
)

type Templates struct {
	Name    string
	Subject string
	Text    string
	HTML    string
}

type GUIDGenerationFunc func() (*uuid.UUID, error)

type HTML struct {
	BodyContent    string
	BodyAttributes string
	Head           string
	Doctype        string
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
}
