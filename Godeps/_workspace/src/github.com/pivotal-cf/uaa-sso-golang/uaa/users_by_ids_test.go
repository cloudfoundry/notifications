package uaa_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "regexp"
    "strings"

    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var users = map[string]map[string]interface{}{
    "87dfc5b4-daf9-49fd-9aa8-bb1e21d28929": map[string]interface{}{
        "id": "87dfc5b4-daf9-49fd-9aa8-bb1e21d28929",
        "meta": map[string]interface{}{
            "version":      6,
            "created":      "2014-05-22T22:36:36.941Z",
            "lastModified": "2014-06-25T23:10:03.845Z",
        },
        "userName": "admin",
        "name": map[string]string{
            "familyName": "Admin",
            "givenName":  "Mister",
        },
        "emails": []map[string]string{
            {"value": "fake-user@example.com"},
        },
        "groups": []map[string]string{
            {
                "value":   "e7f74565-4c7e-44ba-b068-b16072cbf08f",
                "display": "clients.read",
                "type":    "DIRECT",
            },
        },
        "approvals": []string{},
        "active":    true,
        "verified":  false,
        "schemas":   []string{"urn:scim:schemas:core:1.0"},
    },
    "21f1c87b-0c4b-4dbb-a1e0-1ba479e8aed3": map[string]interface{}{
        "id": "21f1c87b-0c4b-4dbb-a1e0-1ba479e8aed3",
        "meta": map[string]interface{}{
            "version":      6,
            "created":      "2014-05-22T22:36:36.941Z",
            "lastModified": "2014-06-25T23:10:03.845Z",
        },
        "userName": "some-user",
        "name": map[string]string{
            "familyName": "Some",
            "givenName":  "User",
        },
        "emails": []map[string]string{
            {"value": "some-user@example.com"},
        },
        "groups": []map[string]string{
            {
                "value":   "e7f74565-4c7e-44ba-b068-b16072cbf08f",
                "display": "clients.read",
                "type":    "DIRECT",
            },
        },
        "approvals": []string{},
        "active":    true,
        "verified":  false,
        "schemas":   []string{"urn:scim:schemas:core:1.0"},
    },
    "baf908c9-3248-451f-ab3c-103d921cd61e": map[string]interface{}{
        "id": "baf908c9-3248-451f-ab3c-103d921cd61e",
        "meta": map[string]interface{}{
            "version":      6,
            "created":      "2014-05-22T22:36:36.941Z",
            "lastModified": "2014-06-25T23:10:03.845Z",
        },
        "userName": "other",
        "name": map[string]string{
            "familyName": "Other",
            "givenName":  "User",
        },
        "emails": []map[string]string{
            {"value": "other-user@example.com"},
        },
        "groups": []map[string]string{
            {
                "value":   "e7f74565-4c7e-44ba-b068-b16072cbf08f",
                "display": "clients.read",
                "type":    "DIRECT",
            },
        },
        "approvals": []string{},
        "active":    true,
        "verified":  false,
        "schemas":   []string{"urn:scim:schemas:core:1.0"},
    },
}

var _ = Describe("UsersByIds", func() {
    var fakeUAAServer *httptest.Server
    var auth uaa.UAA
    var requestCount int

    Context("when UAA is responding normally", func() {
        BeforeEach(func() {
            requestCount = 0
            fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                requestCount += 1
                if req.URL.Path == "/Users" && req.Method == "GET" && strings.Contains(req.Header.Get("Authorization"), "Bearer my-special-token") {
                    responseObj := map[string]interface{}{
                        "resources":    []interface{}{},
                        "startIndex":   1,
                        "itemsPerPage": 100,
                        "totalResults": 3,
                        "schemas":      []string{"urn:scim:schemas:core:1.0"},
                    }

                    err := req.ParseForm()
                    if err != nil {
                        panic(err)
                    }

                    filter := req.FormValue("filter")
                    matcher := regexp.MustCompile(`Id eq "([a-zA-Z0-9\-]*)"`)
                    matches := matcher.FindAllStringSubmatch(filter, -1)

                    usersList := []interface{}{}
                    for _, match := range matches {
                        if user, ok := users[match[1]]; ok {
                            usersList = append(usersList, user)
                        }
                    }

                    responseObj["resources"] = usersList

                    response, err := json.Marshal(responseObj)
                    if err != nil {
                        panic(err)
                    }

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

        It("returns slice of Users from UAA", func() {
            users, err := uaa.UsersByIDs(auth, "87dfc5b4-daf9-49fd-9aa8-bb1e21d28929", "baf908c9-3248-451f-ab3c-103d921cd61e")
            if err != nil {
                panic(err)
            }

            user1 := uaa.User{
                Username: "admin",
                ID:       "87dfc5b4-daf9-49fd-9aa8-bb1e21d28929",
                Name: uaa.Name{
                    FamilyName: "Admin",
                    GivenName:  "Mister",
                },
                Emails:   []string{"fake-user@example.com"},
                Active:   true,
                Verified: false,
            }

            user2 := uaa.User{
                Username: "other",
                ID:       "baf908c9-3248-451f-ab3c-103d921cd61e",
                Name: uaa.Name{
                    FamilyName: "Other",
                    GivenName:  "User",
                },
                Emails:   []string{"other-user@example.com"},
                Active:   true,
                Verified: false,
            }

            Expect(users).To(Equal([]uaa.User{user1, user2}))
        })

        It("respects the maximum length of a URL", func() {
            users, err := uaa.UsersByIDsWithMaxLength(auth, 100, "87dfc5b4-daf9-49fd-9aa8-bb1e21d28929", "baf908c9-3248-451f-ab3c-103d921cd61e")
            if err != nil {
                panic(err)
            }

            Expect(len(users)).To(Equal(2))
            Expect(users[0].ID).To(Equal("87dfc5b4-daf9-49fd-9aa8-bb1e21d28929"))
            Expect(users[1].ID).To(Equal("baf908c9-3248-451f-ab3c-103d921cd61e"))

            Expect(requestCount).To(Equal(2))
        })
    })

    Context("when UAA is not responding normally", func() {
        BeforeEach(func() {
            fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
                if req.URL.Path == "/Users" && req.Method == "GET" && strings.Contains(req.Header.Get("Authorization"), "Bearer my-special-token") {
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
            _, err := uaa.UsersByIDs(auth, "1234", "5678")
            Expect(err).To(BeAssignableToTypeOf(uaa.Failure{}))
            Expect(err.Error()).To(Equal(`UAA Failure: 401 {"errors": "Unauthorized"}`))
        })
    })
})
