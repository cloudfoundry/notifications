package clients_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/pivotal-cf-experimental/warrant/internal/server/clients"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("getHandler", func() {
	var (
		router           http.Handler
		recorder         *httptest.ResponseRecorder
		request          *http.Request
		tokensCollection *domain.Tokens
	)

	BeforeEach(func() {
		tokensCollection = NewTokens()
		token := tokensCollection.Encrypt(domain.Token{
			ClientID:    "my-client-id",
			Authorities: []string{"clients.read"},
			Audiences:   []string{"clients"},
		})

		clientsCollection := domain.NewClients()
		clientsCollection.Add(domain.Client{
			ID:                   "some-client-id",
			Name:                 "banana",
			Scope:                []string{"some-scope"},
			ResourceIDs:          []string{"some-resource-id"},
			Authorities:          []string{"some-authority"},
			AuthorizedGrantTypes: []string{"some-grant-type"},
			AccessTokenValidity:  3600,
			RedirectURI:          []string{"https://example.com/sessions/create"},
			Autoapprove:          []string{"some-approval"},
		})

		var err error
		request, err = http.NewRequest("GET", "/oauth/clients/some-client-id", nil)
		Expect(err).NotTo(HaveOccurred())
		request.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))

		recorder = httptest.NewRecorder()
		router = clients.NewRouter(clientsCollection, tokensCollection)
	})

	It("gets the requested client", func() {
		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"client_id": "some-client-id",
			"name": "banana",
			"scope": ["some-scope"],
			"resource_ids": ["some-resource-id"],
			"authorities": ["some-authority"],
			"authorized_grant_types": ["some-grant-type"],
			"autoapprove": ["some-approval"],
			"access_token_validity": 3600,
			"redirect_uri": ["https://example.com/sessions/create"]
		}`))
	})

	It("returns a 404 if the client cannot be found", func() {
		request.URL.Path = "/oauth/clients/missing-client"

		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusNotFound))
		Expect(recorder.Body).To(MatchJSON(`{
			"error_description": "Client missing-client does not exist",
			"error": "scim_resource_not_found"
		}`))
	})

	It("requires a token", func() {
		request.Header.Del("Authorization")

		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		Expect(recorder.Body).To(MatchJSON(`{
			"error_description": "Full authentication is required to access this resource",
			"error": "unauthorized"
		}`))
	})

	It("requires a token with the correct scopes", func() {
		token := tokensCollection.Encrypt(domain.Token{
			ClientID:    "my-client-id",
			Authorities: []string{"bananas.read"},
			Audiences:   []string{"bananas"},
		})
		request.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))

		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		Expect(recorder.Body).To(MatchJSON(`{
			"error_description": "Full authentication is required to access this resource",
			"error": "unauthorized"
		}`))
	})
})
