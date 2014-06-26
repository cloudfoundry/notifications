package uaa_test

import (
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Exchange", func() {
    var fakeUAAServer *httptest.Server
    var auth uaa.UAA

    Context("when UAA is responding normally", func() {
        BeforeEach(func() {
            fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                if req.URL.Path == "/oauth/token" && req.Method == "POST" && strings.Contains(req.Header.Get("Authorization"), "Basic") {
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

        It("returns a token received in exchange for a code from UAA", func() {
            token, err := uaa.Exchange(auth, "1234")
            if err != nil {
                panic(err)
            }

            Expect(token).To(Equal(uaa.Token{
                Access:  "access-token",
                Refresh: "refresh-token",
            }))
        })
    })

    Context("when UAA is not responding normally", func() {
        BeforeEach(func() {
            fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                if req.URL.Path == "/oauth/token" && req.Method == "POST" && strings.Contains(req.Header.Get("Authorization"), "Basic") {
                    w.WriteHeader(http.StatusUnauthorized)
                    w.Write([]byte(`{"errors": "Unauthorized"}`))
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
            _, err := uaa.Exchange(auth, "1234")
            Expect(err).To(BeAssignableToTypeOf(uaa.Failure{}))
            Expect(err.Error()).To(Equal(`UAA Failure: 401 {"errors": "Unauthorized"}`))
        })
    })
})
