package acceptance

import (
	"time"

	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client Lifecycle", func() {
	var (
		warrantClient warrant.Warrant
		client        warrant.Client
	)

	BeforeEach(func() {
		client = warrant.Client{
			ID:                   UAADefaultClientID,
			Scope:                []string{"openid"},
			ResourceIDs:          []string{"none"},
			Authorities:          []string{"scim.read", "scim.write"},
			AuthorizedGrantTypes: []string{"client_credentials"},
			AccessTokenValidity:  5000 * time.Second,
		}

		warrantClient = warrant.New(warrant.Config{
			Host:          UAAHost,
			SkipVerifySSL: true,
			TraceWriter:   TraceWriter,
		})
	})

	AfterEach(func() {
		err := warrantClient.Clients.Delete(client.ID, UAAToken)
		Expect(err).NotTo(HaveOccurred())
	})

	It("creates, and retrieves a client", func() {
		By("creating a client", func() {
			err := warrantClient.Clients.Create(client, "secret", UAAToken)
			Expect(err).NotTo(HaveOccurred())
		})

		By("finding the client", func() {
			fetchedClient, err := warrantClient.Clients.Get(client.ID, UAAToken)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedClient).To(Equal(client))
		})
	})

	It("rejects requests from clients that do not have clients.write scope", func() {
		var token string

		By("creating a client", func() {
			err := warrantClient.Clients.Create(client, "secret", UAAToken)
			Expect(err).NotTo(HaveOccurred())
		})

		By("fetching the new client token", func() {
			var err error

			token, err = warrantClient.Clients.GetToken(client.ID, "secret")
			Expect(err).NotTo(HaveOccurred())
		})

		By("using the new client token to delete the client", func() {
			err := warrantClient.Clients.Delete(client.ID, token)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(warrant.UnauthorizedError{}))
		})
	})
})
