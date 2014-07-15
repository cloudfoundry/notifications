package uaa_test

import (
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("GetTokenKey", func() {
    var fakeUAAServer *httptest.Server
    var auth uaa.UAA

    Context("when UAA is responding normally", func() {
        BeforeEach(func() {
            fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                if req.URL.Path == "/token_key" && req.Method == "GET" && strings.Contains(req.Header.Get("Authorization"), "Bearer") {
                    response := `{
                      "alg": "SHA256withRSA",
                      "value": "THIS-IS-THE-PUBLIC-KEY"
                    }`
                    w.WriteHeader(http.StatusOK)
                    w.Write([]byte(response))
                } else if req.URL.Path == "/oauth/token" && req.Method == "POST" && strings.Contains(req.Header.Get("Authorization"), "Basic") {
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
            auth = uaa.NewUAA("", fakeUAAServer.URL, "the-client-id", "the-client-secret", "")
        })

        It("returns the public key that UAA tokens can be validated with", func() {
            key, err := uaa.GetTokenKey(auth)
            if err != nil {
                panic(err)
            }

            Expect(key).To(Equal("THIS-IS-THE-PUBLIC-KEY"))
        })
    })

    Context("when UAA is not responding normally", func() {
        BeforeEach(func() {
            fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                if req.URL.Path == "/token_key" && req.Method == "GET" && strings.Contains(req.Header.Get("Authorization"), "Bearer") {
                    response := `{"errors": "Out to lunch"}`

                    w.WriteHeader(http.StatusGone)
                    w.Write([]byte(response))
                } else if req.URL.Path == "/oauth/token" && req.Method == "POST" && strings.Contains(req.Header.Get("Authorization"), "Basic") {
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
            auth = uaa.NewUAA("", fakeUAAServer.URL, "the-client-id", "the-client-secret", "")
        })

        AfterEach(func() {
            fakeUAAServer.Close()
        })

        It("returns an error message", func() {
            _, err := uaa.GetTokenKey(auth)
            Expect(err).To(BeAssignableToTypeOf(uaa.Failure{}))
            Expect(err.Error()).To(Equal(`UAA Failure: 410 {"errors": "Out to lunch"}`))
        })
    })
})
