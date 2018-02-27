package domain_test

import (
	"github.com/pivotal-cf-experimental/warrant/internal/documents"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientsList", func() {
	Describe("ToDocument", func() {
		It("returns a list of resources when there is a client in the list", func() {
			list := domain.ClientsList{domain.Client{}}

			doc := list.ToDocument()
			Expect(doc.Resources).To(HaveLen(1))
		})

		It("returns an empty resources array when there are no clients", func() {
			list := domain.ClientsList{}

			doc := list.ToDocument()
			Expect(doc.Resources).To(Equal([]documents.ClientResponse{}))
		})
	})
})
