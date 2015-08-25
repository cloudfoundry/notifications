package servers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type CC struct {
	server *httptest.Server
}

func NewCC() CC {
	router := mux.NewRouter()
	router.HandleFunc("/v2/spaces/{guid}", CCGetSpace).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}/users", CCGetOrgUsers).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}/managers", CCGetOrgManagers).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}/auditors", CCGetOrgAuditors).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}/billing_managers", CCGetOrgBillingManagers).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}", CCGetOrg).Methods("GET")
	router.HandleFunc("/v2/users", CCGetSpaceUsers).Methods("GET")
	router.HandleFunc("/{anything:.*}", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("CC ROUTE REQUEST ---> %+v\n", req)
		w.WriteHeader(http.StatusTeapot)
	}))

	return CC{
		server: httptest.NewUnstartedServer(router),
	}
}

func (s CC) Boot() {
	s.server.Start()
	os.Setenv("CC_HOST", s.server.URL)
}

func (s CC) Close() {
	s.server.Close()
}

var CCGetSpace = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	guid := vars["guid"]
	if guid == "space-123" {
		json := `{
		   "metadata": {
			  "guid": "space-123",
			  "url": "/v2/spaces/space-123",
			  "created_at": "2014-08-01T17:36:18+00:00",
			  "updated_at": null
		   },
		   "entity": {
			  "name": "notifications-service",
			  "organization_guid": "org-123",
			  "organization_url": "/v2/organizations/org-123",
			  "developers_url": "/v2/spaces/space-123/developers",
			  "managers_url": "/v2/spaces/space-123/managers",
			  "auditors_url": "/v2/spaces/space-123/auditors",
			  "apps_url": "/v2/spaces/space-123/apps",
			  "routes_url": "/v2/spaces/space-123/routes",
			  "domains_url": "/v2/spaces/space-123/domains",
			  "service_instances_url": "/v2/spaces/space-123/service_instances",
			  "app_events_url": "/v2/spaces/space-123/app_events",
			  "events_url": "/v2/spaces/space-123/events",
			  "security_groups_url": "/v2/spaces/space-123/security_groups"
		   }
		}`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"code":40004,"description":"The app space could not be found","error_code":"CF-SpaceNotFound"}`))
	}
})

var CCGetOrg = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	json := `{
       "metadata": {
          "guid": "org-123",
          "url": "/v2/organizations/org-123",
          "created_at": "2014-08-01T17:36:17+00:00",
          "updated_at": "2014-08-01T17:36:20+00:00"
       },
       "entity": {
          "name": "notifications-service",
          "billing_enabled": false,
          "quota_definition_guid": "73530fc0-17a3-42f1-9692-838860d30ec2",
          "status": "active",
          "quota_definition_url": "/v2/quota_definitions/73530fc0-17a3-42f1-9692-838860d30ec2",
          "spaces_url": "/v2/organizations/org-123/spaces",
          "domains_url": "/v2/organizations/org-123/domains",
          "private_domains_url": "/v2/organizations/org-123/private_domains",
          "users_url": "/v2/organizations/org-123/users",
          "managers_url": "/v2/organizations/org-123/managers",
          "billing_managers_url": "/v2/organizations/org-123/billing_managers",
          "auditors_url": "/v2/organizations/org-123/auditors",
          "app_events_url": "/v2/organizations/org-123/app_events"
       }
    }`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(json))
})

var CCGetOrgUsers = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	token := strings.Split(req.Header.Get("Authorization"), " ")[1]
	jwtToken, _ := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(helpers.UAAPublicKey), nil
	})

	var json string
	if regexp.MustCompile(string(os.Getenv("UAA_HOST"))).MatchString(jwtToken.Claims["iss"].(string)) {
		json = `{
		   "total_results": 2,
		   "total_pages": 1,
		   "prev_url": null,
		   "next_url": null,
		   "resources": [
			  {
				 "metadata": {
					"guid": "user-456",
					"url": "/v2/users/user-456",
					"created_at": "2014-07-16T21:58:29+00:00",
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
			  },
			  {
				 "metadata": {
					"guid": "user-789",
					"url": "/v2/users/user-789",
					"created_at": "2014-07-18T22:08:18+00:00",
					"updated_at": null
				 },
				 "entity": {
					"admin": false,
					"active": false,
					"default_space_guid": null,
					"spaces_url": "/v2/users/user-789/spaces",
					"organizations_url": "/v2/users/user-789/organizations",
					"managed_organizations_url": "/v2/users/user-789/managed_organizations",
					"billing_managed_organizations_url": "/v2/users/user-789/billing_managed_organizations",
					"audited_organizations_url": "/v2/users/user-789/audited_organizations",
					"managed_spaces_url": "/v2/users/user-789/managed_spaces",
					"audited_spaces_url": "/v2/users/user-789/audited_spaces"
				 }
			  },
			  {
				 "metadata": {
					"guid": "user-000",
					"url": "/v2/users/user-000",
					"created_at": "2014-07-16T21:58:29+00:00",
					"updated_at": null
				 },
				 "entity": {
					"admin": false,
					"active": true,
					"default_space_guid": null,
					"spaces_url": "/v2/users/user-000/spaces",
					"organizations_url": "/v2/users/user-000/organizations",
					"managed_organizations_url": "/v2/users/user-000/managed_organizations",
					"billing_managed_organizations_url": "/v2/users/user-000/billing_managed_organizations",
					"audited_organizations_url": "/v2/users/user-000/audited_organizations",
					"managed_spaces_url": "/v2/users/user-000/managed_spaces",
					"audited_spaces_url": "/v2/users/user-000/audited_spaces"
				 }
			  }
		   ]
		}`
	} else {
		json = `{
		   "total_results": 1,
		   "total_pages": 1,
		   "prev_url": null,
		   "next_url": null,
		   "resources": [
			  {
				 "metadata": {
					"guid": "another-user-in-zone",
					"url": "/v2/users/another-user-in-zone",
					"created_at": "2014-07-16T21:58:29+00:00",
					"updated_at": null
				 },
				 "entity": {
					"admin": false,
					"active": true,
					"default_space_guid": null,
					"spaces_url": "/v2/users/another-user-in-zone/spaces",
					"organizations_url": "/v2/users/another-user-in-zone/organizations",
					"managed_organizations_url": "/v2/users/another-user-in-zone/managed_organizations",
					"billing_managed_organizations_url": "/v2/users/another-user-in-zone/billing_managed_organizations",
					"audited_organizations_url": "/v2/users/another-user-in-zone/audited_organizations",
					"managed_spaces_url": "/v2/users/another-user-in-zone/managed_spaces",
					"audited_spaces_url": "/v2/users/another-user-in-zone/audited_spaces"
				 }
			  }
			]
		}`
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(json))
})

var CCGetOrgManagers = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
})

var CCGetOrgAuditors = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
})

var CCGetOrgBillingManagers = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
})

var CCGetSpaceUsers = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	uaaJSON := `{
       "total_results": 3,
       "total_pages": 1,
       "prev_url": null,
       "next_url": null,
       "resources": [
          {
             "metadata": {
                "guid": "user-456",
                "url": "/v2/users/user-456",
                "created_at": "2014-07-16T21:58:29+00:00",
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
          },
          {
             "metadata": {
                "guid": "user-789",
                "url": "/v2/users/user-789",
                "created_at": "2014-07-18T22:08:18+00:00",
                "updated_at": null
             },
             "entity": {
                "admin": false,
                "active": false,
                "default_space_guid": null,
                "spaces_url": "/v2/users/user-789/spaces",
                "organizations_url": "/v2/users/user-789/organizations",
                "managed_organizations_url": "/v2/users/user-789/managed_organizations",
                "billing_managed_organizations_url": "/v2/users/user-789/billing_managed_organizations",
                "audited_organizations_url": "/v2/users/user-789/audited_organizations",
                "managed_spaces_url": "/v2/users/user-789/managed_spaces",
                "audited_spaces_url": "/v2/users/user-789/audited_spaces"
             }
          },
          {
             "metadata": {
                "guid": "user-000",
                "url": "/v2/users/user-000",
                "created_at": "2014-07-16T21:58:29+00:00",
                "updated_at": null
             },
             "entity": {
                "admin": false,
                "active": true,
                "default_space_guid": null,
                "spaces_url": "/v2/users/user-000/spaces",
                "organizations_url": "/v2/users/user-000/organizations",
                "managed_organizations_url": "/v2/users/user-000/managed_organizations",
                "billing_managed_organizations_url": "/v2/users/user-000/billing_managed_organizations",
                "audited_organizations_url": "/v2/users/user-000/audited_organizations",
                "managed_spaces_url": "/v2/users/user-000/managed_spaces",
                "audited_spaces_url": "/v2/users/user-000/audited_spaces"
             }
          }
       ]
    }`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(uaaJSON))
})
