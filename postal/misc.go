package postal

const (
	StatusUnavailable   = "unavailable"
	StatusFailed        = "failed"
	StatusRetry         = "retry"
	StatusDelivered     = "delivered"
	StatusQueued        = "queued"
	StatusUndeliverable = "undeliverable"
)

type Templates struct {
	Name    string
	Subject string
	Text    string
	HTML    string
}

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
	Endorsement       string
	TemplateID        string
}
