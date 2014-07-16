package cf_test

import (
    "net/http"
    "net/http/httptest"

    "github.com/cloudfoundry-incubator/notifications/cf"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var SpacesEndpoint = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    if req.URL.Path != "/v2/spaces/space-guid" {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{
       "metadata": {
          "guid": "space-guid",
          "url": "/v2/spaces/7e3b8b40-cced-4714-8d4a-2b6bddc10fda",
          "created_at": "2014-06-25T18:25:05+00:00",
          "updated_at": null
       },
       "entity": {
          "name": "duh space",
          "organization_guid": "first-rate",
          "organization_url": "/v2/organizations/cd1d0c26-0da8-42d8-9478-8c1d32235279",
          "developers_url": "/v2/spaces/7e3b8b40-cced-4714-8d4a-2b6bddc10fda/developers",
          "managers_url": "/v2/spaces/7e3b8b40-cced-4714-8d4a-2b6bddc10fda/managers",
          "auditors_url": "/v2/spaces/7e3b8b40-cced-4714-8d4a-2b6bddc10fda/auditors",
          "apps_url": "/v2/spaces/7e3b8b40-cced-4714-8d4a-2b6bddc10fda/apps",
          "routes_url": "/v2/spaces/7e3b8b40-cced-4714-8d4a-2b6bddc10fda/routes",
          "domains_url": "/v2/spaces/7e3b8b40-cced-4714-8d4a-2b6bddc10fda/domains",
          "service_instances_url": "/v2/spaces/7e3b8b40-cced-4714-8d4a-2b6bddc10fda/service_instances",
          "app_events_url": "/v2/spaces/7e3b8b40-cced-4714-8d4a-2b6bddc10fda/app_events",
          "events_url": "/v2/spaces/7e3b8b40-cced-4714-8d4a-2b6bddc10fda/events",
          "security_groups_url": "/v2/spaces/7e3b8b40-cced-4714-8d4a-2b6bddc10fda/security_groups"
       }
    }`))
})

var _ = Describe("LoadSpace", func() {
    var CCServer *httptest.Server

    BeforeEach(func() {
        CCServer = httptest.NewServer(SpacesEndpoint)
    })

    AfterEach(func() {
        CCServer.Close()
    })

    It("loads the space from cloud controller", func() {
        cc := cf.NewCloudController(CCServer.URL)

        space, err := cc.LoadSpace("space-guid", "notification-token")
        if err != nil {
            panic(err)
        }

        Expect(space.Guid).To(Equal("space-guid"))
        Expect(space.Name).To(Equal("duh space"))
        Expect(space.OrganizationGuid).To(Equal("first-rate"))
    })
})
