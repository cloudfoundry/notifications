package acceptance

import (
	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Groups", func() {
	var (
		client warrant.Warrant
		group  warrant.Group
		err    error
	)

	BeforeEach(func() {
		client = warrant.New(warrant.Config{
			Host:          UAAHost,
			SkipVerifySSL: true,
			TraceWriter:   TraceWriter,
		})
	})

	AfterEach(func() {
		err = client.Groups.Delete(group.ID, UAAToken)
		Expect(err).NotTo(HaveOccurred())

		_, err = client.Groups.Get(group.ID, UAAToken)
		Expect(err).To(BeAssignableToTypeOf(warrant.NotFoundError{}))
	})

	It("creates, lists, retrieves, and deletes a group", func() {
		By("creating a new group", func() {
			group, err = client.Groups.Create("banana.read", UAAToken)
			Expect(err).NotTo(HaveOccurred())

			Expect(group.ID).NotTo(BeEmpty())
			Expect(group.DisplayName).To(Equal("banana.read"))
		})

		By("listing the groups", func() {
			groups, err := client.Groups.List(warrant.Query{}, UAAToken)
			Expect(err).NotTo(HaveOccurred())
			Expect(groups).To(ContainElement(group))
		})

		By("getting the created group", func() {
			group, err = client.Groups.Get(group.ID, UAAToken)
			Expect(err).NotTo(HaveOccurred())

			Expect(group.ID).NotTo(BeEmpty())
			Expect(group.DisplayName).To(Equal("banana.read"))
		})
	})
})
