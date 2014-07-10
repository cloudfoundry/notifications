package cf_test

import (
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var testSpaceGuid = "test-space-guid"
var UsersEndpoint = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    err := req.ParseForm()
    if err != nil {
        panic(err)
    }

    if query := req.Form.Get("q"); query != "" {
        spaceGuid := strings.TrimPrefix(query, "space_guid:")
        if spaceGuid != testSpaceGuid {
            w.WriteHeader(http.StatusNotFound)
            w.Write([]byte(`{"total_results":0,"total_pages":1,"prev_url":null,"next_url":null,"resources":[]}`))
        }
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{
      "total_results": 1,
      "total_pages": 1,
      "prev_url": null,
      "next_url": null,
      "resources": [
        {
          "metadata": {
            "guid": "user-123",
            "url": "/v2/users/user-123",
            "created_at": "2013-04-30T21:00:49+00:00",
            "updated_at": null
          },
          "entity": {
            "admin": true,
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
        },
        {
          "metadata": {
            "guid": "user-456",
            "url": "/v2/users/user-456",
            "created_at": "2013-04-31T21:00:49+00:00",
            "updated_at": null
          },
          "entity": {
            "admin": true,
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
    }`))
})

var _ = Describe("CloudController", func() {
    var CCServer *httptest.Server

    Describe("GetUsersBySpaceGuid", func() {
        BeforeEach(func() {
            CCServer = httptest.NewServer(UsersEndpoint)
        })

        AfterEach(func() {
            CCServer.Close()
        })

        It("returns a list of users for the given space guid", func() {
            cloudController := cf.NewCloudController(CCServer.URL)
            users, err := cloudController.GetUsersBySpaceGuid(testSpaceGuid)
            if err != nil {
                panic(err)
            }

            Expect(len(users)).To(Equal(2))

            Expect(users).To(ContainElement(cf.CloudControllerUser{
                Guid: "user-123",
            }))

            Expect(users).To(ContainElement(cf.CloudControllerUser{
                Guid: "user-456",
            }))
        })
    })
})
