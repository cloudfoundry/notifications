package uaa_test

import (
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("UserByID", func() {
    var fakeUAAServer *httptest.Server
    var auth uaa.UAA

    Context("when UAA is responding normally", func() {
        BeforeEach(func() {
            fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                if req.URL.Path == "/Users/87dfc5b4-daf9-49fd-9aa8-bb1e21d28929" && req.Method == "GET" && strings.Contains(req.Header.Get("Authorization"), "Bearer my-special-token") {
                    response := `{
                      "id": "87dfc5b4-daf9-49fd-9aa8-bb1e21d28929",
                      "meta": {
                        "version": 6,
                        "created": "2014-05-22T22:36:36.941Z",
                        "lastModified": "2014-06-25T23:10:03.845Z"
                      },
                      "userName": "admin",
                      "name": {
                        "familyName": "Admin",
                        "givenName": "Mister"
                      },
                      "emails": [
                        {
                          "value": "fake-user@example.com"
                        }
                      ],
                      "groups": [
                        {
                          "value": "e7f74565-4c7e-44ba-b068-b16072cbf08f",
                          "display": "clients.read",
                          "type": "DIRECT"
                        }
                      ],
                      "approvals": [],
                      "active": true,
                      "verified": false,
                      "schemas": [
                        "urn:scim:schemas:core:1.0"
                      ]
                    }`
                    w.WriteHeader(http.StatusOK)
                    w.Write([]byte(response))
                } else {
                    w.WriteHeader(http.StatusNotFound)
                }
            }))
            auth = uaa.NewUAA("http://uaa.example.com", fakeUAAServer.URL, "the-client-id", "the-client-secret", "my-special-token")
        })

        AfterEach(func() {
            fakeUAAServer.Close()
        })

        It("returns a User from UAA", func() {
            user, err := uaa.UserByID(auth, "87dfc5b4-daf9-49fd-9aa8-bb1e21d28929")
            if err != nil {
                panic(err)
            }

            Expect(user).To(Equal(uaa.User{
                Username: "admin",
                ID:       "87dfc5b4-daf9-49fd-9aa8-bb1e21d28929",
                Name: uaa.Name{
                    FamilyName: "Admin",
                    GivenName:  "Mister",
                },
                Emails:   []string{"fake-user@example.com"},
                Active:   true,
                Verified: false,
            }))
        })
    })

    Context("when UAA is not responding normally", func() {
        BeforeEach(func() {
            fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                if req.URL.Path == "/Users/1234" && req.Method == "GET" && strings.Contains(req.Header.Get("Authorization"), "Bearer my-special-token") {
                    w.WriteHeader(http.StatusUnauthorized)
                    w.Write([]byte(`{"errors": "Unauthorized"}`))
                } else {
                    w.WriteHeader(http.StatusNotFound)
                }
            }))
            auth = uaa.NewUAA("http://uaa.example.com", fakeUAAServer.URL, "the-client-id", "the-client-secret", "my-special-token")
        })

        AfterEach(func() {
            fakeUAAServer.Close()
        })

        It("returns an error message", func() {
            _, err := uaa.UserByID(auth, "1234")
            Expect(err).To(BeAssignableToTypeOf(uaa.Failure{}))
            Expect(err.Error()).To(Equal(`UAA Failure: 401 {"errors": "Unauthorized"}`))
        })
    })
})
