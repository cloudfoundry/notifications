package rainmaker_test

import (
	"time"

	"github.com/pivotal-golang/rainmaker"
	"github.com/pivotal-golang/rainmaker/internal/documents"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Organization", func() {
	var config rainmaker.Config

	Describe("NewOrganizationFromResponse", func() {
		BeforeEach(func() {
			config = rainmaker.Config{
				Host: fakeCloudController.URL(),
			}
		})

		It("converts a response into an organization", func() {
			createdAt := time.Now().Add(-10 * time.Minute)
			updatedAt := time.Now().Add(-9 * time.Minute)

			document := documents.OrganizationResponse{}
			document.Metadata.GUID = "org-001"
			document.Metadata.URL = "/v2/organizations/org-001"
			document.Metadata.CreatedAt = &createdAt
			document.Metadata.UpdatedAt = &updatedAt
			document.Entity.Name = "rainmaker-organization"
			document.Entity.BillingEnabled = true
			document.Entity.Status = "active"
			document.Entity.QuotaDefinitionGUID = "quota-definition-guid"
			document.Entity.QuotaDefinitionURL = "/v2/quota_definitions/quota-definition-guid"
			document.Entity.SpacesURL = "/v2/organizations/org-001/spaces"
			document.Entity.DomainsURL = "/v2/organizations/org-001/domains"
			document.Entity.PrivateDomainsURL = "/v2/organizations/org-001/private_domains"
			document.Entity.UsersURL = "/v2/organizations/org-001/users"
			document.Entity.ManagersURL = "/v2/organizations/org-001/managers"
			document.Entity.BillingManagersURL = "/v2/organizations/org-001/billing_managers"
			document.Entity.AuditorsURL = "/v2/organizations/org-001/auditors"
			document.Entity.AppEventsURL = "/v2/organizations/org-001/app_events"
			document.Entity.SpaceQuotaDefinitionsURL = "/v2/organizations/org-001/space_quota_definitions"

			organization := rainmaker.NewOrganizationFromResponse(config, document)
			expectedOrganization := rainmaker.NewOrganization(config, "org-001")
			expectedOrganization.Name = "rainmaker-organization"
			expectedOrganization.URL = "/v2/organizations/org-001"
			expectedOrganization.BillingEnabled = true
			expectedOrganization.Status = "active"
			expectedOrganization.QuotaDefinitionGUID = "quota-definition-guid"
			expectedOrganization.QuotaDefinitionURL = "/v2/quota_definitions/quota-definition-guid"
			expectedOrganization.SpacesURL = "/v2/organizations/org-001/spaces"
			expectedOrganization.DomainsURL = "/v2/organizations/org-001/domains"
			expectedOrganization.PrivateDomainsURL = "/v2/organizations/org-001/private_domains"
			expectedOrganization.UsersURL = "/v2/organizations/org-001/users"
			expectedOrganization.ManagersURL = "/v2/organizations/org-001/managers"
			expectedOrganization.BillingManagersURL = "/v2/organizations/org-001/billing_managers"
			expectedOrganization.AuditorsURL = "/v2/organizations/org-001/auditors"
			expectedOrganization.AppEventsURL = "/v2/organizations/org-001/app_events"
			expectedOrganization.SpaceQuotaDefinitionsURL = "/v2/organizations/org-001/space_quota_definitions"
			expectedOrganization.CreatedAt = createdAt
			expectedOrganization.UpdatedAt = updatedAt

			Expect(organization).To(Equal(expectedOrganization))
		})
	})
})
