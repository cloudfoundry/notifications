package servers

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "os"

    "github.com/gorilla/mux"
)

type CC struct {
    server *httptest.Server
}

func NewCC() CC {
    router := mux.NewRouter()
    router.HandleFunc("/v2/spaces/{guid}", CCGetSpace).Methods("GET")
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

var CCGetSpaceUsers = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    json := `{
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
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(json))
})
