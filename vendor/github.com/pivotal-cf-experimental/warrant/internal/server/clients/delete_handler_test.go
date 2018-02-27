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

var _ = Describe("deleteHandler", func() {
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

		clientsCollection.Add(domain.Client{
			ID: "some-client",
		})

		var err error
		request, err = http.NewRequest("DELETE", "/oauth/clients/some-client", nil)
		Expect(err).NotTo(HaveOccurred())

		token := tokensCollection.Encrypt(domain.Token{
			Authorities: []string{"clients.write"},
			Audiences:   []string{"clients"},
		})
		request.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	})

	It("deletes the client", func() {
		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusOK))

		_, ok := clientsCollection.Get("some-client")
		Expect(ok).To(BeFalse())
	})

	It("requires an Authorization header", func() {
		request.Header.Del("Authorization")

		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
		Expect(recorder.Body).To(MatchJSON(`{
			"error": "unauthorized",
			"error_description": "Full authentication is required to access this resource"
		}`))
	})

	It("requires a token with valid permissions", func() {
		token := tokensCollection.Encrypt(domain.Token{
			Authorities: []string{"bananas.eat"},
			Audiences:   []string{"bananas"},
		})
		request.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))

		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusForbidden))
		Expect(recorder.Body).To(MatchJSON(`{
			"error": "access_denied",
			"error_description": "Invalid token does not contain resource id (clients)"
		}`))
	})
})
