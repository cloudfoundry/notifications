package rainmaker_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/pivotal-cf-experimental/rainmaker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var unsupportedJSONType = func() {}

var _ = Describe("Client", func() {
	var client rainmaker.Client

	BeforeEach(func() {
		client = rainmaker.NewClient(rainmaker.Config{Host: "http://example.com"})
	})

	It("has an organizations service", func() {
		Expect(client.Organizations).To(BeAssignableToTypeOf(&rainmaker.OrganizationsService{}))
		Expect(client.Organizations).NotTo(BeNil())
	})

	It("has an spaces service", func() {
		Expect(client.Spaces).To(BeAssignableToTypeOf(&rainmaker.SpacesService{}))
		Expect(client.Spaces).NotTo(BeNil())
	})

	It("has a service instance service", func() {
		Expect(client.ServiceInstances).To(BeAssignableToTypeOf(&rainmaker.ServiceInstancesService{}))
		Expect(client.ServiceInstances).NotTo(BeNil())
	})

	Describe("makeRequest", func() {
		Context("when an error occurs", func() {
			It("returns a RequestBodyMarshalError when the request body cannot be marshalled", func() {
				requestArgs := rainmaker.NewRequestArguments("GET", "/path", "token", unsupportedJSONType, []int{http.StatusOK})
				_, _, err := client.MakeRequest(requestArgs)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(rainmaker.RequestBodyMarshalError{}))
			})

			It("returns a RequestConfigurationError when the request params are bad", func() {
				client.Config.Host = "://example.com"
				requestArgs := rainmaker.NewRequestArguments("GET", "/path", "token", nil, []int{http.StatusOK})
				_, _, err := client.MakeRequest(requestArgs)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(rainmaker.RequestConfigurationError{}))
			})

			It("returns a RequestHTTPError when the request fails", func() {
				client.Config.Host = "banana://example.com"
				requestArgs := rainmaker.NewRequestArguments("GET", "/path", "token", nil, []int{http.StatusOK})
				_, _, err := client.MakeRequest(requestArgs)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(rainmaker.RequestHTTPError{}))
			})

			It("returns a ResponseReadError when the response cannot be read", func() {
				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Header().Set("Content-Length", "100")
					w.Write([]byte("{}"))
				}))
				client.Config.Host = fakeServer.URL

				requestArgs := rainmaker.NewRequestArguments("GET", "/path", "token", nil, []int{http.StatusOK})
				_, _, err := client.MakeRequest(requestArgs)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(rainmaker.ResponseReadError{}))

				fakeServer.Close()
			})

			It("returns an UnexpectedStatusError when the response status is not an expected value", func() {
				fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusTeapot)
					w.Write([]byte(`I'm a little teapot`))
				}))
				client.Config.Host = fakeServer.URL

				requestArgs := rainmaker.NewRequestArguments("GET", "/path", "token", nil, []int{http.StatusOK})
				_, _, err := client.MakeRequest(requestArgs)
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(rainmaker.UnexpectedStatusError{}))
				Expect(err).To(MatchError("Rainmaker UnexpectedStatusError: 418 I'm a little teapot"))

				fakeServer.Close()
			})
		})
	})

	Describe("unmarshal", func() {
		It("returns a JSON unmarshalled representation of the given bytes", func() {
			response := make(map[string]interface{})
			body := []byte(`{"greeting":"hello", "answer":true}`)
			err := client.Unmarshal(body, &response)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).To(HaveKeyWithValue("greeting", "hello"))
			Expect(response).To(HaveKeyWithValue("answer", true))
		})

		It("returns a ResponseBodyUnmarshalError when the body cannot be unmarshaled", func() {
			var response map[string]interface{}
			body := []byte(`This is not JSON`)
			err := client.Unmarshal(body, &response)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(rainmaker.ResponseBodyUnmarshalError{}))
		})
	})
})
