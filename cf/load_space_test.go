package cf_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/cf"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var SpacesEndpoint = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/v2/spaces/nacho-space" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code":1234,"description":"This is not allowed.","error_code":"CF-SpaceNotFound"}`))
		return
	} else if req.URL.Path != "/v2/spaces/space-guid" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"code":40004,"description":"The app space could not be found: ` + strings.TrimPrefix(req.URL.Path, "/v2/spaces/") + `","error_code":"CF-SpaceNotFound"}`))
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
	var cc cf.CloudController

	BeforeEach(func() {
		CCServer = httptest.NewServer(SpacesEndpoint)
		cc = cf.NewCloudController(CCServer.URL, false)
	})

	AfterEach(func() {
		CCServer.Close()
	})

	It("loads the space from cloud controller", func() {
		space, err := cc.LoadSpace("space-guid", "notification-token")
		if err != nil {
			panic(err)
		}

		Expect(space.GUID).To(Equal("space-guid"))
		Expect(space.Name).To(Equal("duh space"))
		Expect(space.OrganizationGUID).To(Equal("first-rate"))
	})

	It("returns a NotFoundError when the space cannot be found", func() {
		_, err := cc.LoadSpace("banana", "notification-token")
		Expect(err).To(BeAssignableToTypeOf(cf.NotFoundError{}))
		Expect(err.Error()).To(Equal(`CloudController Failure: Space "banana" could not be found`))
	})

	It("returns a 0 error code for any other error", func() {
		_, err := cc.LoadSpace("nacho-space", "notification-token")
		Expect(err).To(BeAssignableToTypeOf(cf.Failure{}))
		Expect(err.(cf.Failure).Code).To(Equal(0))
	})
})
