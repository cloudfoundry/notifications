package tokens_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
	"github.com/pivotal-cf-experimental/warrant/internal/server/tokens"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type hasURL struct{}

func (hasURL) URL() string {
	return "https://uaa.example.com"
}

var _ = Describe("authorizeHandler", func() {
	var (
		router           http.Handler
		recorder         *httptest.ResponseRecorder
		request          *http.Request
		tokensCollection *domain.Tokens
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()

		values := url.Values{
			"client_id":     {"some-client-id"},
			"scope":         {"openid"},
			"response_type": {"token"},
			"redirect_uri":  {"https://uaa.example.com"},
		}
		u := &url.URL{
			Path:     "/oauth/authorize",
			RawQuery: values.Encode(),
		}

		var err error
		request, err = http.NewRequest("POST", u.String(), strings.NewReader(url.Values{
			"username": {"some-user"},
			"password": {"password"},
			"source":   {"credentials"},
		}.Encode()))
		Expect(err).NotTo(HaveOccurred())

		request.Header.Set("Accept", "application/json")
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		tokensCollection = domain.NewTokens(common.TestPublicKey, common.TestPrivateKey, []string{})
		usersCollection := domain.NewUsers()
		clientsCollection := domain.NewClients()

		clientsCollection.Add(domain.NewClientFromDocument(documents.CreateUpdateClientRequest{
			ClientID: "some-client-id",
		}))

		usersCollection.Add(domain.NewUserFromUpdateDocument(documents.UpdateUserRequest{
			ID:       "some-user",
			UserName: "some-user",
		}))

		user, ok := usersCollection.Get("some-user")
		Expect(ok).To(BeTrue())

		user.Password = "password"
		usersCollection.Add(user)

		router = tokens.NewRouter(tokensCollection,
			usersCollection, clientsCollection, common.TestPublicKey, common.TestPrivateKey, hasURL{})
	})

	It("returns a valid token when there is no overlap between client and user scopes", func() {
		router.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusFound))

		location, err := url.Parse(recorder.HeaderMap.Get("Location"))
		Expect(err).NotTo(HaveOccurred())

		query, err := url.ParseQuery(location.Fragment)
		Expect(err).NotTo(HaveOccurred())
		Expect(query.Get("scope")).To(Equal(""))

		token, err := tokensCollection.Decrypt(query.Get("access_token"))
		Expect(err).NotTo(HaveOccurred())
		Expect(token.Scopes).NotTo(BeNil())
	})
})
