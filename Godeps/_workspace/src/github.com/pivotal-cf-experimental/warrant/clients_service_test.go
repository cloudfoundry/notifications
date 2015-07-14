package warrant_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientsService", func() {
	var (
		service warrant.ClientsService
		token   string
		config  warrant.Config
	)

	BeforeEach(func() {
		config = warrant.Config{
			Host:          fakeUAA.URL(),
			SkipVerifySSL: true,
			TraceWriter:   TraceWriter,
		}
		service = warrant.NewClientsService(config)
		token = fakeUAA.ClientTokenFor("admin", []string{"clients.write", "clients.read"}, []string{"clients"})
	})

	Describe("Create/Get", func() {
		var client warrant.Client

		BeforeEach(func() {
			client = warrant.Client{
				ID:                   "client-id",
				Scope:                []string{"openid"},
				ResourceIDs:          []string{"none"},
				Authorities:          []string{"scim.read", "scim.write"},
				AuthorizedGrantTypes: []string{"client_credentials"},
				AccessTokenValidity:  5000 * time.Second,
			}
		})

		It("an error does not occur and the new client can be fetched", func() {
			err := service.Create(client, "client-secret", token)
			Expect(err).NotTo(HaveOccurred())

			foundClient, err := service.Get(client.ID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(foundClient).To(Equal(client))
		})

		It("responds with an error when the client cannot be created", func() {
			client.AuthorizedGrantTypes = []string{"invalid-grant-type"}
			err := service.Create(client, "client-secret", token)
			Expect(err).To(BeAssignableToTypeOf(warrant.BadRequestError{}))
			Expect(err.Error()).To(Equal(`bad request: {"message":"invalid-grant-type is not an allowed grant type. Must be one of: [implicit refresh_token authorization_code client_credentials password]","error":"invalid_client"}`))
		})

		It("responds with an error when the client cannot be found", func() {
			_, err := service.Get("unknown-client", token)
			Expect(err).To(BeAssignableToTypeOf(warrant.NotFoundError{}))
		})

		Context("failure cases", func() {
			It("returns an error if the json response is malformed", func() {
				malformedJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write([]byte("this is not JSON"))
				}))
				service = warrant.NewClientsService(warrant.Config{
					Host:          malformedJSONServer.URL,
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				_, err := service.Get("some-client", "some-token")
				Expect(err).To(BeAssignableToTypeOf(warrant.MalformedResponseError{}))
				Expect(err).To(MatchError("malformed response: invalid character 'h' in literal true (expecting 'r')"))
			})
		})
	})

	Describe("GetToken", func() {
		var (
			client       warrant.Client
			clientSecret string
		)

		BeforeEach(func() {
			client = warrant.Client{
				ID:                   "client-id",
				Scope:                []string{"openid"},
				ResourceIDs:          []string{"none"},
				Authorities:          []string{"scim.read", "scim.write"},
				AuthorizedGrantTypes: []string{"client_credentials"},
				AccessTokenValidity:  5000 * time.Second,
			}
			clientSecret = "client-secret"

			err := service.Create(client, clientSecret, token)
			Expect(err).NotTo(HaveOccurred())
		})

		It("retrieves a token for the client given a valid secret", func() {
			clientToken, err := service.GetToken(client.ID, clientSecret)
			Expect(err).NotTo(HaveOccurred())

			tokensService := warrant.NewTokensService(config)
			decodedToken, err := tokensService.Decode(clientToken)
			Expect(err).NotTo(HaveOccurred())
			Expect(decodedToken.ClientID).To(Equal(client.ID))
		})

		Context("failure cases", func() {
			It("returns an error if the json response is malformed", func() {
				malformedJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write([]byte("this is not JSON"))
				}))
				service = warrant.NewClientsService(warrant.Config{
					Host:          malformedJSONServer.URL,
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				_, err := service.GetToken("some-client", "some-secret")
				Expect(err).To(BeAssignableToTypeOf(warrant.MalformedResponseError{}))
				Expect(err).To(MatchError("malformed response: invalid character 'h' in literal true (expecting 'r')"))
			})
		})
	})

	Describe("Delete", func() {
		var client warrant.Client

		BeforeEach(func() {
			client = warrant.Client{
				ID:                   "client-id",
				Scope:                []string{"openid"},
				ResourceIDs:          []string{"none"},
				Authorities:          []string{"scim.read", "scim.write"},
				AuthorizedGrantTypes: []string{"client_credentials"},
				AccessTokenValidity:  5000 * time.Second,
			}

			err := service.Create(client, "secret", token)
			Expect(err).NotTo(HaveOccurred())
		})

		It("deletes the client", func() {
			err := service.Delete(client.ID, token)
			Expect(err).NotTo(HaveOccurred())

			_, err = service.Get(client.ID, token)
			Expect(err).To(BeAssignableToTypeOf(warrant.NotFoundError{}))
		})

		It("errors when the token is unauthorized", func() {
			token = fakeUAA.ClientTokenFor("admin", []string{"clients.foo", "clients.boo"}, []string{"clients"})
			err := service.Delete(client.ID, token)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(warrant.UnauthorizedError{}))
		})
	})
})
