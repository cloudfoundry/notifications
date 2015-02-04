package postal_test

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TokenLoader", func() {
	var tokenLoader postal.TokenLoader
	var uaaClient *fakes.UAAClient
	var clientToken string
	var tokenHeader map[string]interface{}
	var tokenClaims map[string]interface{}

	Describe("Load", func() {
		BeforeEach(func() {
			tokenHeader = map[string]interface{}{
				"alg": "FAST",
			}

			tokenClaims = map[string]interface{}{
				"client_id": "mister-client",
				"exp":       time.Now().Add(5 * time.Minute).Unix(),
				"scope":     []string{"notifications.write"},
			}

			clientToken = fakes.BuildToken(tokenHeader, tokenClaims)

			uaaClient = fakes.NewUAAClient()
			uaaClient.ClientToken = uaa.Token{Access: clientToken}

			tokenLoader = postal.NewTokenLoader(uaaClient)
		})

		AfterEach(func() {
			postal.ResetLoader()
		})

		It("returns the client token from UAA", func() {
			token, err := tokenLoader.Load()
			if err != nil {
				panic(err)
			}

			Expect(token).To(Equal(clientToken))
		})

		It("assigns the access token on the uaa client", func() {
			_, err := tokenLoader.Load()
			if err != nil {
				panic(err)
			}

			Expect(uaaClient.AccessToken).To(Equal(clientToken))
		})

		Context("When the current client token is not expired", func() {
			It("Does not ask UAA for a new token", func() {
				token, err := tokenLoader.Load()
				if err != nil {
					panic(err)
				}
				Expect(token).To(Equal(clientToken))

				uaaClient.ClientToken.Access = "the-wrong-token"

				token, err = tokenLoader.Load()
				if err != nil {
					panic(err)
				}

				Expect(token).To(Equal(clientToken))
			})
		})

		Context("When the current client token is expired", func() {
			It("Does ask UAA for a new token", func() {
				tokenClaims["exp"] = time.Now().Add(-5 * time.Minute).Unix()
				expiredToken := fakes.BuildToken(tokenHeader, tokenClaims)
				uaaClient.ClientToken.Access = expiredToken

				token, err := tokenLoader.Load()
				if err != nil {
					panic(err)
				}
				Expect(token).To(Equal(expiredToken))

				uaaClient.ClientToken.Access = "the-correct-token"

				token, err = tokenLoader.Load()
				if err != nil {
					panic(err)
				}

				Expect(token).To(Equal("the-correct-token"))

			})
		})

		Context("error handling", func() {
			It("identifies UAA being down, returning an error", func() {
				uaaClient.ClientTokenError = uaa.NewFailure(http.StatusNotFound, []byte("404 Not Found: Requested route ('uaa.10.244.0.34.xip.io') does not exist."))

				_, err := tokenLoader.Load()

				Expect(err).To(BeAssignableToTypeOf(postal.UAADownError("")))
				Expect(err.Error()).To(Equal("UAA is unavailable"))
			})

			It("returns a generic error when UAA returns a 404 that does not indicate that it is down", func() {
				uaaClient.ClientTokenError = uaa.NewFailure(http.StatusNotFound, []byte("Not found"))

				_, err := tokenLoader.Load()

				Expect(err).To(BeAssignableToTypeOf(postal.UAAGenericError("")))
				Expect(err.Error()).To(Equal("UAA Unknown 404 error message: Not found"))
			})

			It("handles non-404 UAAFailure errors", func() {
				failure := uaa.NewFailure(http.StatusInternalServerError, []byte("Banana!"))
				uaaClient.ClientTokenError = failure

				_, err := tokenLoader.Load()

				Expect(err).To(BeAssignableToTypeOf(postal.UAADownError("")))
				Expect(err.Error()).To(Equal(failure.Message()))
			})

			It("returns an error when it cannot make a connection to UAA", func() {
				uaaClient.ClientTokenError = &url.Error{}

				_, err := tokenLoader.Load()

				Expect(err).To(BeAssignableToTypeOf(postal.UAADownError("")))
				Expect(err.Error()).To(Equal("UAA is unavailable"))
			})

			It("handles all other error cases", func() {
				uaaClient.ClientTokenError = errors.New("BOOM!")

				_, err := tokenLoader.Load()

				Expect(err).To(BeAssignableToTypeOf(postal.UAAGenericError("")))
				Expect(err.Error()).To(Equal("UAA Unknown Error: BOOM!"))
			})
		})
	})
})
