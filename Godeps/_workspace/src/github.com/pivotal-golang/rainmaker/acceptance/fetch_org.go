package acceptance

import (
	"os"

	"github.com/pivotal-golang/rainmaker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fetch an organization", func() {
	It("fetches the organization", func() {
		token := os.Getenv("UAA_TOKEN")

		client := rainmaker.NewClient(rainmaker.Config{
			Host:          os.Getenv("CC_HOST"),
			SkipVerifySSL: true,
		})

		org, err := client.Organizations.Create(NewOrgName("org"), token)
		Expect(err).NotTo(HaveOccurred())

		fetchedOrg, err := client.Organizations.Get(org.GUID, token)
		Expect(err).NotTo(HaveOccurred())
		Expect(fetchedOrg).To(Equal(org))
	})
})
