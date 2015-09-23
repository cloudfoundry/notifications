package templates_test

import (
	"encoding/json"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/templates"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplatesListResponse", func() {
	var response templates.TemplatesListResponse

	BeforeEach(func() {
		response = templates.NewTemplatesListResponse([]collections.Template{
			{
				ID:       "some-template-id",
				Name:     "some-template",
				Text:     "template-text",
				HTML:     "template-html",
				Subject:  "template-subject",
				Metadata: `{ "template": "metadata" }`,
			},
			{
				ID:       "another-template-id",
				Name:     "another-template",
				Text:     "another-template-text",
				HTML:     "another-template-html",
				Subject:  "another-template-subject",
				Metadata: `{ "template": "another-metadata" }`,
			},
		})
	})

	It("provides a JSON representation of a list of template resources", func() {
		metadata1 := json.RawMessage(`{ "template": "metadata" }`)
		metadata2 := json.RawMessage(`{ "template": "another-metadata" }`)

		Expect(response).To(Equal(templates.TemplatesListResponse{
			Templates: []templates.TemplateResponse{
				{
					ID:       "some-template-id",
					Name:     "some-template",
					Text:     "template-text",
					HTML:     "template-html",
					Subject:  "template-subject",
					Metadata: &metadata1,
					Links: templates.TemplateResponseLinks{
						Self: templates.Link{"/templates/some-template-id"},
					},
				},
				{
					ID:       "another-template-id",
					Name:     "another-template",
					Text:     "another-template-text",
					HTML:     "another-template-html",
					Subject:  "another-template-subject",
					Metadata: &metadata2,
					Links: templates.TemplateResponseLinks{
						Self: templates.Link{"/templates/another-template-id"},
					},
				},
			},
			Links: templates.TemplatesListResponseLinks{
				Self: templates.Link{"/templates"},
			},
		}))
	})

	It("can marshal into JSON", func() {
		output, err := json.Marshal(response)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(MatchJSON(`{
			"templates": [
				{
					"id": "some-template-id",
					"name": "some-template",
					"text": "template-text",
					"html": "template-html",
					"subject": "template-subject",
					"metadata": {
						"template": "metadata"
					},
					"_links": {
						"self": {
							"href": "/templates/some-template-id"
						}
					}
				},
				{
					"id": "another-template-id",
					"name": "another-template",
					"text": "another-template-text",
					"html": "another-template-html",
					"subject": "another-template-subject",
					"metadata": {
						"template": "another-metadata"
					},
					"_links": {
						"self": {
							"href": "/templates/another-template-id"
						}
					}
				}
			],
			"_links": {
				"self": {
					"href": "/templates"
				}
			}
		}`))
	})

	Context("when the list is empty", func() {
		It("returns an empty list (not null)", func() {
			response = templates.NewTemplatesListResponse([]collections.Template{})

			output, err := json.Marshal(response)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(MatchJSON(`{
				"templates": [],
				"_links": {
					"self": {
						"href": "/templates"
					}
				}
			}`))
		})
	})
})
