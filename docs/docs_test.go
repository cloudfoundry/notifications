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

const (
	tokenFixture      = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQzMjI1NTEsImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NjM3NzYvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.GbFa6KEgSarJYqU2yohnJj7G1ztoG4az8Ded1EcDVMbHf9OGKzuzXUKXPqruNzX3JhjwBlBN4f6P8yZXvhPmvnQSTG7a2lqE817lPs7L7c5J5ka6MBRWZqlJwquYOalv9ytn7ZpaodWNqGK_uH5ctpP9LFIoxG2WmKQsd1N7f_Wv8GyEbrHL6qMze8x1QhNuBos_zdOPuQ3DyMsJBlx7b35h2nY_ZXXDxVErGlALfrVfrUnAPbt7fgkRbX9JYtbzgnkhsrDXcPt_5xO3ideqYZARfrjKtCFQiR4kPqRADIGvjgaSDyX5fNvKpwmLfElmu_rCyr9tN7uZszShhwy9fQ"
	authHeaderFixture = "bearer " + tokenFixture
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
					"Authorization":           []string{authHeaderFixture},
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
				Expect(resourceEntry.ListMethodEntries[0].Request.Headers).To(HaveKeyWithValue("Authorization", []string{"bearer " + tokenFixture}))

				Expect(resourceEntry.ListMethodEntries[0].Responses).To(HaveLen(1))
				Expect(resourceEntry.ListMethodEntries[0].Responses[0].Code).To(Equal(http.StatusOK))
				Expect(resourceEntry.ListMethodEntries[0].Responses[0].Headers).To(HaveKeyWithValue("X-Cf-Requestid", []string{"some-request-id"}))
				Expect(resourceEntry.ListMethodEntries[0].Responses[0].Body).To(Equal(`{ "some": "response" }`))

				Expect(resourceEntry.ListMethodEntries[0].Request.Scopes).To(Equal([]string{
					"notifications.manage",
					"notifications.write",
					"emails.write",
					"notification_preferences.admin",
					"critical_notifications.write",
					"notification_templates.admin",
					"notification_templates.write",
					"notification_templates.read"}))
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
					"Authorization":           []string{authHeaderFixture},
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
					"Authorization":           []string{authHeaderFixture},
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
