package uaa_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/uaa"

	"github.com/dgrijalva/jwt-go"
	"github.com/pivotal-cf-experimental/warrant"
	"github.com/pivotal-golang/lager"
)

var _ = Describe("TokenValidator", func() {
	var (
		validator  *uaa.TokenValidator
		rawToken   string
		token      *jwt.Token
		err        error
		keyFetcher *mocks.KeyFetcher
	)

	BeforeEach(func() {
		keyFetcher = &mocks.KeyFetcher{}
		keyFetcher.GetSigningKeysCall.Returns.Keys = []warrant.SigningKey{
			{
				KeyId:     "some-key",
				Algorithm: "RS256",
				Value:     helpers.UAAPublicKey,
			},
		}
		validator = uaa.NewTokenValidator(lager.NewLogger("test"), keyFetcher)
	})

	Describe("loading signing keys", func() {
		It("returns an error when loading keys fails", func() {
			keyFetcher.GetSigningKeysCall.Returns.Error = errors.New("network failure")
			err := validator.LoadSigningKeys()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("parsing tokens", func() {
		JustBeforeEach(func() {

			token, err = validator.Parse(rawToken)
		})

		Context("when the request contains a valid auth token", func() {
			BeforeEach(func() {
				tokenHeader := map[string]interface{}{
					"alg": "RS256",
					"kid": "some-key",
				}
				tokenClaims := map[string]interface{}{
					"jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
					"client_id": "mister-client",
					"cid":       "mister-client",
					"exp":       3404281214,
					"scope":     []string{"gaben.scope"},
				}
				rawToken = helpers.BuildToken(tokenHeader, tokenClaims)
			})

			It("returns no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns the token", func() {
				Expect(token.Claims).To(Equal(map[string]interface{}{
					"jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
					"client_id": "mister-client",
					"cid":       "mister-client",
					"exp":       float64(3404281214),
					"scope":     []interface{}{"gaben.scope"},
				}))
			})
		})

		Context("when the request uses an expired auth token", func() {
			BeforeEach(func() {
				tokenHeader := map[string]interface{}{
					"alg": "RS256",
					"kid": "some-key",
				}
				tokenClaims := map[string]interface{}{
					"jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
					"client_id": "mister-client",
					"cid":       "mister-client",
					"exp":       1404281214,
				}
				rawToken = helpers.BuildToken(tokenHeader, tokenClaims)
			})

			It("returns an error", func() {
				Expect(err.Error()).To(ContainSubstring("expired"))
			})
		})

		Context("with a token signed using the public key (used symmetrically)", func() {
			BeforeEach(func() {
				tokenHeader := map[string]interface{}{
					"alg": "HS256",
					"kid": "some-key",
				}
				tokenClaims := map[string]interface{}{
					"jti":       "c5f6a266-5cf0-4ae2-9647-2615e7d28fa1",
					"client_id": "mister-client",
					"cid":       "mister-client",
					"exp":       3404281214,
					"scope":     []string{"gaben.scope"},
				}

				rawToken = helpers.BuildTokenWithKey(tokenHeader, tokenClaims, helpers.UAAPublicKey)
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
