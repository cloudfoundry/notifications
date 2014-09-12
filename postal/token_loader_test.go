package postal_test

import (
    "errors"
    "net/http"
    "net/url"

    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/test_helpers/fakes"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("TokenLoader", func() {
    var tokenLoader postal.TokenLoader
    var fakeUAA fakes.FakeUAAClient

    BeforeEach(func() {
        fakeUAA = fakes.FakeUAAClient{
            ClientToken: uaa.Token{
                Access: "the-client-token",
            },
        }

        tokenLoader = postal.NewTokenLoader(&fakeUAA)
    })

    Describe("Load", func() {
        It("returns the client token from UAA", func() {
            token, err := tokenLoader.Load()
            if err != nil {
                panic(err)
            }

            Expect(token).To(Equal("the-client-token"))
        })

        It("assigns the access token on the uaa client", func() {
            _, err := tokenLoader.Load()
            if err != nil {
                panic(err)
            }

            Expect(fakeUAA.AccessToken).To(Equal("the-client-token"))
        })

        Context("error handling", func() {
            It("identifies UAA being down, returning an error", func() {
                fakeUAA.ClientTokenError = uaa.NewFailure(http.StatusNotFound, []byte("404 Not Found: Requested route ('uaa.10.244.0.34.xip.io') does not exist."))

                _, err := tokenLoader.Load()

                Expect(err).To(BeAssignableToTypeOf(postal.UAADownError("")))
                Expect(err.Error()).To(Equal("UAA is unavailable"))
            })

            It("returns a generic error when UAA returns a 404 that does not indicate that it is down", func() {
                fakeUAA.ClientTokenError = uaa.NewFailure(http.StatusNotFound, []byte("Not found"))

                _, err := tokenLoader.Load()

                Expect(err).To(BeAssignableToTypeOf(postal.UAAGenericError("")))
                Expect(err.Error()).To(Equal("UAA Unknown 404 error message: Not found"))
            })

            It("handles non-404 UAAFailure errors", func() {
                failure := uaa.NewFailure(http.StatusInternalServerError, []byte("Banana!"))
                fakeUAA.ClientTokenError = failure

                _, err := tokenLoader.Load()

                Expect(err).To(BeAssignableToTypeOf(postal.UAADownError("")))
                Expect(err.Error()).To(Equal(failure.Message()))
            })

            It("returns an error when it cannot make a connection to UAA", func() {
                fakeUAA.ClientTokenError = &url.Error{}

                _, err := tokenLoader.Load()

                Expect(err).To(BeAssignableToTypeOf(postal.UAADownError("")))
                Expect(err.Error()).To(Equal("UAA is unavailable"))
            })

            It("handles all other error cases", func() {
                fakeUAA.ClientTokenError = errors.New("BOOM!")

                _, err := tokenLoader.Load()

                Expect(err).To(BeAssignableToTypeOf(postal.UAAGenericError("")))
                Expect(err.Error()).To(Equal("UAA Unknown Error: BOOM!"))
            })
        })
    })
})
