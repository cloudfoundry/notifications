package uaa_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var usersFirstPage = []map[string]interface{}{
	map[string]interface{}{
		"id": "091b6583-0933-4d17-a5b6-66e54666c88e",
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
			{"value": "why-email@example.com"},
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
	map[string]interface{}{
		"id": "943e6076-b1a5-4404-811b-a1ee9253bf56",
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
			{"value": "slayer@example.com"},
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
	map[string]interface{}{
		"id": "646eb628-00d0-4c1e-957f-c54733fefb81",
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
			{"value": "the-yesman@example.com"},
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

var usersSecondPage = []map[string]interface{}{
	map[string]interface{}{
		"id": "8bd730bd-0a66-495d-a009-2bdaacfb2e50",
		"meta": map[string]interface{}{
			"version":      6,
			"created":      "2014-05-22T22:36:36.941Z",
			"lastModified": "2014-06-25T23:10:03.845Z",
		},
		"userName": "nothing",
		"name": map[string]string{
			"familyName": "Nada",
			"givenName":  "Mister",
		},
		"emails": []map[string]string{
			{"value": "my-example@example.com"},
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

var _ = Describe("AllUsers", func() {
	var fakeUAAServer *httptest.Server
	var auth uaa.UAA

	Context("when UAA is responding normally", func() {
		BeforeEach(func() {
			fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path == "/Users" && req.Method == "GET" && strings.Contains(req.Header.Get("Authorization"), "Bearer my-special-token") {
					responseObj := map[string]interface{}{
						"resources":    usersFirstPage,
						"startIndex":   1,
						"itemsPerPage": 3,
						"totalResults": 4,
						"schemas":      []string{"urn:scim:schemas:core:1.0"},
					}

					if strings.Contains(req.URL.RawQuery, "startIndex=4") {
						responseObj = map[string]interface{}{
							"resources":    usersSecondPage,
							"startIndex":   4,
							"itemsPerPage": 3,
							"totalResults": 4,
							"schemas":      []string{"urn:scim:schemas:core:1.0"},
						}
					}

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
			users, err := uaa.AllUsers(auth)
			if err != nil {
				panic(err)
			}

			user1 := uaa.User{
				Username: "admin",
				ID:       "091b6583-0933-4d17-a5b6-66e54666c88e",
				Name: uaa.Name{
					FamilyName: "Admin",
					GivenName:  "Mister",
				},
				Emails:   []string{"why-email@example.com"},
				Active:   true,
				Verified: false,
			}

			user2 := uaa.User{
				Username: "some-user",
				ID:       "943e6076-b1a5-4404-811b-a1ee9253bf56",
				Name: uaa.Name{
					FamilyName: "Some",
					GivenName:  "User",
				},
				Emails:   []string{"slayer@example.com"},
				Active:   true,
				Verified: false,
			}

			user3 := uaa.User{
				Username: "other",
				ID:       "646eb628-00d0-4c1e-957f-c54733fefb81",
				Name: uaa.Name{
					FamilyName: "Other",
					GivenName:  "User",
				},
				Emails:   []string{"the-yesman@example.com"},
				Active:   true,
				Verified: false,
			}

			user4 := uaa.User{
				Username: "nothing",
				ID:       "8bd730bd-0a66-495d-a009-2bdaacfb2e50",
				Name: uaa.Name{
					FamilyName: "Nada",
					GivenName:  "Mister",
				},
				Emails:   []string{"my-example@example.com"},
				Active:   true,
				Verified: false,
			}

			Expect(users).To(ConsistOf(user1, user2, user3, user4))
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
			_, err := uaa.AllUsers(auth)
			Expect(err).To(BeAssignableToTypeOf(uaa.Failure{}))
			Expect(err.Error()).To(Equal(`UAA Failure: 401 {"errors": "Unauthorized"}`))
		})
	})
})
