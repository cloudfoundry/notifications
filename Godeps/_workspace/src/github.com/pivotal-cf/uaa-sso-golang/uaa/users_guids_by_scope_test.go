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

var _ = Describe("UsersGUIDsByScope", func() {
	var fakeUAAServer *httptest.Server
	var auth uaa.UAA
	var users []map[string][]map[string]string
	var requestCount int

	BeforeEach(func() {
		requestCount = 0
		users = []map[string][]map[string]string{
			{
				"members": {
					{
						"origin": "uaa",
						"type":   "USER",
						"value":  "my-water-bottle-guid",
					},
					{
						"origin": "uaa",
						"type":   "USER",
						"value":  "my-other-guid",
					},
					{
						"origin": "uaa",
						"type":   "USER",
						"value":  "my-guid",
					},
				},
			},
		}

		fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			requestCount += 1

			if req.URL.Path == `/Groups` && req.Method == "GET" && strings.Contains(req.Header.Get("Authorization"), "Bearer my-special-token") {
				responseObj := map[string]interface{}{
					"resources":    users,
					"startIndex":   1,
					"itemsPerPage": 100,
					"totalResults": 3,
					"schemas":      []string{"urn:scim:schemas:core:1.0"},
				}

				err := req.ParseForm()
				if err != nil {
					panic(err)
				}

				if req.Form.Get("attributes") != "members" {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte{})
					return
				}

				filter := req.FormValue("filter")
				matcher := regexp.MustCompile(`displayName eq "this\.scope"`)
				matches := matcher.FindAllStringSubmatch(filter, -1)
				if matches == nil {
					w.WriteHeader(http.StatusBadRequest)
					return
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

	It("returns slice of GUIDs from UAA", func() {
		users, err := uaa.UsersGUIDsByScope(auth, "this.scope")

		Expect(err).NotTo(HaveOccurred())
		Expect(users).To(Equal([]string{"my-water-bottle-guid", "my-other-guid", "my-guid"}))
	})

	Context("when the members response is empty", func() {
		It("does not blow up", func() {
			users = []map[string][]map[string]string{
				{
					"Members": {},
				},
			}

			guids, err := uaa.UsersGUIDsByScope(auth, "this.scope")

			Expect(err).NotTo(HaveOccurred())
			Expect(guids).To(Equal([]string{}))
		})
	})

	Context("when UAA is not responding normally", func() {
		BeforeEach(func() {
			fakeUAAServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path == "/Groups" && req.Method == "GET" && strings.Contains(req.Header.Get("Authorization"), "Bearer my-special-token") {
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
			_, err := uaa.UsersGUIDsByScope(auth, "this.scope")
			Expect(err).To(BeAssignableToTypeOf(uaa.Failure{}))
			Expect(err.Error()).To(Equal(`UAA Failure: 401 {"errors": "Unauthorized"}`))
		})
	})
})
