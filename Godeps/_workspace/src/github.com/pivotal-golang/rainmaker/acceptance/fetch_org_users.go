package acceptance

import (
	"os"

	"github.com/pivotal-golang/rainmaker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fetch all the users of an organization", func() {
	It("fetches the user records of all users associated with an organization", func() {
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

		list, err := client.Organizations.ListUsers(org.GUID, token)
		Expect(err).NotTo(HaveOccurred())

		Expect(list.Users).To(HaveLen(1))
		Expect(list.Users[0].GUID).To(Equal(user.GUID))
	})

	It("fetches paginated results of users associated with an organization", func() {
		token := os.Getenv("UAA_TOKEN")

		client := rainmaker.NewClient(rainmaker.Config{
			Host:          os.Getenv("CC_HOST"),
			SkipVerifySSL: true,
		})

		org, err := client.Organizations.Create(NewOrgName("org"), token)
		Expect(err).NotTo(HaveOccurred())

		for _ = range make([]int, 60, 60) {
			user, err := client.Users.Create(NewGUID("user"), token)
			Expect(err).NotTo(HaveOccurred())

			err = org.Users.Associate(user.GUID, token)
			Expect(err).NotTo(HaveOccurred())
		}

		list, err := client.Organizations.ListUsers(org.GUID, token)
		Expect(err).NotTo(HaveOccurred())

		Expect(list.TotalResults).To(Equal(60))
		Expect(list.TotalPages).To(Equal(2))
		Expect(list.Users).To(HaveLen(50))

		users, err := list.AllUsers(token)
		Expect(err).NotTo(HaveOccurred())

		Expect(users).To(HaveLen(60))
	})
})
