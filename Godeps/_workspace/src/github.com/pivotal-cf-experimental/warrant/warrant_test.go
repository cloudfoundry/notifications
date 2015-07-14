package warrant_test

import (
	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Warrant", func() {
	var client warrant.Warrant

	BeforeEach(func() {
		client = warrant.New(warrant.Config{})
	})

	It("has a users service", func() {
		Expect(client.Users).To(BeAssignableToTypeOf(warrant.UsersService{}))
	})

	It("has an clients service", func() {
		Expect(client.Clients).To(BeAssignableToTypeOf(warrant.ClientsService{}))
	})

	It("has a tokens service", func() {
		Expect(client.Tokens).To(BeAssignableToTypeOf(warrant.TokensService{}))
	})

	It("has a groups service", func() {
		Expect(client.Groups).To(BeAssignableToTypeOf(warrant.GroupsService{}))
	})
})
