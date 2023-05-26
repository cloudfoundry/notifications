package docs_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/docs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func toJSON(input interface{}) string {
	output, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(output)
}

const (
	tokenFixture      = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQzMjI1NTEsImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NjM3NzYvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.GbFa6KEgSarJYqU2yohnJj7G1ztoG4az8Ded1EcDVMbHf9OGKzuzXUKXPqruNzX3JhjwBlBN4f6P8yZXvhPmvnQSTG7a2lqE817lPs7L7c5J5ka6MBRWZqlJwquYOalv9ytn7ZpaodWNqGK_uH5ctpP9LFIoxG2WmKQsd1N7f_Wv8GyEbrHL6qMze8x1QhNuBos_zdOPuQ3DyMsJBlx7b35h2nY_ZXXDxVErGlALfrVfrUnAPbt7fgkRbX9JYtbzgnkhsrDXcPt_5xO3ideqYZARfrjKtCFQiR4kPqRADIGvjgaSDyX5fNvKpwmLfElmu_rCyr9tN7uZszShhwy9fQ"
	authHeaderFixture = "bearer " + tokenFixture
)

var _ = Describe("RoundTripRecorder", func() {
	Describe("Record", func() {
		Context("when a request is added to the document generator", func() {
			It("should add a new MethodEntry", func() {
				request := &http.Request{}
				response := &http.Response{}
				roundtripRecorder := docs.NewRoundTripRecorder()

				Expect(roundtripRecorder.Record("get-list-request", request, response)).To(Succeed())

				Expect(roundtripRecorder.RoundTrips).To(HaveKey("get-list-request"))

				roundtrip := roundtripRecorder.RoundTrips["get-list-request"]

				Expect(roundtrip).To(Equal(docs.RoundTrip{
					Request:  request,
					Response: response,
				}))
			})
		})

		Context("when two entries with the same endpoint key are added to the document generator", func() {
			It("return an error", func() {
				request := &http.Request{}
				response := &http.Response{}
				roundtripRecorder := docs.NewRoundTripRecorder()

				Expect(roundtripRecorder.Record("banana", request, response)).To(Succeed())
				err := roundtripRecorder.Record("banana", request, response)
				Expect(err).To(MatchError(errors.New("new roundtrip \"banana\" conflicts with existing roundtrip")))
			})
		})
	})

	Describe("BuildTemplateContext", func() {
		It("returns a context used to generate the API docs", func() {
			resources := []docs.Resource{
				{
					Name:        "Bananas",
					Description: "These are really tasty!",
					Endpoints: []docs.Endpoint{
						{
							Key:         "bananas-create",
							Description: "Create Bananas",
						},
						{
							Key:         "bananas-update",
							Description: "Update a Banana",
						},
					},
				},
			}
			roundtrips := map[string]docs.RoundTrip{
				"bananas-update": {
					Request: &http.Request{
						Method: "PUT",
						URL: &url.URL{
							Path: "/bananas/f2461184-871b-ef4b-e243-bc17bdcd3534",
						},
						Header: http.Header(map[string][]string{
							"X-Notifications-Version": {"2"},
						}),
						Body: ioutil.NopCloser(strings.NewReader(`{"color": "yellow","organic": false}`)),
					},
					Response: &http.Response{
						Status: "200 OK",
						Header: http.Header(map[string][]string{
							"Content-Length": {"123"},
							"Content-Type":   {"application/json"},
						}),
						Body: ioutil.NopCloser(strings.NewReader(`{"id": "banana-1","color": "yellow","organic": false}`)),
					},
				},
				"bananas-create": {
					Request: &http.Request{
						Method: "POST",
						URL: &url.URL{
							Path: "/bananas",
						},
						Header: http.Header(map[string][]string{
							"Authorization":           {"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwic2NvcGUiOlsiYmFuYW5hcy5ncm93Il0sImFkbWluIjp0cnVlfQ.ZaevrDILQ0OEmg5Baj8OSyNcJKoiNCfq_7yaYExTJA0"},
							"X-Notifications-Version": {"2"},
						}),
						Body: ioutil.NopCloser(strings.NewReader(`{"color": "green","organic": true}`)),
					},
					Response: &http.Response{
						Status: "201 Created",
						Header: http.Header(map[string][]string{
							"Content-Length": {"123"},
							"Content-Type":   {"application/json"},
						}),
						Body: ioutil.NopCloser(strings.NewReader(`{"id": "banana-1","color": "green","organic": true}`)),
					},
				},
			}

			context, err := docs.BuildTemplateContext(resources, roundtrips)
			Expect(err).NotTo(HaveOccurred())
			Expect(context).To(Equal(docs.TemplateContext{
				Resources: []docs.TemplateResource{
					{
						Name:        "Bananas",
						Description: "These are really tasty!",
						Endpoints: []docs.TemplateEndpoint{
							{
								Key:            "bananas-create",
								Description:    "Create Bananas",
								Method:         "POST",
								Path:           "/bananas",
								RequiredScopes: "bananas.grow",
								RequestHeaders: []string{
									"Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwic2NvcGUiOlsiYmFuYW5hcy5ncm93Il0sImFkbWluIjp0cnVlfQ.ZaevrDILQ0OEmg5Baj8OSyNcJKoiNCfq_7yaYExTJA0",
									"X-Notifications-Version: 2",
								},
								RequestBody: toJSON(map[string]interface{}{
									"organic": true,
									"color":   "green",
								}),
								ResponseStatus: "201 Created",
								ResponseHeaders: []string{
									"Content-Length: 123",
									"Content-Type: application/json",
								},
								ResponseBody: toJSON(map[string]interface{}{
									"id":      "banana-1",
									"organic": true,
									"color":   "green",
								}),
							},
							{
								Key:         "bananas-update",
								Description: "Update a Banana",
								Method:      "PUT",
								Path:        "/bananas/{id}",
								RequestHeaders: []string{
									"X-Notifications-Version: 2",
								},
								RequestBody: toJSON(map[string]interface{}{
									"organic": false,
									"color":   "yellow",
								}),
								ResponseStatus: "200 OK",
								ResponseHeaders: []string{
									"Content-Length: 123",
									"Content-Type: application/json",
								},
								ResponseBody: toJSON(map[string]interface{}{
									"id":      "banana-1",
									"organic": false,
									"color":   "yellow",
								}),
							},
						},
					},
				},
			}))
		})

		It("returns an error when there is a roundtrip that is not used", func() {
			resources := []docs.Resource{}
			roundtrips := map[string]docs.RoundTrip{
				"bananas-list": {
					Request: &http.Request{
						Method: "GET",
					},
				},
				"bananas-create": {
					Request: &http.Request{
						Method: "POST",
					},
				},
			}

			_, err := docs.BuildTemplateContext(resources, roundtrips)
			Expect(err).To(MatchError(errors.New("unused roundtrips [bananas-create bananas-list]")))
		})

		It("returns an error when a roundtrip is missing", func() {
			resources := []docs.Resource{
				{
					Name: "Bananas",
					Endpoints: []docs.Endpoint{
						{
							Key:         "bananas-create",
							Description: "Makin' naners",
						},
					},
				},
			}
			roundtrips := map[string]docs.RoundTrip{}

			_, err := docs.BuildTemplateContext(resources, roundtrips)
			Expect(err).To(MatchError(errors.New("missing roundtrip \"bananas-create\"")))
		})
	})

	Describe("GenerateMarkdown", func() {
		It("generates a markdown representation of the API docs", func() {
			context := docs.TemplateContext{
				Resources: []docs.TemplateResource{
					{
						Name:        "Bananas",
						Description: "These are really tasty!",
						Endpoints: []docs.TemplateEndpoint{
							{
								Key:            "bananas-create",
								Description:    "Create Bananas",
								Method:         "POST",
								Path:           "/bananas",
								RequiredScopes: "bananas.grow",
								RequestHeaders: []string{
									"Authorization: Bearer some-token",
									"X-Notifications-Version: 2",
								},
								RequestBody: toJSON(map[string]bool{
									"organic": true,
								}),
								ResponseStatus: "201 Created",
								ResponseHeaders: []string{
									"Content-Length: 123",
									"Content-Type: application/json",
								},
								ResponseBody: toJSON(map[string]interface{}{
									"id":      "banana-1",
									"organic": true,
								}),
							},
							{
								Key:         "bananas-update",
								Description: "Update a Banana",
								Method:      "PUT",
								Path:        "/bananas/{id}",
								RequestHeaders: []string{
									"Authorization: Bearer some-token",
									"X-Notifications-Version: 2",
								},
								RequestBody: toJSON(map[string]bool{
									"organic": false,
								}),
								ResponseStatus: "200 OK",
								ResponseHeaders: []string{
									"Content-Length: 123",
									"Content-Type: application/json",
								},
								ResponseBody: toJSON(map[string]interface{}{
									"id":      "banana-1",
									"organic": false,
								}),
							},
						},
					},
					{
						Name:        "Bicycles",
						Description: "Fun when they go fast.",
						Endpoints: []docs.TemplateEndpoint{
							{
								Key:            "bicycles-create",
								Description:    "Create Bicycles",
								Method:         "POST",
								Path:           "/bicycles",
								RequiredScopes: "bicycles.manufacture",
								RequestHeaders: []string{
									"Authorization: Bearer some-token",
									"X-Notifications-Version: 2",
								},
								RequestBody: toJSON(map[string]interface{}{
									"color": "blue",
									"size":  46,
								}),
								ResponseStatus: "201 Created",
								ResponseHeaders: []string{
									"Content-Length: 123",
									"Content-Type: application/json",
								},
								ResponseBody: toJSON(map[string]interface{}{
									"id":    "bicycle-15",
									"color": "blue",
									"size":  46,
								}),
							},
							{
								Key:            "bicycles-delete",
								Description:    "Delete a Bicycle",
								Method:         "DELETE",
								Path:           "/bicycles/{id}",
								RequiredScopes: "bicycles.crash",
								RequestHeaders: []string{
									"Authorization: Bearer some-token",
									"X-Notifications-Version: 2",
								},
								ResponseStatus: "204 No Content",
								ResponseHeaders: []string{
									"Date: today",
								},
							},
						},
					},
				},
			}

			md, err := docs.GenerateMarkdown(context)
			Expect(err).NotTo(HaveOccurred())

			fixture, err := ioutil.ReadFile("../testing/fixtures/docs.md")
			Expect(err).NotTo(HaveOccurred())
			Expect(strings.TrimSpace(md)).To(Equal(strings.TrimSpace(string(fixture))))
		})
	})
})
