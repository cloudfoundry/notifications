package acceptance

import (
	"os"

	"github.com/pivotal-cf-experimental/rainmaker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fetch a space", func() {
	It("fetches the space", func() {
		token := os.Getenv("UAA_TOKEN")

		client := rainmaker.NewClient(rainmaker.Config{
			Host:          os.Getenv("CC_HOST"),
			SkipVerifySSL: true,
		})

		org, err := client.Organizations.Create(NewOrgName("org"), token)
		Expect(err).NotTo(HaveOccurred())

		space, err := client.Spaces.Create("space-1", org.GUID, token)

		fetchedSpace, err := client.Spaces.Get(space.GUID, token)
		Expect(err).NotTo(HaveOccurred())
		Expect(fetchedSpace).To(Equal(space))
	})
})
