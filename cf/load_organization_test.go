package cf_test

import (
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/cloudfoundry-incubator/notifications/cf"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var OrganizationsEndpoint = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    if req.URL.Path != "/v2/organizations/org-guid" {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"code": 30003, "description": "The organization could not be found: ` + strings.TrimPrefix(req.URL.Path, "/v2/organizations/") + `", "error_code": "CF-OrganizationNotFound"}`))
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{
       "metadata": {
          "guid": "org-guid",
          "url": "/v2/organizations/org-guid",
          "created_at": "2014-06-25T18:24:48+00:00",
          "updated_at": "2014-06-25T18:27:13+00:00"
       },
       "entity": {
          "name": "Initech",
          "billing_enabled": false,
          "quota_definition_guid": "caf592e1-bdac-40ec-b863-5d44ab785b7e",
          "status": "active",
          "quota_definition_url": "/v2/quota_definitions/caf592e1-bdac-40ec-b863-5d44ab785b7e",
          "spaces_url": "/v2/organizations/cd1d0c26-0da8-42d8-9478-8c1d32235279/spaces",
          "domains_url": "/v2/organizations/cd1d0c26-0da8-42d8-9478-8c1d32235279/domains",
          "private_domains_url": "/v2/organizations/cd1d0c26-0da8-42d8-9478-8c1d32235279/private_domains",
          "users_url": "/v2/organizations/cd1d0c26-0da8-42d8-9478-8c1d32235279/users",
          "managers_url": "/v2/organizations/cd1d0c26-0da8-42d8-9478-8c1d32235279/managers",
          "billing_managers_url": "/v2/organizations/cd1d0c26-0da8-42d8-9478-8c1d32235279/billing_managers",
          "auditors_url": "/v2/organizations/cd1d0c26-0da8-42d8-9478-8c1d32235279/auditors",
          "app_events_url": "/v2/organizations/cd1d0c26-0da8-42d8-9478-8c1d32235279/app_events"
       }
    }`))
})

var _ = Describe("LoadOrganization", func() {
    var CCServer *httptest.Server
    var cc cf.CloudController

    BeforeEach(func() {
        CCServer = httptest.NewServer(OrganizationsEndpoint)
        cc = cf.NewCloudController(CCServer.URL, false)
    })

    AfterEach(func() {
        CCServer.Close()
    })

    It("loads the organization from cloud controller", func() {

        org, err := cc.LoadOrganization("org-guid", "notification-token")
        if err != nil {
            panic(err)
        }

        Expect(org.Guid).To(Equal("org-guid"))
        Expect(org.Name).To(Equal("Initech"))
    })

    It("returns a Failure instance when the org cannot be found", func() {
        _, err := cc.LoadOrganization("banana", "notification-token")

        Expect(err).To(BeAssignableToTypeOf(cf.Failure{}))

        failure := err.(cf.Failure)
        Expect(failure.Code).To(Equal(http.StatusNotFound))
        Expect(failure.Message).To(Equal(`{"code": 30003, "description": "The organization could not be found: banana", "error_code": "CF-OrganizationNotFound"}`))
        Expect(failure.Error()).To(Equal(`CloudController Failure (404): {"code": 30003, "description": "The organization could not be found: banana", "error_code": "CF-OrganizationNotFound"}`))
    })
})
