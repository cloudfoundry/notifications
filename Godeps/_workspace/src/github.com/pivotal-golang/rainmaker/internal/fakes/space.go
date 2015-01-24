package fakes

import (
	"encoding/json"
	"time"
)

type Space struct {
	GUID             string
	Name             string
	Users            *Users
	OrganizationGUID string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Developers       *Users
}

func NewSpace(guid string) Space {
	return Space{
		GUID:       guid,
		Users:      NewUsers(),
		Developers: NewUsers(),
	}
}

func (space Space) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"metadata": map[string]interface{}{
			"guid":       space.GUID,
			"url":        "/v2/spaces/" + space.GUID,
			"created_at": space.CreatedAt,
			"updated_at": space.UpdatedAt,
		},
		"entity": map[string]interface{}{
			"name":                        space.Name,
			"organization_guid":           space.OrganizationGUID,
			"space_quota_definition_guid": nil,
			"organization_url":            "/v2/organizations/" + space.OrganizationGUID,
			"developers_url":              "/v2/spaces/" + space.GUID + "/developers",
			"managers_url":                "/v2/spaces/" + space.GUID + "/managers",
			"auditors_url":                "/v2/spaces/" + space.GUID + "/auditors",
			"apps_url":                    "/v2/spaces/" + space.GUID + "/apps",
			"routes_url":                  "/v2/spaces/" + space.GUID + "/routes",
			"domains_url":                 "/v2/spaces/" + space.GUID + "/domains",
			"service_instances_url":       "/v2/spaces/" + space.GUID + "/service_instances",
			"app_events_url":              "/v2/spaces/" + space.GUID + "/app_events",
			"events_url":                  "/v2/spaces/" + space.GUID + "/events",
			"security_groups_url":         "/v2/spaces/" + space.GUID + "/security_groups",
		},
	})
}
