package cf_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/cf"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudController", func() {
	var testOrganizationGuid = "test-organization-guid"
	var CCServer *httptest.Server
	var UsersEndpoint http.HandlerFunc
	var cloudController cf.CloudController

	Describe("GetUsersByOrgGuid", func() {
		BeforeEach(func() {
			UsersEndpoint = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
				if token != testUAAToken {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"code":10002,"description":"Authentication error","error_code":"CF-NotAuthenticated"}`))
					return
				}

				err := req.ParseForm()
				if err != nil {
					panic(err)
				}

				organizationGuid := strings.Split(req.URL.String(), "/")[3]
				if organizationGuid != testOrganizationGuid {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"total_results":0,"total_pages":1,"prev_url":null,"next_url":null,"resources":[]}`))
					return
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
                    }
                  ]
            }`))
			})

			CCServer = httptest.NewServer(UsersEndpoint)
			cloudController = cf.NewCloudController(CCServer.URL, false)
		})

		AfterEach(func() {
			CCServer.Close()
		})

		It("returns a list of users for the given organization guid", func() {
			users, err := cloudController.GetUsersByOrgGuid(testOrganizationGuid, testUAAToken)
			if err != nil {
				panic(err)
			}

			Expect(len(users)).To(Equal(1))

			Expect(users).To(ContainElement(cf.CloudControllerUser{
				GUID: "user-123",
			}))
		})

		It("returns an error when the Cloud Controller returns an error status code", func() {
			_, err := cloudController.GetUsersByOrgGuid("my-nonexistant-guid", testUAAToken)

			Expect(err).To(BeAssignableToTypeOf(cf.Failure{}))
		})
	})
})
