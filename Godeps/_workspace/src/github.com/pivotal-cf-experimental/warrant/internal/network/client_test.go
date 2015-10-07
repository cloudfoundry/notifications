package network_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pivotal-cf-experimental/warrant/internal/network"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var unsupportedJSONType = func() {}

type Request struct {
	Body   string
	Header http.Header
}

var _ = Describe("Client", func() {
	var (
		token           string
		fakeServer      *httptest.Server
		client          network.Client
		receivedRequest *Request
	)

	BeforeEach(func() {
		token = "TOKEN"
		receivedRequest = &Request{}
		fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Body != nil {
				body, err := ioutil.ReadAll(req.Body)
				if err != nil {
					panic(err)
				}

				receivedRequest.Body = string(body)
			}

			receivedRequest.Header = req.Header
		}))

		client = network.NewClient(network.Config{
			Host:          fakeServer.URL,
			SkipVerifySSL: true,
			TraceWriter:   TraceWriter,
		})
	})

	AfterEach(func() {
		fakeServer.Close()
	})

	Describe("makeRequest", func() {
		It("can make requests", func() {
			jsonBody := map[string]interface{}{
				"hello": "goodbye",
			}

			_, err := client.MakeRequest(network.Request{
				Method:        "GET",
				Path:          "/path",
				Authorization: network.NewTokenAuthorization(token),
				Body:          network.NewJSONRequestBody(jsonBody),
				AcceptableStatusCodes: []int{http.StatusOK},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(receivedRequest.Body).To(MatchJSON(`{"hello": "goodbye"}`))
			Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept", []string{"application/json"}))
			Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept-Encoding", []string{"gzip"}))
			Expect(receivedRequest.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer TOKEN"}))
			Expect(receivedRequest.Header).To(HaveKeyWithValue("Content-Length", []string{"19"}))
			Expect(receivedRequest.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
		})

		It("can make more requests than the total allowed number of open files", func() {
			var output []byte

			_, err := exec.LookPath("ulimit")
			if err != nil {
				var err error
				output, err = ioutil.ReadFile("/proc/sys/fs/nr_open")
				Expect(err).NotTo(HaveOccurred())
			} else {
				cmd := exec.Command("ulimit", "-n")

				var err error
				output, err = cmd.Output()
				Expect(err).NotTo(HaveOccurred())
			}
			fdCount, err := strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64)
			Expect(err).NotTo(HaveOccurred())

			for i := 0; i < int(fdCount)+10; i++ {
				_, err := client.MakeRequest(network.Request{
					Method: "GET",
					Path:   "/path",
					AcceptableStatusCodes: []int{http.StatusOK},
				})
				Expect(err).NotTo(HaveOccurred())
			}
		})

		Context("Following redirects", func() {
			var requestArgs network.Request

			BeforeEach(func() {
				redirectingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					if req.URL.Path == "/redirect" {
						w.Header().Set("Location", "/noredirect")
						w.WriteHeader(http.StatusFound)
						return
					}

					w.Write([]byte("did not redirect"))
				}))

				client = network.NewClient(network.Config{
					Host:          redirectingServer.URL,
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				requestArgs = network.Request{
					Method:                "GET",
					Path:                  "/redirect",
					Authorization:         network.NewTokenAuthorization(token),
					AcceptableStatusCodes: []int{http.StatusFound, http.StatusOK},
				}
			})

			Context("when DoNotFollowRedirects is not set", func() {
				It("follows redirects to their location", func() {
					resp, err := client.MakeRequest(requestArgs)
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.Code).To(Equal(http.StatusOK))
					Expect(resp.Headers.Get("Location")).To(Equal(""))
					Expect(resp.Body).To(ContainSubstring("did not redirect"))
				})
			})

			Context("when DoNotFollowRedirects is set", func() {
				It("does not follow redirects to their location", func() {
					requestArgs.DoNotFollowRedirects = true
					resp, err := client.MakeRequest(requestArgs)
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.Code).To(Equal(http.StatusFound))
					Expect(resp.Headers.Get("Location")).To(Equal("/noredirect"))
				})
			})
		})

		Context("Headers", func() {
			Context("authorization", func() {
				It("does not include Authorization header when there is no authorization", func() {
					requestArgs := network.Request{
						Method: "GET",
						Path:   "/path",
						Body:   network.NewJSONRequestBody(map[string]string{"hello": "world"}),
						AcceptableStatusCodes: []int{http.StatusOK},
					}

					_, err := client.MakeRequest(requestArgs)
					Expect(err).NotTo(HaveOccurred())
					Expect(receivedRequest.Body).To(MatchJSON(`{"hello": "world"}`))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept", []string{"application/json"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept-Encoding", []string{"gzip"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Content-Length", []string{"17"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
				})

				It("includes a bearer Authorization header when there is token authorization", func() {
					requestArgs := network.Request{
						Method:        "GET",
						Path:          "/path",
						Authorization: network.NewTokenAuthorization("TOKEN"),
						Body:          network.NewJSONRequestBody(map[string]string{"hello": "world"}),
						AcceptableStatusCodes: []int{http.StatusOK},
					}

					_, err := client.MakeRequest(requestArgs)
					Expect(err).NotTo(HaveOccurred())
					Expect(receivedRequest.Body).To(MatchJSON(`{"hello": "world"}`))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept", []string{"application/json"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept-Encoding", []string{"gzip"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer TOKEN"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Content-Length", []string{"17"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
				})

				It("includes a basic Authorization header when there is basic authorization", func() {
					requestArgs := network.Request{
						Method:        "GET",
						Path:          "/path",
						Authorization: network.NewBasicAuthorization("username", "password"),
						Body:          network.NewJSONRequestBody(map[string]string{"hello": "world"}),
						AcceptableStatusCodes: []int{http.StatusOK},
					}
					_, err := client.MakeRequest(requestArgs)
					Expect(err).NotTo(HaveOccurred())
					Expect(receivedRequest.Body).To(MatchJSON(`{"hello": "world"}`))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept", []string{"application/json"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept-Encoding", []string{"gzip"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Authorization", []string{"Basic dXNlcm5hbWU6cGFzc3dvcmQ="}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Content-Length", []string{"17"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
				})
			})

			Context("when there is a JSON body", func() {
				It("includes the Content-Type header in the request", func() {
					requestArgs := network.Request{
						Method:        "GET",
						Path:          "/path",
						Authorization: network.NewTokenAuthorization(token),
						Body:          network.NewJSONRequestBody(map[string]string{"hello": "world"}),
						AcceptableStatusCodes: []int{http.StatusOK},
					}
					_, err := client.MakeRequest(requestArgs)
					Expect(err).NotTo(HaveOccurred())
					Expect(receivedRequest.Body).To(MatchJSON(`{"hello": "world"}`))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept", []string{"application/json"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept-Encoding", []string{"gzip"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer TOKEN"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Content-Length", []string{"17"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
				})
			})

			Context("when there is no JSON body", func() {
				It("does not include the Content-Type or Content-Length headers in the request", func() {
					requestArgs := network.Request{
						Method:        "GET",
						Path:          "/path",
						Authorization: network.NewTokenAuthorization(token),
						Body:          nil,
						AcceptableStatusCodes: []int{http.StatusOK},
					}

					_, err := client.MakeRequest(requestArgs)
					Expect(err).NotTo(HaveOccurred())
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept", []string{"application/json"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept-Encoding", []string{"gzip"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer TOKEN"}))
					Expect(receivedRequest.Header).NotTo(HaveKey("Content-Type"))
					Expect(receivedRequest.Header).NotTo(HaveKey("Content-Length"))
				})
			})

			Context("when the If-Match argument is assigned", func() {
				It("includes the header in the request", func() {
					requestArgs := network.Request{
						Method:        "GET",
						Path:          "/path",
						Authorization: network.NewTokenAuthorization(token),
						IfMatch:       "45",
						Body:          nil,
						AcceptableStatusCodes: []int{http.StatusOK},
					}

					_, err := client.MakeRequest(requestArgs)
					Expect(err).NotTo(HaveOccurred())
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept", []string{"application/json"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept-Encoding", []string{"gzip"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer TOKEN"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("If-Match", []string{"45"}))
				})
			})

			Context("when the If-Match argument is not assigned", func() {
				It("does not include the header in the request", func() {
					requestArgs := network.Request{
						Method:        "GET",
						Path:          "/path",
						Authorization: network.NewTokenAuthorization(token),
						Body:          nil,
						AcceptableStatusCodes: []int{http.StatusOK},
					}

					_, err := client.MakeRequest(requestArgs)
					Expect(err).NotTo(HaveOccurred())
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept", []string{"application/json"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Accept-Encoding", []string{"gzip"}))
					Expect(receivedRequest.Header).To(HaveKeyWithValue("Authorization", []string{"Bearer TOKEN"}))
					Expect(receivedRequest.Header).NotTo(HaveKey("If-Match"))
				})
			})
		})

		Context("when errors occur", func() {
			It("returns a RequestBodyEncodeError when the request body cannot be encoded", func() {
				requestArgs := network.Request{
					Method:        "GET",
					Path:          "/path",
					Authorization: network.NewTokenAuthorization(token),
					Body:          network.NewJSONRequestBody(unsupportedJSONType),
					AcceptableStatusCodes: []int{http.StatusOK},
				}

				_, err := client.MakeRequest(requestArgs)
				Expect(err).To(BeAssignableToTypeOf(network.RequestBodyEncodeError{}))
			})

			It("returns a RequestConfigurationError when the request params are bad", func() {
				client = network.NewClient(network.Config{
					Host:          "://example.com",
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				requestArgs := network.Request{
					Method:        "GET",
					Path:          "/path",
					Authorization: network.NewTokenAuthorization(token),
					Body:          nil,
					AcceptableStatusCodes: []int{http.StatusOK},
				}
				_, err := client.MakeRequest(requestArgs)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(network.RequestConfigurationError{}))
			})

			It("returns a RequestHTTPError when the request fails", func() {
				client = network.NewClient(network.Config{
					Host:          "banana://example.com",
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				requestArgs := network.Request{
					Method:        "GET",
					Path:          "/path",
					Authorization: network.NewTokenAuthorization(token),
					Body:          nil,
					AcceptableStatusCodes: []int{http.StatusOK},
				}
				_, err := client.MakeRequest(requestArgs)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(network.RequestHTTPError{}))
			})

			It("returns a ResponseReadError when the response cannot be read", func() {
				unintelligibleServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Header().Set("Content-Length", "100")
					w.Write([]byte("{}"))
				}))

				client = network.NewClient(network.Config{
					Host:          unintelligibleServer.URL,
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				requestArgs := network.Request{
					Method:        "GET",
					Path:          "/path",
					Authorization: network.NewTokenAuthorization(token),
					Body:          nil,
					AcceptableStatusCodes: []int{http.StatusOK},
				}
				_, err := client.MakeRequest(requestArgs)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(network.ResponseReadError{}))

				unintelligibleServer.Close()
			})

			It("returns an UnexpectedStatusError when the response status is not an expected value", func() {
				requestArgs := network.Request{
					Method:        "GET",
					Path:          "/path",
					Authorization: network.NewTokenAuthorization(token),
					Body:          nil,
					AcceptableStatusCodes: []int{http.StatusTeapot},
				}
				_, err := client.MakeRequest(requestArgs)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(network.UnexpectedStatusError{}))
			})

			It("returns a NotFoundError when the response status is 404", func() {
				missingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))

				client = network.NewClient(network.Config{
					Host:          missingServer.URL,
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				requestArgs := network.Request{
					Method:        "GET",
					Path:          "/path",
					Authorization: network.NewTokenAuthorization(token),
					Body:          nil,
					AcceptableStatusCodes: []int{http.StatusOK},
				}
				_, err := client.MakeRequest(requestArgs)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(network.NotFoundError{}))

				missingServer.Close()
			})

			It("returns an UnauthorizedError when the response status is 401 Unauthorized", func() {
				lockedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
				}))

				client = network.NewClient(network.Config{
					Host:          lockedServer.URL,
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				requestArgs := network.Request{
					Method:        "GET",
					Path:          "/path",
					Authorization: network.NewTokenAuthorization(token),
					Body:          nil,
					AcceptableStatusCodes: []int{http.StatusOK},
				}
				_, err := client.MakeRequest(requestArgs)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(network.UnauthorizedError{}))

				lockedServer.Close()
			})

			It("returns an UnauthorizedError when the response status is 403 Forbidden", func() {
				forbiddenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusForbidden)
				}))

				client = network.NewClient(network.Config{
					Host:          forbiddenServer.URL,
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				requestArgs := network.Request{
					Method:        "GET",
					Path:          "/path",
					Authorization: network.NewTokenAuthorization(token),
					Body:          nil,
					AcceptableStatusCodes: []int{http.StatusOK},
				}
				_, err := client.MakeRequest(requestArgs)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(network.UnauthorizedError{}))

				forbiddenServer.Close()
			})
		})
	})
})
