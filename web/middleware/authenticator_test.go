package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Authenticator", func() {
	var ware middleware.Authenticator
	var request *http.Request
	var writer *httptest.ResponseRecorder
	var rawToken string
	var context stack.Context

	BeforeEach(func() {
		var err error

		ware = middleware.NewAuthenticator(application.UAAPublicKey, "fake.scope", "gaben.scope")
		writer = httptest.NewRecorder()
		request, err = http.NewRequest("GET", "/some/path", nil)
		if err != nil {
			panic(err)
		}
		context = stack.NewContext()
	})

	Context("when the request contains a valid auth token", func() {
		BeforeEach(func() {
			tokenHeader := map[string]interface{}{
				"alg": "FAST",
			}
			tokenClaims := map[string]interface{}{
				"jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
				"client_id": "mister-client",
				"cid":       "mister-client",
				"exp":       3404281214,
				"scope":     []string{"gaben.scope"},
			}
			rawToken = fakes.BuildToken(tokenHeader, tokenClaims)

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
			request.Header.Set("Authorization", "Bearer "+rawToken)
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
			Expect(contextToken).NotTo(BeNil())

			token, err := jwt.Parse(rawToken, func(*jwt.Token) (interface{}, error) {
				return []byte(application.UAAPublicKey), nil
			})
			if err != nil {
				panic(err)
			}

			Expect(*(contextToken.(*jwt.Token))).To(Equal(*token))
		})

		Context("when the prefix to the token has different capitalization", func() {
			It("still sets the token", func() {
				request.Header.Set("Authorization", "bearer "+rawToken)
				ware.ServeHTTP(writer, request, context)

				contextToken := context.Get("token")
				Expect(contextToken).NotTo(BeNil())

				token, err := jwt.Parse(rawToken, func(*jwt.Token) (interface{}, error) {
					return []byte(application.UAAPublicKey), nil
				})
				if err != nil {
					panic(err)
				}

				Expect(*(contextToken.(*jwt.Token))).To(Equal(*token))
			})
		})
	})

	Context("when the request uses an expired auth token", func() {
		BeforeEach(func() {
			tokenHeader := map[string]interface{}{
				"alg": "FAST",
			}
			tokenClaims := map[string]interface{}{
				"jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
				"client_id": "mister-client",
				"cid":       "mister-client",
				"exp":       1404281214,
			}
			rawToken = fakes.BuildToken(tokenHeader, tokenClaims)

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
			request.Header.Set("Authorization", "Bearer "+rawToken)
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

			Expect(parsed["errors"]).To(ContainElement("Authorization header is invalid: expired"))
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
			tokenHeader := map[string]interface{}{
				"alg": "FAST",
			}
			tokenClaims := map[string]interface{}{
				"jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
				"client_id": "mister-client",
				"cid":       "mister-client",
				"exp":       3404281214,
				"scope":     []string{"cloud_controller.admin"},
			}
			rawToken = fakes.BuildToken(tokenHeader, tokenClaims)

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
			request.Header.Set("Authorization", "Bearer "+rawToken)
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
			tokenHeader := map[string]interface{}{
				"alg": "FAST",
			}
			tokenClaims := map[string]interface{}{
				"jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
				"client_id": "mister-client",
				"cid":       "mister-client",
				"exp":       3404281214,
			}
			rawToken = fakes.BuildToken(tokenHeader, tokenClaims)

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
			request.Header.Set("Authorization", "Bearer "+rawToken)
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

			Expect(parsed["errors"]).To(ContainElement("Authorization header is invalid: corrupt"))
		})
	})
})
