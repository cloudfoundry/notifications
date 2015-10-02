package docs_test

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/docs"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DocsGenerator", func() {
	Describe("Add", func() {
		Context("when a list request is added to the document generator", func() {
			It("should add a new MethodEntry to ListMethodEntries", func() {
				requestInspector := mocks.NewRequestInspector()
				requestInspector.GetResourceInfoCall.Returns.ResourceInfo = docs.ResourceInfo{
					ResourceType: "some-resource-types",
					ListName:     "Some resource types",
					ItemName:     "Some resource type",
				}

				docGenerator := docs.NewDocGenerator(requestInspector)

				request, err := http.NewRequest("GET", "/some-resource-types", nil)
				Expect(err).NotTo(HaveOccurred())

				request.Header = http.Header{
					"X-NOTIFICATIONS-VERSION": []string{"2"},
					"Authorization":           []string{"bearer some-client-token"},
				}

				response := &http.Response{
					StatusCode: http.StatusOK,
					Header: http.Header{
						"X-Cf-Requestid": []string{"some-request-id"},
					},
					Body:    ioutil.NopCloser(bytes.NewReader([]byte(`{ "some": "response" }`))),
					Request: request,
				}

				Expect(docGenerator.Add(request, response)).To(Succeed())

				Expect(docGenerator.Resources).To(HaveKey("some-resource-types"))

				resourceEntry := docGenerator.Resources["some-resource-types"]
				Expect(resourceEntry.ListResourceName).To(Equal("Some resource types"))
				Expect(resourceEntry.ItemResourceName).To(Equal("Some resource type"))

				Expect(requestInspector.GetResourceInfoCall.Receives.Request).To(Equal(request))

				Expect(resourceEntry.ListMethodEntries).To(HaveLen(1))
				Expect(resourceEntry.ItemMethodEntries).To(HaveLen(0))
				Expect(resourceEntry.ListMethodEntries[0].Verb).To(Equal("GET"))

				Expect(resourceEntry.ListMethodEntries[0].Request.Headers).To(HaveKeyWithValue("X-NOTIFICATIONS-VERSION", []string{"2"}))
				Expect(resourceEntry.ListMethodEntries[0].Request.Headers).To(HaveKeyWithValue("Authorization", []string{"bearer some-client-token"}))

				Expect(resourceEntry.ListMethodEntries[0].Responses).To(HaveLen(1))
				Expect(resourceEntry.ListMethodEntries[0].Responses[0].Code).To(Equal(http.StatusOK))
				Expect(resourceEntry.ListMethodEntries[0].Responses[0].Headers).To(HaveKeyWithValue("X-Cf-Requestid", []string{"some-request-id"}))
				Expect(resourceEntry.ListMethodEntries[0].Responses[0].Body).To(Equal(`{ "some": "response" }`))
			})

			It("should save the authorization header", func() {
				requestInspector := mocks.NewRequestInspector()
				requestInspector.GetResourceInfoCall.Returns.ResourceInfo = docs.ResourceInfo{
					ResourceType:       "some-resource-types",
					ListName:           "Some resource types",
					ItemName:           "Some resource type",
					AuthorizationToken: "some-client-token",
				}

				docGenerator := docs.NewDocGenerator(requestInspector)

				request, err := http.NewRequest("GET", "/some-resource-types", nil)
				Expect(err).NotTo(HaveOccurred())

				response := &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{ "some": "response" }`))),
					Request:    request,
				}

				Expect(docGenerator.Add(request, response)).To(Succeed())
				Expect(docGenerator.SampleAuthorizationToken).To(Equal("some-client-token"))
			})
		})

		Context("when an item request is added to the document generator", func() {
			It("should add a new MethodEntry to ItemMethodEntries", func() {
				requestInspector := mocks.NewRequestInspector()
				requestInspector.GetResourceInfoCall.Returns.ResourceInfo = docs.ResourceInfo{
					ResourceType: "some-resource-types",
					ListName:     "Some resource types",
					ItemName:     "Some resource type",
					IsItem:       true,
				}

				docGenerator := docs.NewDocGenerator(requestInspector)

				request, err := http.NewRequest("GET", "/some-resource-types/79b6bc02-60a1-11e5-b1e2-6f64bf254a5d", nil)
				Expect(err).NotTo(HaveOccurred())

				request.Header = http.Header{
					"X-NOTIFICATIONS-VERSION": []string{"2"},
					"Authorization":           []string{"bearer some-client-token"},
				}

				response := &http.Response{
					StatusCode: http.StatusOK,
					Header: http.Header{
						"X-Cf-Requestid": []string{"some-request-id"},
					},
					Body:    ioutil.NopCloser(bytes.NewReader([]byte(`{ "some": "response" }`))),
					Request: request,
				}

				Expect(docGenerator.Add(request, response)).To(Succeed())

				Expect(docGenerator.Resources).To(HaveKey("some-resource-types"))

				Expect(docGenerator.Resources["some-resource-types"].ListMethodEntries).To(HaveLen(0))
				Expect(docGenerator.Resources["some-resource-types"].ItemMethodEntries).To(HaveLen(1))
			})
		})

		Context("when two item requests are added to the document generator", func() {
			It("should add two entries to ItemMethodEntries", func() {
				requestInspector := mocks.NewRequestInspector()
				requestInspector.GetResourceInfoCall.Returns.ResourceInfo = docs.ResourceInfo{
					ResourceType: "some-resource-types",
					ListName:     "Some resource types",
					ItemName:     "Some resource type",
					IsItem:       true,
				}

				docGenerator := docs.NewDocGenerator(requestInspector)

				request, err := http.NewRequest("GET", "/some-resource-types/79b6bc02-60a1-11e5-b1e2-6f64bf254a5d", nil)
				Expect(err).NotTo(HaveOccurred())

				request.Header = http.Header{
					"X-NOTIFICATIONS-VERSION": []string{"2"},
					"Authorization":           []string{"bearer some-client-token"},
				}

				response := &http.Response{
					StatusCode: http.StatusOK,
					Header: http.Header{
						"X-Cf-Requestid": []string{"some-request-id"},
					},
					Body:    ioutil.NopCloser(bytes.NewReader([]byte(`{ "some": "response" }`))),
					Request: request,
				}

				Expect(docGenerator.Add(request, response)).To(Succeed())
				Expect(docGenerator.Add(request, response)).To(Succeed())

				Expect(docGenerator.Resources).To(HaveKey("some-resource-types"))

				Expect(docGenerator.Resources["some-resource-types"].ListMethodEntries).To(HaveLen(0))
				Expect(docGenerator.Resources["some-resource-types"].ItemMethodEntries).To(HaveLen(2))
			})
		})
	})
})
