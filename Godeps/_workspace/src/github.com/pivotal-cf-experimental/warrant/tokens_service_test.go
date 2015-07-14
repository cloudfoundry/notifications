package warrant_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TokensService", func() {
	var (
		service warrant.TokensService
		config  warrant.Config
	)

	BeforeEach(func() {
		config = warrant.Config{
			Host:          fakeUAA.URL(),
			SkipVerifySSL: true,
			TraceWriter:   TraceWriter,
		}
		service = warrant.NewTokensService(config)
	})

	Describe("Decode", func() {
		It("returns a decoded token given an encoded token string", func() {
			encodedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlci1pZCIsInNjb3BlIjpbInNjaW0ucmVhZCIsImNsb3VkX2NvbnRyb2xsZXIuYWRtaW4iLCJwYXNzd29yZC53cml0ZSJdfQ.QWNTRFahfCn7qSWxEHTCn6QeZMJxNMq9a_TP8aANc4k"
			token, err := service.Decode(encodedToken)
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal(warrant.Token{
				UserID: "user-id",
				Scopes: []string{
					"scim.read",
					"cloud_controller.admin",
					"password.write",
				},
			}))
		})

		Context("failure cases", func() {
			It("returns an error if there is an invalid number of segments", func() {
				encodedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
				_, err := service.Decode(encodedToken)
				Expect(err).To(BeAssignableToTypeOf(warrant.InvalidTokenError{}))
				Expect(err).To(MatchError("invalid number of segments in token (1/3)"))
			})

			It("returns an error if the claims segment cannot be decoded", func() {
				encodedToken := "headers.invalid#claims.signature"
				_, err := service.Decode(encodedToken)
				Expect(err).To(BeAssignableToTypeOf(warrant.InvalidTokenError{}))
				Expect(err).To(MatchError("claims cannot be decoded: illegal base64 data at input byte 7"))
			})

			It("returns an error if the token cannot be unmarshaled", func() {
				encodedToken := "headers.invalid-claims.signature"
				_, err := service.Decode(encodedToken)
				Expect(err).To(BeAssignableToTypeOf(warrant.InvalidTokenError{}))
				Expect(err).To(MatchError("token cannot be parsed: invalid character '\\u008a' looking for beginning of value"))
			})
		})
	})

	Describe("GetSigningKey", func() {
		It("returns the public key, used to sign tokens", func() {
			key, err := service.GetSigningKey()
			Expect(err).NotTo(HaveOccurred())
			Expect(key).To(Equal(warrant.SigningKey{
				Algorithm: "SHA256withRSA",
				Value:     fakeUAAPublicKey,
			}))
		})

		Context("failure cases", func() {
			It("returns an error if the HTTP request fails", func() {
				erroringServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))

				service = warrant.NewTokensService(warrant.Config{
					Host:          erroringServer.URL,
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				_, err := service.GetSigningKey()
				Expect(err).To(BeAssignableToTypeOf(warrant.UnexpectedStatusError{}))
			})

			It("returns an error if the response JSON cannot be parsed", func() {
				malformedJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write([]byte("this is not JSON"))
				}))

				service = warrant.NewTokensService(warrant.Config{
					Host:          malformedJSONServer.URL,
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				_, err := service.GetSigningKey()
				Expect(err).To(BeAssignableToTypeOf(warrant.MalformedResponseError{}))
				Expect(err).To(MatchError("malformed response: invalid character 'h' in literal true (expecting 'r')"))
			})
		})
	})
})
