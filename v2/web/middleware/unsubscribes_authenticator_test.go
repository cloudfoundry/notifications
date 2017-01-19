package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/web/middleware"
	"github.com/ryanmoran/stack"

	"errors"

	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unsubscribes Authenticator", func() {
	var (
		auth    middleware.UnsubscribesAuthenticator
		request *http.Request
		writer  *httptest.ResponseRecorder
		context stack.Context

		clientAuthenticator *mocks.Authenticator
		userAuthenticator   *mocks.Authenticator
		validator           *mocks.TokenValidator
		expectedToken       *jwt.Token
	)

	BeforeEach(func() {
		var err error

		writer = httptest.NewRecorder()

		request, err = http.NewRequest("GET", "/some/path", nil)
		if err != nil {
			panic(err)
		}
		request.Header.Set("Authorization", "Bearer valid-token")

		context = stack.NewContext()

		clientAuthenticator = &mocks.Authenticator{}
		userAuthenticator = &mocks.Authenticator{}
		validator = &mocks.TokenValidator{}

		auth = middleware.UnsubscribesAuthenticator{
			Validator:           validator,
			ClientAuthenticator: clientAuthenticator,
			UserAuthenticator:   userAuthenticator,
		}

		expectedToken = jwt.New(jwt.SigningMethodRS256)
		expectedToken.Claims = map[string]interface{}{
			"jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
			"client_id": "mister-client",
			"cid":       "mister-client",
			"exp":       3404281214,
			"scope":     []interface{}{"gaben.scope"},
		}
		validator.ParseCall.Returns.Token = expectedToken
	})

	Context("when the token is a client token", func() {
		It("delegates to the clientAuthenticator", func() {
			clientAuthenticator.ServeHTTPCall.Returns.Continue = true

			keepGoing := auth.ServeHTTP(writer, request, context)
			Expect(clientAuthenticator.ServeHTTPCall.Receives.Writer).To(Equal(writer))
			Expect(clientAuthenticator.ServeHTTPCall.Receives.Request).To(Equal(request))
			Expect(clientAuthenticator.ServeHTTPCall.Receives.Context).To(Equal(context))

			Expect(keepGoing).To(BeTrue())

			Expect(userAuthenticator.ServeHTTPCall.Receives.Request).To(BeNil())
		})
	})

	Context("when the token is a user token", func() {
		BeforeEach(func() {
			expectedToken = jwt.New(jwt.SigningMethodRS256)
			expectedToken.Claims = map[string]interface{}{
				"jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
				"client_id": "mister-client",
				"user_id":   "some-user-guid",
				"cid":       "mister-client",
				"exp":       3404281214,
				"scope":     []interface{}{"gaben.scope"},
			}

			validator.ParseCall.Returns.Token = expectedToken

			var err error
			request, err = http.NewRequest("GET", "/some/path/some-user-guid", nil)
			if err != nil {
				panic(err)
			}
			request.Header.Set("Authorization", "Bearer some-token")
		})

		Context("when the user_id in the token matches the user_guid in the route", func() {
			It("delegates to the userAuthenticator", func() {
				userAuthenticator.ServeHTTPCall.Returns.Continue = true

				keepGoing := auth.ServeHTTP(writer, request, context)
				Expect(userAuthenticator.ServeHTTPCall.Receives.Writer).To(Equal(writer))
				Expect(userAuthenticator.ServeHTTPCall.Receives.Request).To(Equal(request))
				Expect(userAuthenticator.ServeHTTPCall.Receives.Context).To(Equal(context))

				Expect(keepGoing).To(BeTrue())

				Expect(clientAuthenticator.ServeHTTPCall.Receives.Request).To(BeNil())
			})
		})

		Context("when the user_id does not match the user_guid in the route", func() {
			It("returns a 403 status code and error message", func() {
				var err error
				request.URL, err = url.Parse("/some/path/a-different-user")
				Expect(err).NotTo(HaveOccurred())

				keepGoing := auth.ServeHTTP(writer, request, context)
				Expect(userAuthenticator.ServeHTTPCall.Receives.Request).To(BeNil())
				Expect(clientAuthenticator.ServeHTTPCall.Receives.Request).To(BeNil())
				Expect(keepGoing).To(BeFalse())

				Expect(writer.Code).To(Equal(http.StatusForbidden))
				Expect(writer.Body.String()).To(MatchJSON(`{
					"errors": [
						"You are not authorized to perform the requested action"
					]
				}`))
			})
		})
	})

	Context("when the token is missing", func() {
		It("returns a 401 status code and error message", func() {
			request.Header.Set("Authorization", "")

			keepGoing := auth.ServeHTTP(writer, request, context)
			Expect(userAuthenticator.ServeHTTPCall.Receives.Request).To(BeNil())
			Expect(clientAuthenticator.ServeHTTPCall.Receives.Request).To(BeNil())
			Expect(keepGoing).To(BeFalse())

			Expect(writer.Code).To(Equal(http.StatusUnauthorized))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": [
					"Authorization header is invalid: missing"
				]
			}`))
		})
	})

	Context("when the token is malformed", func() {
		It("returns a 401 status code and error message", func() {
			request.Header.Set("Authorization", "bearer some-junk")

			validator.ParseCall.Returns.Error = errors.New("broken token")

			keepGoing := auth.ServeHTTP(writer, request, context)
			Expect(userAuthenticator.ServeHTTPCall.Receives.Request).To(BeNil())
			Expect(clientAuthenticator.ServeHTTPCall.Receives.Request).To(BeNil())
			Expect(keepGoing).To(BeFalse())

			Expect(writer.Code).To(Equal(http.StatusUnauthorized))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": [
					"Authorization header is invalid: corrupt"
				]
			}`))
		})
	})
})
