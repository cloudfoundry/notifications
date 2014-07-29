package uaa_test

import (
    "reflect"

    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("UAA", func() {
    var auth uaa.UAA

    BeforeEach(func() {
        auth = uaa.NewUAA("http://login.example.com", "http://uaa.example.com", "the-client-id", "the-client-secret", "")
    })

    Describe("NewUAA", func() {
        It("defaults VerifySSL to true", func() {
            Expect(auth.VerifySSL).To(BeTrue())
        })
    })

    Describe("AuthorizeURL", func() {
        It("returns the URL for the /oauth/authorize endpoint", func() {
            Expect(auth.AuthorizeURL()).To(Equal("http://login.example.com/oauth/authorize"))
        })
    })

    Describe("LoginURL", func() {
        It("returns a url to be used as the redirect for authenticating with UAA", func() {
            auth.ClientID = "fake-client"
            auth.RedirectURL = "http://redirect.example.com"
            auth.Scope = "username,email"
            auth.State = "some-data"
            auth.AccessType = "offline"
            auth.ApprovalPrompt = "yes"
            expected := "http://login.example.com/oauth/authorize?access_type=offline&approval_prompt=yes&client_id=fake-client&redirect_uri=http%3A%2F%2Fredirect.example.com&response_type=code&scope=username%2Cemail&state=some-data"
            Expect(auth.LoginURL()).To(Equal(expected))
        })
    })

    Describe("SetToken", func() {
        It("assigns the given token value to the AccessToken field", func() {
            Expect(auth.AccessToken).To(Equal(""))

            auth.SetToken("the-new-access-token")

            Expect(auth.AccessToken).To(Equal("the-new-access-token"))
        })
    })

    Describe("Exchange", func() {
        var exchangeWasCalledWith string

        It("delegates to the Exchange Command", func() {
            Expect(reflect.ValueOf(auth.ExchangeCommand).Pointer()).To(Equal(reflect.ValueOf(uaa.Exchange).Pointer()))

            auth.ExchangeCommand = func(u uaa.UAA, authCode string) (uaa.Token, error) {
                exchangeWasCalledWith = authCode
                return uaa.Token{}, nil
            }

            auth.Exchange("auth-code")

            Expect(exchangeWasCalledWith).To(Equal("auth-code"))
        })
    })

    Describe("Refresh", func() {
        var refreshWasCalledWith string

        It("delegates to the Refresh Command", func() {
            Expect(reflect.ValueOf(auth.RefreshCommand).Pointer()).To(Equal(reflect.ValueOf(uaa.Refresh).Pointer()))

            auth.RefreshCommand = func(u uaa.UAA, refreshToken string) (uaa.Token, error) {
                refreshWasCalledWith = refreshToken
                return uaa.Token{}, nil
            }

            auth.Refresh("some-token")

            Expect(refreshWasCalledWith).To(Equal("some-token"))
        })
    })

    Describe("GetClientToken", func() {
        var getClientTokenWasCalled bool

        It("delegates to the GetClientToken Command", func() {
            Expect(reflect.ValueOf(auth.GetClientTokenCommand).Pointer()).To(Equal(reflect.ValueOf(uaa.GetClientToken).Pointer()))

            auth.GetClientTokenCommand = func(u uaa.UAA) (uaa.Token, error) {
                getClientTokenWasCalled = true
                return uaa.Token{}, nil
            }

            auth.GetClientToken()

            Expect(getClientTokenWasCalled).To(Equal(true))
        })
    })

    Describe("UserByID", func() {
        var userByIDWasCalledWith string

        It("delegates to the UserByID Command", func() {
            Expect(reflect.ValueOf(auth.UserByIDCommand).Pointer()).To(Equal(reflect.ValueOf(uaa.UserByID).Pointer()))

            auth.UserByIDCommand = func(u uaa.UAA, id string) (uaa.User, error) {
                userByIDWasCalledWith = id
                return uaa.User{}, nil
            }

            auth.UserByID("my-special-id")

            Expect(userByIDWasCalledWith).To(Equal("my-special-id"))
        })
    })

    Describe("GetTokenKey", func() {
        var getTokenKeyWasCalled bool

        It("delegates to the GetTokenKey Command", func() {
            Expect(reflect.ValueOf(auth.GetTokenKeyCommand).Pointer()).To(Equal(reflect.ValueOf(uaa.GetTokenKey).Pointer()))

            auth.GetTokenKeyCommand = func(u uaa.UAA) (string, error) {
                getTokenKeyWasCalled = true
                return "", nil
            }

            auth.GetTokenKey()

            Expect(getTokenKeyWasCalled).To(BeTrue())
        })
    })

    Describe("UsersByIDs", func() {
        var usersByIDsWasCalledWith []string

        It("delegates to the UsersByIDs command", func() {
            Expect(reflect.ValueOf(auth.UsersByIDsCommand).Pointer()).To(Equal(reflect.ValueOf(uaa.UsersByIDs).Pointer()))

            auth.UsersByIDsCommand = func(u uaa.UAA, ids ...string) ([]uaa.User, error) {
                usersByIDsWasCalledWith = ids
                return []uaa.User{}, nil
            }

            auth.UsersByIDs([]string{"something", "another-thing"}...)

            Expect(usersByIDsWasCalledWith).To(Equal([]string{"something", "another-thing"}))
        })
    })
})
