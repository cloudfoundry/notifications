package uaa_test

import (
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Refresh", func() {
    var auth uaa.UAA
    var fakeUAAServer *httptest.Server

    BeforeEach(func() {
        auth = uaa.NewUAA("http://login.example.com", "http://uaa.example.com", "the-client-id", "the-client-secret", "")
    })

    Context("when UAA is responding normally", func() {
        BeforeEach(func() {
            fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                if req.URL.Path == "/oauth/token" && req.Method == "POST" && strings.Contains(req.Header.Get("Authorization"), "Basic") {
                    err := req.ParseForm()
                    if err != nil {
                        panic(err)
                    }

                    if req.Form.Get("refresh_token") != "refresh-token" {
                        w.WriteHeader(http.StatusUnauthorized)
                        w.Write([]byte(`{"error":"invalid_token"}`))
                        return
                    }

                    response := `{
                            "access_token": "access-token",
                            "refresh_token": "refresh-token",
                            "token_type": "bearer"
                        }`

                    w.WriteHeader(http.StatusOK)
                    w.Write([]byte(response))
                } else {
                    w.WriteHeader(http.StatusNotFound)
                }
            }))
            auth = uaa.NewUAA("http://login.example.com", fakeUAAServer.URL, "the-client-id", "the-client-secret", "")
        })

        AfterEach(func() {
            fakeUAAServer.Close()
        })

        It("returns a token received in exchange for a refresh token", func() {
            token, err := uaa.Refresh(auth, "refresh-token")
            Expect(err).To(BeNil())
            Expect(token.Access).To(Equal("access-token"))
        })

        It("returns an invalid refresh token error for invalid token", func() {
            _, err := uaa.Refresh(auth, "bad-refresh-token")
            Expect(err).To(Equal(uaa.InvalidRefreshToken))
        })
    })

    Context("when UAA is not responding normally", func() {
        BeforeEach(func() {
            fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                if req.URL.Path == "/oauth/token" && req.Method == "POST" && strings.Contains(req.Header.Get("Authorization"), "Basic") {
                    response := `{"errors": "client_error"}`

                    w.WriteHeader(http.StatusMethodNotAllowed)
                    w.Write([]byte(response))
                } else {
                    w.WriteHeader(http.StatusNotFound)
                }
            }))
            auth = uaa.NewUAA("http://login.example.com", fakeUAAServer.URL, "the-client-id", "the-client-secret", "")
        })

        AfterEach(func() {
            fakeUAAServer.Close()
        })

        It("returns an error message", func() {
            _, err := uaa.Refresh(auth, "refresh-token")
            Expect(err).To(BeAssignableToTypeOf(uaa.Failure{}))
            Expect(err.Error()).To(Equal(`UAA Failure: 405 {"errors": "client_error"}`))
        })
    })
})
