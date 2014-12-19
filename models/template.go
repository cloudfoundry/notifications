package models

import "time"

const (
	UserBodyTemplateName         = "user_body"
	SpaceBodyTemplateName        = "space_body"
	EmailBodyTemplateName        = "email_body"
	UAAScopeBodyTemplateName     = "uaa_scope_body"
	OrganizationBodyTemplateName = "organization_body"
	EveryoneBodyTemplateName     = "everyone_body"
	SubjectMissingTemplateName   = "subject.missing"
	SubjectProvidedTemplateName  = "subject.provided"
)

var TemplateNames = []string{
	UserBodyTemplateName,
	SpaceBodyTemplateName,
	EmailBodyTemplateName,
	UAAScopeBodyTemplateName,
	OrganizationBodyTemplateName,
	EveryoneBodyTemplateName,
	SubjectMissingTemplateName,
	SubjectProvidedTemplateName,
}

type Template struct {
	Primary   int       `db:"primary"`
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Subject   string    `db:"subject"`
	Text      string    `db:"text"`
	HTML      string    `db:"html"`
	Metadata  string    `db:"metadata"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
