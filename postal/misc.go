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

type TypedGUID interface {
    IsTypeEmail() bool
    String() string
}

type Templates struct {
    Subject string
    Text    string
    HTML    string
}

type EmailID string

func (guid EmailID) IsTypeEmail() bool {
    return true
}

func (guid EmailID) String() string {
    return string(guid)
}

type UAAGUID string

func (guid UAAGUID) IsTypeEmail() bool {
    return false
}

func (guid UAAGUID) String() string {
    return string(guid)
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
