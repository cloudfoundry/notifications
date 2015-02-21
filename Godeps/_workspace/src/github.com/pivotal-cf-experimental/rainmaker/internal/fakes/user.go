package fakes

import (
	"encoding/json"
	"time"
)

type User struct {
	GUID             string
	Admin            bool
	Active           bool
	DefaultSpaceGUID string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewUser(guid string) User {
	return User{
		GUID: guid,
	}
}

func (user User) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"metadata": map[string]interface{}{
			"guid":       user.GUID,
			"url":        "/v2/users/" + user.GUID,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
		"entity": map[string]interface{}{
			"admin":                             user.Admin,
			"active":                            user.Active,
			"default_space_guid":                user.DefaultSpaceGUID,
			"spaces_url":                        "/v2/users/" + user.GUID + "/spaces",
			"organizations_url":                 "/v2/users/" + user.GUID + "/organizations",
			"managed_organizations_url":         "/v2/users/" + user.GUID + "/managed_organizations",
			"billing_managed_organizations_url": "/v2/users/" + user.GUID + "/billing_managed_organizations",
			"audited_organizations_url":         "/v2/users/" + user.GUID + "/audited_organizations",
			"managed_spaces_url":                "/v2/users/" + user.GUID + "/managed_spaces",
			"audited_spaces_url":                "/v2/users/" + user.GUID + "/audited_spaces",
		},
	})
}
