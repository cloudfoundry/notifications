package params_test

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/params"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func buildTemplateRequestBody(template params.Template) io.Reader {
	body, err := json.Marshal(template)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(body)
}

var _ = Describe("Template", func() {
	Describe("NewTemplate", func() {
		Context("When creating a new template", func() {
			It("contructs parameters from a reader", func() {
				body, err := json.Marshal(map[string]interface{}{
					"name":    `Foo Bar Baz`,
					"text":    `its foobar of course`,
					"html":    `<p>its foobar</p>`,
					"subject": `Stuff and Things`,
					"metadata": map[string]interface{}{
						"some_property": "some_value",
					},
				})
				if err != nil {
					panic(err)
				}

				parameters, err := params.NewTemplate(bytes.NewBuffer(body))

				Expect(parameters).To(BeAssignableToTypeOf(params.Template{}))
				Expect(parameters.Name).To(Equal("Foo Bar Baz"))
				Expect(parameters.Text).To(Equal("its foobar of course"))
				Expect(parameters.HTML).To(Equal("<p>its foobar</p>"))
				Expect(parameters.Subject).To(Equal("Stuff and Things"))
				Expect(string(parameters.Metadata)).To(MatchJSON(`{"some_property": "some_value"}`))
			})

			It("gracefully handles non-required missing parameters", func() {
				body, err := json.Marshal(map[string]interface{}{
					"name": "Foo Bar Baz",
					"html": "<p>its foobar</p>",
				})
				if err != nil {
					panic(err)
				}

				parameters, err := params.NewTemplate(bytes.NewBuffer(body))
				Expect(err).NotTo(HaveOccurred())
				Expect(parameters).To(BeAssignableToTypeOf(params.Template{}))
				Expect(parameters.Name).To(Equal("Foo Bar Baz"))
				Expect(parameters.Text).To(Equal(""))
				Expect(parameters.HTML).To(Equal("<p>its foobar</p>"))
				Expect(parameters.Subject).To(Equal("{{.Subject}}"))
				Expect(parameters.Metadata).To(Equal(json.RawMessage("{}")))
			})

			Context("when the template has invalid syntax", func() {
				Context("when subject template has invalid syntax", func() {
					It("returns a validation error", func() {
						body := buildTemplateRequestBody(params.Template{
							Name:    "Template name",
							Text:    "Textual template",
							HTML:    "HTML template",
							Subject: "{{.bad}",
						})
						_, err := params.NewTemplate(body)
						Expect(err).To(HaveOccurred())
						Expect(err).To(BeAssignableToTypeOf(params.ValidationError([]string{})))
					})
				})

				Context("when text template has invalid syntax", func() {
					It("returns a validation error", func() {
						body := buildTemplateRequestBody(params.Template{
							Name:    "Template name",
							Text:    "You should feel {{.BAD}",
							HTML:    "<h1> Amazing </h1>",
							Subject: "Great Subject",
						})
						_, err := params.NewTemplate(body)
						Expect(err).To(HaveOccurred())
						Expect(err).To(BeAssignableToTypeOf(params.ValidationError([]string{})))
					})
				})

				Context("when html template has invalid syntax", func() {
					It("returns a validation error", func() {
						body := buildTemplateRequestBody(params.Template{
							Name:    "Template name",
							Text:    "Textual template",
							HTML:    "{{.bad}",
							Subject: "Great Subject",
						})
						_, err := params.NewTemplate(body)
						Expect(err).To(HaveOccurred())
						Expect(err).To(BeAssignableToTypeOf(params.ValidationError([]string{})))
					})
				})
			})
		})
	})

	Describe("ToModel", func() {
		It("turns a params.Template into a models.Template", func() {
			theTemplate := params.Template{
				Name:     "The Foo to the Bar",
				Text:     "its foobar of course",
				HTML:     "<p>its foobar</p>",
				Subject:  "Foobar Yah",
				Metadata: json.RawMessage(`{"some_property": "some_value"}`),
			}
			theModel := theTemplate.ToModel()

			Expect(theModel).To(BeAssignableToTypeOf(models.Template{}))
			Expect(theModel.Name).To(Equal("The Foo to the Bar"))
			Expect(theModel.Text).To(Equal("its foobar of course"))
			Expect(theModel.HTML).To(Equal("<p>its foobar</p>"))
			Expect(theModel.Subject).To(Equal("Foobar Yah"))
			Expect(theModel.Metadata).To(MatchJSON(`{"some_property": "some_value"}`))
			Expect(theModel.CreatedAt).To(BeZero())
			Expect(theModel.UpdatedAt).To(BeZero())
		})
	})
})
