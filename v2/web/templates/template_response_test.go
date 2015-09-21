package templates_test

import (
	"encoding/json"

	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/templates"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplateResponse", func() {
	It("provides a JSON representation of a template resource", func() {
		template := collections.Template{
			ID:       "some-template-id",
			Name:     "some-template",
			Text:     "template-text",
			HTML:     "template-html",
			Subject:  "template-subject",
			Metadata: `{ "template": "metadata" }`,
		}

		metadata := json.RawMessage(template.Metadata)
		response := templates.NewTemplateResponse(template)
		Expect(response).To(Equal(templates.TemplateResponse{
			ID:       "some-template-id",
			Name:     "some-template",
			Text:     "template-text",
			HTML:     "template-html",
			Subject:  "template-subject",
			Metadata: &metadata,
			Links: templates.TemplateResponseLinks{
				Self: templates.Link{"/templates/some-template-id"},
			},
		}))
	})

	It("can marshal into JSON", func() {
		template := collections.Template{
			ID:       "some-template-id",
			Name:     "some-template",
			Text:     "template-text",
			HTML:     "template-html",
			Subject:  "template-subject",
			Metadata: `{ "template": "metadata" }`,
		}

		output, err := json.Marshal(templates.NewTemplateResponse(template))
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(MatchJSON(`{
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
		}`))
	})
})
