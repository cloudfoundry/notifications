package servers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type CC struct {
	server          *httptest.Server
	userNameToIdMap map[string]string
}

func NewCC(userNameToIdMap map[string]string) CC {
	router := mux.NewRouter()
	cc := CC{
		server:          httptest.NewUnstartedServer(router),
		userNameToIdMap: userNameToIdMap,
	}

	router.HandleFunc("/v2/spaces/{guid}", cc.GetSpace).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}/users", cc.GetOrgUsers).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}/managers", cc.GetOrgManagers).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}/auditors", cc.GetOrgAuditors).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}/billing_managers", cc.GetOrgBillingManagers).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}", cc.GetOrg).Methods("GET")
	router.HandleFunc("/v2/users", cc.GetSpaceUsers).Methods("GET")
	router.HandleFunc("/{anything:.*}", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("CC ROUTE REQUEST ---> %+v\n", req)
		w.WriteHeader(http.StatusTeapot)
	}))

	return cc
}

func (s CC) Boot() {
	s.server.Start()
	os.Setenv("CC_HOST", s.server.URL)
}

func (s CC) Close() {
	s.server.Close()
}

func (cc CC) GetSpace(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	guid := vars["guid"]
	if guid == "space-123" || guid == "space-456" {
		output := map[string]interface{}{
			"metadata": map[string]interface{}{
				"guid":       guid,
				"url":        fmt.Sprintf("/v2/spaces/%s", guid),
				"created_at": "2014-08-01T17:36:18+00:00",
				"updated_at": nil,
			},
			"entity": map[string]interface{}{
				"name":                  "notifications-service",
				"organization_guid":     "org-123",
				"organization_url":      "/v2/organizations/org-123",
				"developers_url":        fmt.Sprintf("/v2/spaces/%s/developers", guid),
				"managers_url":          fmt.Sprintf("/v2/spaces/%s/managers", guid),
				"auditors_url":          fmt.Sprintf("/v2/spaces/%s/auditors", guid),
				"apps_url":              fmt.Sprintf("/v2/spaces/%s/apps", guid),
				"routes_url":            fmt.Sprintf("/v2/spaces/%s/routes", guid),
				"domains_url":           fmt.Sprintf("/v2/spaces/%s/domains", guid),
				"service_instances_url": fmt.Sprintf("/v2/spaces/%s/service_instances", guid),
				"app_events_url":        fmt.Sprintf("/v2/spaces/%s/app_events", guid),
				"events_url":            fmt.Sprintf("/v2/spaces/%s/events", guid),
				"security_groups_url":   fmt.Sprintf("/v2/spaces/%s/security_groups", guid),
			},
		}
		response, err := json.Marshal(output)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"code":40004,"description":"The app space could not be found","error_code":"CF-SpaceNotFound"}`))
	}
}

func (cc CC) GetOrg(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	guid := vars["guid"]
	if guid == "org-123" || guid == "org-456" {
		json, err := json.Marshal(map[string]interface{}{
			"metadata": map[string]string{
				"guid":       guid,
				"url":        fmt.Sprintf("/v2/organizations/%s", guid),
				"created_at": "2014-08-01T17:36:17+00:00",
				"updated_at": "2014-08-01T17:36:20+00:00",
			},
			"entity": map[string]interface{}{
				"name":                  "notifications-service",
				"billing_enabled":       false,
				"quota_definition_guid": "73530fc0-17a3-42f1-9692-838860d30ec2",
				"status":                "active",
				"quota_definition_url":  "/v2/quota_definitions/73530fc0-17a3-42f1-9692-838860d30ec2",
				"spaces_url":            fmt.Sprintf("/v2/organizations/%s/spaces", guid),
				"domains_url":           fmt.Sprintf("/v2/organizations/%s/domains", guid),
				"private_domains_url":   fmt.Sprintf("/v2/organizations/%s/private_domains", guid),
				"users_url":             fmt.Sprintf("/v2/organizations/%s/users", guid),
				"managers_url":          fmt.Sprintf("/v2/organizations/%s/managers", guid),
				"billing_managers_url":  fmt.Sprintf("/v2/organizations/%s/billing_managers", guid),
				"auditors_url":          fmt.Sprintf("/v2/organizations/%s/auditors", guid),
				"app_events_url":        fmt.Sprintf("/v2/organizations/%s/app_events", guid),
			},
		})
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(json)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"code":40004,"description":"The org could not be found","error_code":"CF-OrgNotFound"}`))
	}
}

func (cc CC) GetOrgUsers(w http.ResponseWriter, req *http.Request) {
	orgGUID := strings.Split(req.URL.Path, "/")[3]

	var desiredUsers []string
	switch orgGUID {
	case "org-123":
		desiredUsers = []string{
			"user-456",
			"user-789",
			"user-000",
		}
	case "org-456":
		desiredUsers = []string{
			"user-123",
			"user-456",
		}
	default:
		desiredUsers = []string{}
	}

	users := []map[string]interface{}{}
	for _, userName := range desiredUsers {
		guid, ok := cc.userNameToIdMap[userName]
		if !ok {
			guid = userName
		}

		users = append(users, map[string]interface{}{
			"metadata": map[string]interface{}{
				"guid":       guid,
				"url":        fmt.Sprintf("/v2/users/%s", guid),
				"created_at": "2014-07-16T21:58:29+00:00",
				"updated_at": nil,
			},
			"entity": map[string]interface{}{
				"admin":                             false,
				"active":                            true,
				"default_space_guid":                nil,
				"spaces_url":                        fmt.Sprintf("/v2/users/%s/spaces", guid),
				"organizations_url":                 fmt.Sprintf("/v2/users/%s/organizations", guid),
				"managed_organizations_url":         fmt.Sprintf("/v2/users/%s/managed_organizations", guid),
				"billing_managed_organizations_url": fmt.Sprintf("/v2/users/%s/billing_managed_organizations", guid),
				"audited_organizations_url":         fmt.Sprintf("/v2/users/%s/audited_organizations", guid),
				"managed_spaces_url":                fmt.Sprintf("/v2/users/%s/managed_spaces", guid),
				"audited_spaces_url":                fmt.Sprintf("/v2/users/%s/audited_spaces", guid),
			},
		})
	}

	output := map[string]interface{}{
		"total_results": len(users),
		"total_pages":   1,
		"prev_url":      nil,
		"next_url":      nil,
		"resources":     users,
	}

	json, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(json))
}

func (cc CC) GetOrgManagers(w http.ResponseWriter, req *http.Request) {
	json := `{
       "total_results": 1,
       "total_pages": 1,
       "prev_url": null,
       "next_url": null,
       "resources": [
          {
             "metadata": {
                "guid": "user-456",
                "url": "/v2/users/user-456",
                "created_at": "2014-10-16T21:05:40+00:00",
                "updated_at": null
             },
             "entity": {
                "admin": false,
                "active": true,
                "default_space_guid": null,
                "spaces_url": "/v2/users/user-456/spaces",
                "organizations_url": "/v2/users/user-456/organizations",
                "managed_organizations_url": "/v2/users/user-456/managed_organizations",
                "billing_managed_organizations_url": "/v2/users/user-456/billing_managed_organizations",
                "audited_organizations_url": "/v2/users/user-456/audited_organizations",
                "managed_spaces_url": "/v2/users/user-456/managed_spaces",
                "audited_spaces_url": "/v2/users/user-456/audited_spaces"
             }
          }
       ]
    }`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(json))
}

func (cc CC) GetOrgAuditors(w http.ResponseWriter, req *http.Request) {
	json := `{
       "total_results": 1,
       "total_pages": 1,
       "prev_url": null,
       "next_url": null,
       "resources": [
          {
             "metadata": {
                "guid": "user-123",
                "url": "/v2/users/user-123",
                "created_at": "2014-10-16T21:05:40+00:00",
                "updated_at": null
             },
             "entity": {
                "admin": false,
                "active": true,
                "default_space_guid": null,
                "spaces_url": "/v2/users/user-123/spaces",
                "organizations_url": "/v2/users/user-123/organizations",
                "managed_organizations_url": "/v2/users/user-123/managed_organizations",
                "billing_managed_organizations_url": "/v2/users/user-123/billing_managed_organizations",
                "audited_organizations_url": "/v2/users/user-123/audited_organizations",
                "managed_spaces_url": "/v2/users/user-123/managed_spaces",
                "audited_spaces_url": "/v2/users/user-123/audited_spaces"
             }
          }
       ]
    }`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(json))
}

func (cc CC) GetOrgBillingManagers(w http.ResponseWriter, req *http.Request) {
	json := `{
       "total_results": 1,
       "total_pages": 1,
       "prev_url": null,
       "next_url": null,
       "resources": [
          {
             "metadata": {
                "guid": "user-111",
                "url": "/v2/users/user-111",
                "created_at": "2014-10-16T21:05:40+00:00",
                "updated_at": null
             },
             "entity": {
                "admin": false,
                "active": true,
                "default_space_guid": null,
                "spaces_url": "/v2/users/user-111/spaces",
                "organizations_url": "/v2/users/user-111/organizations",
                "managed_organizations_url": "/v2/users/user-111/managed_organizations",
                "billing_managed_organizations_url": "/v2/users/user-111/billing_managed_organizations",
                "audited_organizations_url": "/v2/users/user-111/audited_organizations",
                "managed_spaces_url": "/v2/users/user-111/managed_spaces",
                "audited_spaces_url": "/v2/users/user-111/audited_spaces"
             }
          }
       ]
    }`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(json))
}

func (cc CC) GetSpaceUsers(w http.ResponseWriter, req *http.Request) {
	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		panic(err)
	}
	filter := query.Get("q")
	spaceGUID := strings.TrimPrefix(filter, "space_guid:")

	var desiredUsers []string
	switch spaceGUID {
	case "space-123":
		desiredUsers = []string{
			"user-456",
			"user-789",
			"user-000",
		}
	case "space-456":
		desiredUsers = []string{
			"user-123",
			"user-456",
		}
	default:
		desiredUsers = []string{}
	}

	users := []map[string]interface{}{}
	for _, userName := range desiredUsers {
		guid, ok := cc.userNameToIdMap[userName]
		if !ok {
			guid = userName
		}

		users = append(users, map[string]interface{}{
			"metadata": map[string]interface{}{
				"guid":       guid,
				"url":        fmt.Sprintf("/v2/users/%s", guid),
				"created_at": "2014-07-16T21:58:29+00:00",
				"updated_at": nil,
			},
			"entity": map[string]interface{}{
				"admin":                             false,
				"active":                            true,
				"default_space_guid":                nil,
				"spaces_url":                        fmt.Sprintf("/v2/users/%s/spaces", guid),
				"organizations_url":                 fmt.Sprintf("/v2/users/%s/organizations", guid),
				"managed_organizations_url":         fmt.Sprintf("/v2/users/%s/managed_organizations", guid),
				"billing_managed_organizations_url": fmt.Sprintf("/v2/users/%s/billing_managed_organizations", guid),
				"audited_organizations_url":         fmt.Sprintf("/v2/users/%s/audited_organizations", guid),
				"managed_spaces_url":                fmt.Sprintf("/v2/users/%s/managed_spaces", guid),
				"audited_spaces_url":                fmt.Sprintf("/v2/users/%s/audited_spaces", guid),
			},
		})
	}

	output := map[string]interface{}{
		"total_results": len(users),
		"total_pages":   1,
		"prev_url":      nil,
		"next_url":      nil,
		"resources":     users,
	}

	uaaJSON, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(uaaJSON))
}
