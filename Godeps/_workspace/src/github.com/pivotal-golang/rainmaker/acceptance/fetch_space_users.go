package acceptance

import (
	"os"

	"github.com/pivotal-golang/rainmaker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fetch all the users of a space", func() {
	It("fetches the user records of all users associated with a space", func() {
		token := os.Getenv("UAA_TOKEN")

		client := rainmaker.NewClient(rainmaker.Config{
			Host:          os.Getenv("CC_HOST"),
			SkipVerifySSL: true,
		})

		user, err := client.Users.Create(NewGUID("user"), token)
		Expect(err).NotTo(HaveOccurred())

		org, err := client.Organizations.Create(NewOrgName("org"), token)
		Expect(err).NotTo(HaveOccurred())

		err = org.Users.Associate(user.GUID, token)
		Expect(err).NotTo(HaveOccurred())

		space, err := client.Spaces.Create("my-space", org.GUID, token)
		Expect(err).NotTo(HaveOccurred())

		err = space.Developers.Associate(user.GUID, token)
		Expect(err).NotTo(HaveOccurred())

		list, err := client.Spaces.ListUsers(space.GUID, token)
		Expect(err).NotTo(HaveOccurred())

		Expect(list.Users).To(HaveLen(1))
		Expect(list.Users[0].GUID).To(Equal(user.GUID))
	})
})
