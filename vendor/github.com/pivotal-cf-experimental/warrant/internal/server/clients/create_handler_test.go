package clients_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/server/clients"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("createHandler", func() {
	var (
		router            http.Handler
		recorder          *httptest.ResponseRecorder
		request           *http.Request
		clientsCollection *domain.Clients
		tokensCollection  *domain.Tokens
	)

	BeforeEach(func() {
		clientsCollection = domain.NewClients()
		tokensCollection = NewTokens()
		router = clients.NewRouter(clientsCollection, tokensCollection)
		recorder = httptest.NewRecorder()

		requestBody, err := json.Marshal(map[string]interface{}{
			"client_id":              "some-client",
			"client_secret":          "some-client-secret",
			"name":                   "banana",
			"scope":                  []string{"some.scope"},
			"resource_ids":           []string{"resource"},
			"authorities":            []string{"cloud_controller.read"},
			"authorized_grant_types": []string{"client_credentials"},
			"access_token_validity":  43200,
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("POST", "/oauth/clients", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		authorization := tokensCollection.Encrypt(domain.Token{
			ClientID:    "the-client",
			Authorities: []string{"clients.write"},
			Audiences:   []string{"clients"},
		})
		request.Header.Set("Authorization", fmt.Sprintf("bearer %s", authorization))
	})

	It("creates a client", func() {
		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusCreated))
		Expect(recorder.Body).To(MatchJSON(`{
			"client_id": "some-client",
			"name":                   "banana",
			"scope": ["some.scope"],
			"resource_ids": ["resource"],
			"authorities": ["cloud_controller.read"],
			"authorized_grant_types": ["client_credentials"],
			"access_token_validity": 43200,
			"redirect_uri": null,
			"autoapprove": []
		}`))
	})

	It("stores the client in the collection", func() {
		router.ServeHTTP(recorder, request)
		client, ok := clientsCollection.Get("some-client")
		Expect(ok).To(BeTrue())
		Expect(client).To(Equal(domain.Client{
			ID:                   "some-client",
			Secret:               "some-client-secret",
			Name:                 "banana",
			Scope:                []string{"some.scope"},
			ResourceIDs:          []string{"resource"},
			Authorities:          []string{"cloud_controller.read"},
			AuthorizedGrantTypes: []string{"client_credentials"},
			AccessTokenValidity:  43200,
			RedirectURI:          nil,
			Autoapprove:          nil,
		}))
	})

	It("requires an authorization token", func() {
		request.Header.Del("Authorization")

		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		Expect(recorder.Body).To(MatchJSON(`{
			"error": "unauthorized",
			"error_description":"Full authentication is required to access this resource"
		}`))
	})

	It("requires an authorization token with the correct authorities", func() {
		authorization := tokensCollection.Encrypt(domain.Token{
			ClientID:    "the-client",
			Authorities: []string{"banana.write"},
			Audiences:   []string{"banana"},
		})
		request.Header.Set("Authorization", fmt.Sprintf("bearer %s", authorization))

		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		Expect(recorder.Body).To(MatchJSON(`{
			"error": "unauthorized",
			"error_description":"Full authentication is required to access this resource"
		}`))
	})

	It("validates the client create request", func() {
		requestBody, err := json.Marshal(map[string]interface{}{
			"client_id":              "some-client",
			"client_secret":          "some-client-secret",
			"name":                   "banana",
			"scope":                  []string{"some.scope"},
			"resource_ids":           []string{"resource"},
			"authorities":            []string{"cloud_controller.read"},
			"authorized_grant_types": []string{"bananas"},
			"access_token_validity":  43200,
		})
		Expect(err).NotTo(HaveOccurred())
		request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		Expect(recorder.Body).To(MatchJSON(`{
			"error_description": "bananas is not an allowed grant type. Must be one of: [implicit refresh_token authorization_code client_credentials password]",
			"error": "invalid_client"
		}`))
	})

	It("returns an error response when the request cannot be unmarshalled", func() {
		request.Body = ioutil.NopCloser(strings.NewReader("%%%"))

		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		Expect(recorder.Body).To(ContainSubstring("The request sent by the client was syntactically incorrect."))
	})
})
