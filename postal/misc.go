package postal

import (
    "github.com/nu7hatch/gouuid"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const (
    StatusUnavailable = "unavailable"
    StatusFailed      = "failed"
    StatusDelivered   = "delivered"
    StatusNotFound    = "notfound"
    StatusNoAddress   = "noaddress"
    StatusQueued      = "queued"
)

type Response struct {
    Status         string `json:"status"`
    Recipient      string `json:"recipient"`
    NotificationID string `json:"notification_id"`
}

type TypedGUID interface {
    BelongsToSpace() bool
    String() string
}

type SpaceGUID string

func (guid SpaceGUID) BelongsToSpace() bool {
    return true
}

func (guid SpaceGUID) String() string {
    return string(guid)
}

type UserGUID string

func (guid UserGUID) BelongsToSpace() bool {
    return false
}

func (guid UserGUID) String() string {
    return string(guid)
}

type GUIDGenerationFunc func() (*uuid.UUID, error)

type UAAInterface interface {
    uaa.GetClientTokenInterface
    uaa.SetTokenInterface
    uaa.UsersByIDsInterface
}

type Options struct {
    ReplyTo           string
    Subject           string
    KindDescription   string
    SourceDescription string
    Text              string
    HTML              string
    KindID            string
}
