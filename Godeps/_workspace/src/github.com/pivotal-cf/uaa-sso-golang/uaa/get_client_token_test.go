package uaa_test

import (
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("GetClientToken", func() {
    var fakeUAAServer *httptest.Server
    var auth uaa.UAA

    Context("when UAA is responding normally", func() {
        BeforeEach(func() {
            fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                if req.URL.Path == "/oauth/token" && req.Method == "POST" && strings.Contains(req.Header.Get("Authorization"), "Basic") {
                    err := req.ParseForm()
                    if err != nil {
                        panic(err)
                    }

                    if req.Form.Get("grant_type") != "client_credentials" {
                        w.WriteHeader(http.StatusNotAcceptable)
                        w.Write([]byte(`{"error":"unacceptable"}`))
                        return
                    }

                    response := `{
                            "access_token": "client-access-token",
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

        It("returns the client auth token", func() {
            token, err := uaa.GetClientToken(auth)
            Expect(err).To(BeNil())
            Expect(token.Access).To(Equal("client-access-token"))
        })
    })

    Context("when UAA is not responding normally", func() {
        BeforeEach(func() {
            fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                if req.URL.Path == "/oauth/token" && req.Method == "POST" && strings.Contains(req.Header.Get("Authorization"), "Basic") {
                    response := `{"errors": "Out to lunch"}`

                    w.WriteHeader(http.StatusGone)
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
            _, err := uaa.GetClientToken(auth)
            Expect(err).To(BeAssignableToTypeOf(uaa.Failure{}))
            Expect(err.Error()).To(Equal(`UAA Failure: 410 {"errors": "Out to lunch"}`))
        })
    })
})
