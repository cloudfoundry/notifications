package fakes

import (
	"encoding/json"
	"time"
)

type ServiceInstance struct {
	GUID      string
	Name      string
	PlanGUID  string
	SpaceGUID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewServiceInstance(guid string) ServiceInstance {
	return ServiceInstance{
		GUID: guid,
	}
}

func (instance ServiceInstance) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"metadata": map[string]interface{}{
			"guid":       instance.GUID,
			"url":        "/v2/service_instances/" + instance.GUID,
			"created_at": instance.CreatedAt,
			"updated_at": instance.UpdatedAt,
		},
		"entity": map[string]interface{}{
			"name":              instance.Name,
			"credentials":       map[string]interface{}{},
			"service_plan_guid": instance.PlanGUID,
			"space_guid":        instance.SpaceGUID,
			"gateway_data":      "CONFIGURATION",
			"dashboard_url":     "",
			"type":              "managed_service_instance",
			"last_operation": map[string]interface{}{
				"type":        "create",
				"state":       "succeeded",
				"description": "",
				"updated_at":  instance.UpdatedAt,
			},
			"space_url":            "/v2/spaces/" + instance.SpaceGUID,
			"service_plan_url":     "/v2/service_plans/" + instance.PlanGUID,
			"service_bindings_url": "/v2/service_instances/" + instance.GUID + "/service_bindings",
		},
	})
}
