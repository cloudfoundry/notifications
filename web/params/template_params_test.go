package params_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/cloudfoundry-incubator/notifications/web/params"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func buildTemplateRequestBody(template params.TemplateParams) io.Reader {
	body, err := json.Marshal(template)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(body)
}

var _ = Describe("TemplateParams", func() {
	Describe("NewTemplateParams", func() {
		Context("when creating a new template", func() {
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
				Expect(err).NotTo(HaveOccurred())

				parameters, err := params.NewTemplateParams(ioutil.NopCloser(bytes.NewBuffer(body)))
				Expect(err).NotTo(HaveOccurred())
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
				Expect(err).NotTo(HaveOccurred())

				parameters, err := params.NewTemplateParams(ioutil.NopCloser(bytes.NewBuffer(body)))
				Expect(err).NotTo(HaveOccurred())
				Expect(parameters.Name).To(Equal("Foo Bar Baz"))
				Expect(parameters.Text).To(Equal(""))
				Expect(parameters.HTML).To(Equal("<p>its foobar</p>"))
				Expect(parameters.Subject).To(Equal("{{.Subject}}"))
				Expect(parameters.Metadata).To(Equal(json.RawMessage("{}")))
			})

			Context("when the template has invalid syntax", func() {
				Context("when subject template has invalid syntax", func() {
					It("returns a validation error", func() {
						body := buildTemplateRequestBody(params.TemplateParams{
							Name:    "Template name",
							Text:    "Textual template",
							HTML:    "HTML template",
							Subject: "{{.bad}",
						})
						_, err := params.NewTemplateParams(ioutil.NopCloser(body))
						Expect(err).To(HaveOccurred())
						Expect(err).To(BeAssignableToTypeOf(params.ValidationError([]string{})))
					})
				})

				Context("when text template has invalid syntax", func() {
					It("returns a validation error", func() {
						body := buildTemplateRequestBody(params.TemplateParams{
							Name:    "Template name",
							Text:    "You should feel {{.BAD}",
							HTML:    "<h1> Amazing </h1>",
							Subject: "Great Subject",
						})
						_, err := params.NewTemplateParams(ioutil.NopCloser(body))
						Expect(err).To(HaveOccurred())
						Expect(err).To(BeAssignableToTypeOf(params.ValidationError([]string{})))
					})
				})

				Context("when html template has invalid syntax", func() {
					It("returns a validation error", func() {
						body := buildTemplateRequestBody(params.TemplateParams{
							Name:    "Template name",
							Text:    "Textual template",
							HTML:    "{{.bad}",
							Subject: "Great Subject",
						})
						_, err := params.NewTemplateParams(ioutil.NopCloser(body))
						Expect(err).To(HaveOccurred())
						Expect(err).To(BeAssignableToTypeOf(params.ValidationError([]string{})))
					})
				})
			})
		})
	})

	Describe("ToModel", func() {
		It("turns a params.Template into a models.Template", func() {
			templateParams := params.TemplateParams{
				Name:     "The Foo to the Bar",
				Text:     "its foobar of course",
				HTML:     "<p>its foobar</p>",
				Subject:  "Foobar Yah",
				Metadata: json.RawMessage(`{"some_property": "some_value"}`),
			}
			templateModel := templateParams.ToModel()

			Expect(templateModel.Name).To(Equal("The Foo to the Bar"))
			Expect(templateModel.Text).To(Equal("its foobar of course"))
			Expect(templateModel.HTML).To(Equal("<p>its foobar</p>"))
			Expect(templateModel.Subject).To(Equal("Foobar Yah"))
			Expect(templateModel.Metadata).To(MatchJSON(`{"some_property": "some_value"}`))
			Expect(templateModel.CreatedAt).To(BeZero())
			Expect(templateModel.UpdatedAt).To(BeZero())
		})
	})
})
