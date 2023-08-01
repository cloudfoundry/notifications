package middleware_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/web/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Authenticator", func() {
	var (
		ware    middleware.Authenticator
		request *http.Request
		writer  *httptest.ResponseRecorder
		context stack.Context

		validator     *mocks.TokenValidator
		expectedToken *jwt.Token
	)

	BeforeEach(func() {
		var err error

		validator = &mocks.TokenValidator{}
		ware = middleware.NewAuthenticator(validator, "fake.scope", "gaben.scope")
		writer = httptest.NewRecorder()
		request, err = http.NewRequest("GET", "/some/path", nil)
		if err != nil {
			panic(err)
		}
		context = stack.NewContext()
	})

	Context("when the request contains a valid auth token", func() {
		BeforeEach(func() {
			expectedToken = jwt.New(jwt.SigningMethodRS256)
			claims := jwt.MapClaims{}
			claims["jti"] = "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1"
			claims["client_id"] = "mister-client"
			claims["cid"] = "mister-client"
			claims["exp"] = 3404281214
			claims["scope"] = []interface{}{"gaben.scope"}
			expectedToken.Claims = claims

			validator.ParseCall.Returns.Token = expectedToken

			requestBody, err := json.Marshal(map[string]string{
				"kind": "forgot_password",
				"text": "Try to remember your password next time",
			})
			if err != nil {
				panic(err)
			}

			request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
			if err != nil {
				panic(err)
			}
			request.Header.Set("Authorization", "Bearer valid-token")
		})

		It("validates the token", func() {
			ware.ServeHTTP(writer, request, context)
			Expect(validator.ParseCall.Receives.Token).To(Equal("valid-token"))
		})

		It("allows the request through", func() {
			returnValue := ware.ServeHTTP(writer, request, context)

			Expect(returnValue).To(BeTrue())
			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(len(writer.Body.Bytes())).To(Equal(0))
		})

		It("sets the token on the context", func() {
			ware.ServeHTTP(writer, request, context)

			contextToken := context.Get("token")
			Expect(contextToken).To(Equal(expectedToken))
		})

		It("sets the client_id on the context", func() {
			ware.ServeHTTP(writer, request, context)

			Expect(context.Get("client_id")).To(Equal("mister-client"))
		})

		Context("when the prefix to the token has different capitalization", func() {
			It("still sets the token", func() {
				request.Header.Set("Authorization", "bearer some-token")
				ware.ServeHTTP(writer, request, context)

				contextToken := context.Get("token")
				Expect(contextToken).To(Equal(expectedToken))
			})
		})
	})

	Context("when the request does not contain an auth token", func() {
		BeforeEach(func() {
			requestBody, err := json.Marshal(map[string]string{
				"kind": "forgot_password",
				"text": "Try to remember your password next time",
			})
			if err != nil {
				panic(err)
			}

			request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
			if err != nil {
				panic(err)
			}
		})

		It("returns a 401 status code and error message", func() {
			returnValue := ware.ServeHTTP(writer, request, context)

			Expect(returnValue).To(BeFalse())
			Expect(writer.Code).To(Equal(http.StatusUnauthorized))

			parsed := map[string][]string{}
			err := json.Unmarshal(writer.Body.Bytes(), &parsed)
			if err != nil {
				panic(err)
			}

			Expect(parsed["errors"]).To(ContainElement("Authorization header is invalid: missing"))
		})
	})

	Context("when the auth token does not contain the correct scope", func() {
		BeforeEach(func() {
			expectedToken = jwt.New(jwt.SigningMethodRS256)
			claims := jwt.MapClaims{}
			claims["jti"] = "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1"
			claims["client_id"] = "mister-client"
			claims["cid"] = "mister-client"
			claims["exp"] = 3404281214
			claims["scope"] = []interface{}{"cloud_controller.admin"}
			expectedToken.Claims = claims

			validator.ParseCall.Returns.Token = expectedToken

			requestBody, err := json.Marshal(map[string]string{
				"kind": "forgot_password",
				"text": "Try to remember your password next time",
			})
			if err != nil {
				panic(err)
			}

			request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
			if err != nil {
				panic(err)
			}
			request.Header.Set("Authorization", "Bearer bad-scope-token")
		})

		It("returns a 403 status code and error message", func() {
			returnValue := ware.ServeHTTP(writer, request, context)

			Expect(returnValue).To(BeFalse())
			Expect(writer.Code).To(Equal(http.StatusForbidden))

			parsed := map[string][]string{}
			err := json.Unmarshal(writer.Body.Bytes(), &parsed)
			if err != nil {
				panic(err)
			}

			Expect(parsed["errors"]).To(ContainElement("You are not authorized to perform the requested action"))
		})
	})

	Context("when the auth token does not contain any scopes", func() {
		BeforeEach(func() {
			expectedToken = jwt.New(jwt.SigningMethodRS256)
			claims := jwt.MapClaims{}
			claims["jti"] = "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1"
			claims["client_id"] = "mister-client"
			claims["cid"] = "mister-client"
			claims["exp"] = 3404281214
			expectedToken.Claims = claims

			validator.ParseCall.Returns.Token = expectedToken

			requestBody, err := json.Marshal(map[string]string{
				"kind": "forgot_password",
				"text": "Try to remember your password next time",
			})
			if err != nil {
				panic(err)
			}

			request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
			if err != nil {
				panic(err)
			}
			request.Header.Set("Authorization", "Bearer missing-scope-token")
		})

		It("returns a 403 status code and error message", func() {
			returnValue := ware.ServeHTTP(writer, request, context)

			Expect(returnValue).To(BeFalse())
			Expect(writer.Code).To(Equal(http.StatusForbidden))

			parsed := map[string][]string{}
			err := json.Unmarshal(writer.Body.Bytes(), &parsed)
			if err != nil {
				panic(err)
			}

			Expect(parsed["errors"]).To(ContainElement("You are not authorized to perform the requested action"))
		})
	})

	Context("when the request does not contain a auth valid token", func() {
		BeforeEach(func() {
			requestBody, err := json.Marshal(map[string]string{
				"kind": "forgot_password",
				"text": "Try to remember your password next time",
			})
			if err != nil {
				panic(err)
			}

			request, err = http.NewRequest("POST", "/users/user-123", bytes.NewReader(requestBody))
			if err != nil {
				panic(err)
			}
			request.Header.Set("Authorization", "Bearer something-invalid")
			validator.ParseCall.Returns.Error = errors.New("bad token")
		})

		It("returns a 401 status code and error message", func() {
			returnValue := ware.ServeHTTP(writer, request, context)

			Expect(returnValue).To(BeFalse())
			Expect(writer.Code).To(Equal(http.StatusUnauthorized))

			parsed := map[string][]string{}
			err := json.Unmarshal(writer.Body.Bytes(), &parsed)
			if err != nil {
				panic(err)
			}

			Expect(parsed["errors"][0]).To(ContainSubstring("Authorization header is invalid"))
		})
	})
})
