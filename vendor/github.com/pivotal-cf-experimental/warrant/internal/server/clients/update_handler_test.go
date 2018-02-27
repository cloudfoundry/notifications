package clients_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/pivotal-cf-experimental/warrant/internal/server/clients"
	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("updateHandler", func() {
	var (
		router            http.Handler
		recorder          *httptest.ResponseRecorder
		request           *http.Request
		tokensCollection  *domain.Tokens
		clientsCollection *domain.Clients
	)

	BeforeEach(func() {
		tokensCollection = domain.NewTokens(common.TestPublicKey, common.TestPrivateKey, []string{})

		clientsCollection = domain.NewClients()
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

		requestBody, err := json.Marshal(map[string]interface{}{
			"client_id":              "updated-client-id",
			"name":                   "banana",
			"scope":                  []string{"updated-scope"},
			"resource_ids":           []string{"updated-resource-id"},
			"authorities":            []string{"updated-authority"},
			"authorized_grant_types": []string{"authorization_code"},
			"autoapprove":            []string{"updated-approval"},
			"access_token_validity":  7200,
			"redirect_uri":           []string{"https://example.com/sessions/update"},
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("PUT", "/oauth/clients/some-client-id", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())
		token := tokensCollection.Encrypt(domain.Token{
			ClientID:    "my-client-id",
			Authorities: []string{"clients.write"},
			Audiences:   []string{"clients"},
		})
		request.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))

		recorder = httptest.NewRecorder()
		router = clients.NewRouter(clientsCollection, tokensCollection)
	})

	It("updates the requested client", func() {
		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body).To(MatchJSON(`{
			"client_id": "updated-client-id",
			"name": "banana",
			"scope": ["updated-scope"],
			"resource_ids": ["updated-resource-id"],
			"authorities": ["updated-authority"],
			"authorized_grant_types": ["authorization_code"],
			"autoapprove": ["updated-approval"],
			"access_token_validity": 7200,
			"redirect_uri": ["https://example.com/sessions/update"]
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

	It("requires a token with the correct permissions", func() {
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

	It("validates the client updates", func() {
		requestBody, err := json.Marshal(map[string]interface{}{
			"client_id":              "updated-client-id",
			"name":                   "banana",
			"scope":                  []string{"updated-scope"},
			"resource_ids":           []string{"updated-resource-id"},
			"authorities":            []string{"updated-authority"},
			"authorized_grant_types": []string{"invalid-grant-type"},
			"autoapprove":            []string{"updated-approval"},
			"access_token_validity":  7200,
			"redirect_uri":           []string{"https://example.com/sessions/update"},
		})
		Expect(err).NotTo(HaveOccurred())
		request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		Expect(recorder.Body).To(MatchJSON(`{
			"error_description": "invalid-grant-type is not an allowed grant type. Must be one of: [implicit refresh_token authorization_code client_credentials password]",
			"error": "invalid_client"
		}`))
	})

	It("stores the update client in the clients collection", func() {
		router.ServeHTTP(recorder, request)
		client, ok := clientsCollection.Get("updated-client-id")
		Expect(ok).To(BeTrue())
		Expect(client).To(Equal(domain.Client{
			ID:                   "updated-client-id",
			Name:                 "banana",
			Scope:                []string{"updated-scope"},
			ResourceIDs:          []string{"updated-resource-id"},
			Authorities:          []string{"updated-authority"},
			AuthorizedGrantTypes: []string{"authorization_code"},
			AccessTokenValidity:  7200,
			RedirectURI:          []string{"https://example.com/sessions/update"},
			Autoapprove:          []string{"updated-approval"},
		}))

		_, ok = clientsCollection.Get("some-client-id")
		Expect(ok).To(BeFalse())
	})
})
