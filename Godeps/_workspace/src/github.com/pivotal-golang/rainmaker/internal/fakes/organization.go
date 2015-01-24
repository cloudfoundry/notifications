package fakes

import (
	"encoding/json"
	"time"
)

type Organization struct {
	GUID                string
	Name                string
	Status              string
	BillingEnabled      bool
	QuotaDefinitionGUID string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Users               *Users
	BillingManagers     *Users
	Auditors            *Users
	Managers            *Users
}

func NewOrganization(guid string) Organization {
	return Organization{
		GUID:            guid,
		Users:           NewUsers(),
		BillingManagers: NewUsers(),
		Auditors:        NewUsers(),
		Managers:        NewUsers(),
	}
}

func (org Organization) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"metadata": map[string]interface{}{
			"guid":       org.GUID,
			"url":        "/v2/organizations/" + org.GUID,
			"created_at": org.CreatedAt,
			"updated_at": org.UpdatedAt,
		},
		"entity": map[string]interface{}{
			"name":                        org.Name,
			"billing_enabled":             org.BillingEnabled,
			"quota_definition_guid":       org.QuotaDefinitionGUID,
			"status":                      org.Status,
			"quota_definition_url":        "/v2/quota_definitions/" + org.QuotaDefinitionGUID,
			"spaces_url":                  "/v2/organizations/" + org.GUID + "/spaces",
			"domains_url":                 "/v2/organizations/" + org.GUID + "/domains",
			"private_domains_url":         "/v2/organizations/" + org.GUID + "/private_domains",
			"users_url":                   "/v2/organizations/" + org.GUID + "/users",
			"managers_url":                "/v2/organizations/" + org.GUID + "/managers",
			"billing_managers_url":        "/v2/organizations/" + org.GUID + "/billing_managers",
			"auditors_url":                "/v2/organizations/" + org.GUID + "/auditors",
			"app_events_url":              "/v2/organizations/" + org.GUID + "/app_events",
			"space_quota_definitions_url": "/v2/organizations/" + org.GUID + "/space_quota_definitions",
		},
	})
}
