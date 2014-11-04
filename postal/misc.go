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
    Email          string `json:"email"`
}

type TypedGUID interface {
    BelongsToSpace() bool
    BelongsToOrganization() bool
    IsTypeEmail() bool
    String() string
}

type EmailID string

func NewEmailID() EmailID {
    return EmailID("")
}

func (guid EmailID) BelongsToSpace() bool {
    return false
}

func (guid EmailID) BelongsToOrganization() bool {
    return false
}

func (guid EmailID) IsTypeEmail() bool {
    return true
}

func (guid EmailID) String() string {
    return string(guid)
}

type SpaceGUID string

func (guid SpaceGUID) BelongsToSpace() bool {
    return true
}

func (guid SpaceGUID) BelongsToOrganization() bool {
    return false
}

func (guid SpaceGUID) IsTypeEmail() bool {
    return false
}

func (guid SpaceGUID) String() string {
    return string(guid)
}

type UserGUID string

func NewUserGUID() UserGUID {
    return UserGUID("")
}

func (guid UserGUID) BelongsToSpace() bool {
    return false
}

func (guid UserGUID) BelongsToOrganization() bool {
    return false
}

func (guid UserGUID) IsTypeEmail() bool {
    return false
}

func (guid UserGUID) String() string {
    return string(guid)
}

type OrganizationGUID string

func NewOrganizationGUID() OrganizationGUID {
    return OrganizationGUID("")
}

func (guid OrganizationGUID) BelongsToSpace() bool {
    return false
}

func (guid OrganizationGUID) BelongsToOrganization() bool {
    return true
}

func (guid OrganizationGUID) IsTypeEmail() bool {
    return false
}

func (guid OrganizationGUID) String() string {
    return string(guid)
}

type GUIDGenerationFunc func() (*uuid.UUID, error)

type UAAInterface interface {
    uaa.GetClientTokenInterface
    uaa.SetTokenInterface
    uaa.UsersEmailsByIDsInterface
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
}
