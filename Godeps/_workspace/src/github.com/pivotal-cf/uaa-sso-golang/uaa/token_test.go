package uaa_test

import (
    "github.com/dgrijalva/jwt-go"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Token", func() {
    Describe("IsPresent", func() {
        It("returns true if the entire token is populated", func() {
            token := uaa.NewToken()
            Expect(token.IsPresent()).To(BeFalse())

            token.Access = "access-token"
            Expect(token.IsPresent()).To(BeFalse())

            token.Refresh = "refresh-token"
            Expect(token.IsPresent()).To(BeTrue())
        })
    })

    Describe("Type", func() {
        It("returns bearer", func() {
            token := uaa.NewToken()
            Expect(token.Type()).To(Equal("bearer"))
        })
    })

    Describe("IsExpired", func() {
        var token uaa.Token

        Context("when the 'exp' key is in the future", func() {
            It("returns false", func() {
                header := jwt.EncodeSegment([]byte(`{"alg":"RS256"}`))
                body := jwt.EncodeSegment([]byte(`{"exp":32503683661}`))
                token.Access = header + "." + body
                expired, err := token.IsExpired()
                Expect(err).To(BeNil())
                Expect(expired).To(BeFalse())
            })
        })

        Context("when the 'exp' key is in the past", func() {
            It("returns true", func() {
                header := jwt.EncodeSegment([]byte(`{"alg":"RS256"}`))
                body := jwt.EncodeSegment([]byte(`{"exp":915152461}`))
                token.Access = header + "." + body
                expired, err := token.IsExpired()
                Expect(err).To(BeNil())
                Expect(expired).To(BeTrue())
            })
        })

        Context("handling errors", func() {
            It("returns a TokenDecodeError when the token cannot be decoded", func() {
                header := jwt.EncodeSegment([]byte(`{"alg":"RS256"}`))
                body := "bad!token!body"
                token.Access = header + "." + body
                _, err := token.IsExpired()
                Expect(err).To(Equal(uaa.TokenDecodeError))
            })

            It("returns a JSONParseError when the json cannot be parsed", func() {
                header := jwt.EncodeSegment([]byte(`{"alg":"RS256"}`))
                body := jwt.EncodeSegment([]byte("bad-json"))
                token.Access = header + "." + body
                _, err := token.IsExpired()
                Expect(err).To(Equal(uaa.JSONParseError))
            })
        })
    })
})
