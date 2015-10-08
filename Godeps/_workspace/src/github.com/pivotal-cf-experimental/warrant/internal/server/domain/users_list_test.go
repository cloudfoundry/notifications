package domain_test

import (
	"github.com/pivotal-cf-experimental/warrant/internal/documents"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UsersList", func() {
	Describe("ToDocument", func() {
		It("returns a list of resources when there is a user in the list", func() {
			list := domain.UsersList{domain.User{}}

			doc := list.ToDocument()
			Expect(doc.Resources).To(HaveLen(1))
		})

		It("returns an empty resources array when there are no users", func() {
			list := domain.UsersList{}

			doc := list.ToDocument()
			Expect(doc.Resources).To(Equal([]documents.UserResponse{}))
		})
	})
})
