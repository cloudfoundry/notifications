package acceptance

import (
	"os"

	"github.com/pivotal-golang/rainmaker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fetch all the managers of an organization", func() {
	It("fetches the user records of all managers associated with an organization", func() {
		token := os.Getenv("UAA_TOKEN")

		client := rainmaker.NewClient(rainmaker.Config{
			Host:          os.Getenv("CC_HOST"),
			SkipVerifySSL: true,
		})

		user, err := client.Users.Create(NewGUID("user"), token)
		Expect(err).NotTo(HaveOccurred())

		org, err := client.Organizations.Create(NewOrgName("org"), token)
		Expect(err).NotTo(HaveOccurred())

		err = org.Managers.Associate(user.GUID, token)
		Expect(err).NotTo(HaveOccurred())

		list, err := client.Organizations.ListManagers(org.GUID, token)
		Expect(err).NotTo(HaveOccurred())

		Expect(list.Users).To(HaveLen(1))
		Expect(list.Users[0].GUID).To(Equal(user.GUID))
	})
})
