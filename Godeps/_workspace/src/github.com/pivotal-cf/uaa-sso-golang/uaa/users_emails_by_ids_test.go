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

var _ = Describe("UsersEmailsByIds", func() {

	var fakeUAAServer *httptest.Server
	var auth uaa.UAA
	var users map[string]map[string]interface{}
	var requestCount int

	BeforeEach(func() {
		requestCount = 0
		users = map[string]map[string]interface{}{
			"87dfc5b4-daf9-49fd-9aa8-bb1e21d28929": map[string]interface{}{
				"emails": []map[string]string{
					{"value": "fake-user@example.com"},
				},
				"id": "87dfc5b4-daf9-49fd-9aa8-bb1e21d28929",
			},
			"21f1c87b-0c4b-4dbb-a1e0-1ba479e8aed3": map[string]interface{}{
				"emails": []map[string]string{
					{"value": "some-user@example.com"},
				},
				"id": "21f1c87b-0c4b-4dbb-a1e0-1ba479e8aed3",
			},
			"baf908c9-3248-451f-ab3c-103d921cd61e": map[string]interface{}{
				"emails": []map[string]string{
					{"value": "other-user@example.com"},
				},
				"id": "baf908c9-3248-451f-ab3c-103d921cd61e",
			},
		}
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

				if req.Form.Get("attributes") != "emails,id" {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte{})
					return
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

	It("returns slice of Users from UAA", func() {
		users, err := uaa.UsersEmailsByIDs(auth, "87dfc5b4-daf9-49fd-9aa8-bb1e21d28929", "baf908c9-3248-451f-ab3c-103d921cd61e")
		if err != nil {
			panic(err)
		}

		user1 := uaa.User{
			Emails: []string{"fake-user@example.com"},
			ID:     "87dfc5b4-daf9-49fd-9aa8-bb1e21d28929",
		}

		user2 := uaa.User{
			Emails: []string{"other-user@example.com"},
			ID:     "baf908c9-3248-451f-ab3c-103d921cd61e",
		}

		Expect(users).To(Equal([]uaa.User{user1, user2}))
	})

	It("respects the maximum length of a URL", func() {
		users, err := uaa.UsersEmailsByIDsWithMaxLength(auth, 110, "87dfc5b4-daf9-49fd-9aa8-bb1e21d28929", "baf908c9-3248-451f-ab3c-103d921cd61e")
		if err != nil {
			panic(err)
		}

		Expect(len(users)).To(Equal(2))
		Expect(users).To(ContainElement(uaa.User{
			Emails: []string{"fake-user@example.com"},
			ID:     "87dfc5b4-daf9-49fd-9aa8-bb1e21d28929",
		}))
		Expect(users).To(ContainElement(uaa.User{
			Emails: []string{"other-user@example.com"},
			ID:     "baf908c9-3248-451f-ab3c-103d921cd61e",
		}))

		Expect(requestCount).To(Equal(2))
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
			_, err := uaa.UsersEmailsByIDs(auth, "1234", "5678")
			Expect(err).To(BeAssignableToTypeOf(uaa.Failure{}))
			Expect(err.Error()).To(Equal(`UAA Failure: 401 {"errors": "Unauthorized"}`))
		})
	})
})
