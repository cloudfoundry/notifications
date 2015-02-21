package rainmaker_test

import (
	"time"

	"github.com/pivotal-cf-experimental/rainmaker"
	"github.com/pivotal-cf-experimental/rainmaker/internal/documents"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Space", func() {
	var config rainmaker.Config

	BeforeEach(func() {
		config = rainmaker.Config{
			Host: fakeCloudController.URL(),
		}
	})

	Describe("NewSpaceFromResponse", func() {
		var createdAt, updatedAt time.Time
		var document documents.SpaceResponse

		BeforeEach(func() {
			createdAt = time.Now().Add(-15 * time.Minute).UTC()
			updatedAt = time.Now().Add(-5 * time.Minute).UTC()

			document = documents.SpaceResponse{}
			document.Metadata.GUID = "space-001"
			document.Metadata.URL = "/v2/spaces/space-001"
			document.Metadata.CreatedAt = &createdAt
			document.Metadata.UpdatedAt = &updatedAt
			document.Entity.Name = "banana"
			document.Entity.OrganizationGUID = "org-001"
			document.Entity.SpaceQuotaDefinitionGUID = "space-quota-definition-guid"
			document.Entity.OrganizationURL = "/v2/organizations/org-001"
			document.Entity.DevelopersURL = "/v2/spaces/space-001/developers"
			document.Entity.ManagersURL = "/v2/spaces/space-001/managers"
			document.Entity.AuditorsURL = "/v2/spaces/space-001/auditors"
			document.Entity.AppsURL = "/v2/spaces/space-001/apps"
			document.Entity.RoutesURL = "/v2/spaces/space-001/routes"
			document.Entity.DomainsURL = "/v2/spaces/space-001/domains"
			document.Entity.ServiceInstancesURL = "/v2/spaces/space-001/service_instances"
			document.Entity.AppEventsURL = "/v2/spaces/space-001/app_events"
			document.Entity.EventsURL = "/v2/spaces/space-001/events"
			document.Entity.SecurityGroupsURL = "/v2/spaces/space-001/security_groups"
		})

		It("converts a response into a space", func() {
			space := rainmaker.NewSpaceFromResponse(config, document)

			expectedSpace := rainmaker.NewSpace(config, "space-001")
			expectedSpace.URL = "/v2/spaces/space-001"
			expectedSpace.Name = "banana"
			expectedSpace.OrganizationGUID = "org-001"
			expectedSpace.SpaceQuotaDefinitionGUID = "space-quota-definition-guid"
			expectedSpace.OrganizationURL = "/v2/organizations/org-001"
			expectedSpace.DevelopersURL = "/v2/spaces/space-001/developers"
			expectedSpace.ManagersURL = "/v2/spaces/space-001/managers"
			expectedSpace.AuditorsURL = "/v2/spaces/space-001/auditors"
			expectedSpace.AppsURL = "/v2/spaces/space-001/apps"
			expectedSpace.RoutesURL = "/v2/spaces/space-001/routes"
			expectedSpace.DomainsURL = "/v2/spaces/space-001/domains"
			expectedSpace.ServiceInstancesURL = "/v2/spaces/space-001/service_instances"
			expectedSpace.AppEventsURL = "/v2/spaces/space-001/app_events"
			expectedSpace.EventsURL = "/v2/spaces/space-001/events"
			expectedSpace.SecurityGroupsURL = "/v2/spaces/space-001/security_groups"
			expectedSpace.CreatedAt = createdAt
			expectedSpace.UpdatedAt = updatedAt

			Expect(space).To(Equal(expectedSpace))
		})
	})
})
