package params_test

import (
	"bytes"
	"encoding/json"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/params"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Template", func() {
	Describe("NewTemplate", func() {
		Context("When creating a new template", func() {
			It("contructs parameters from a reader", func() {
				body, err := json.Marshal(map[string]interface{}{
					"name":    `Foo Bar Baz`,
					"text":    `its foobar of course`,
					"html":    `<p>its foobar</p>`,
					"subject": `Stuff and Things`,
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
			})

			It("gracefully handles non-required missing parameters", func() {
				body, err := json.Marshal(map[string]interface{}{
					"name": `Foo Bar Baz`,
					"html": `<p>its foobar</p>`,
				})
				if err != nil {
					panic(err)
				}

				parameters, err := params.NewTemplate(bytes.NewBuffer(body))

				Expect(parameters).To(BeAssignableToTypeOf(params.Template{}))
				Expect(parameters.Name).To(Equal("Foo Bar Baz"))
				Expect(parameters.Text).To(Equal(""))
				Expect(parameters.HTML).To(Equal("<p>its foobar</p>"))
				Expect(parameters.Subject).To(Equal("{{.Subject}}"))
			})
		})
	})

	Describe("ToModel", func() {
		It("turns a params.Template into a models.Template", func() {
			theTemplate := params.Template{
				Name:    "The Foo to the Bar",
				Text:    "its foobar of course",
				HTML:    "<p>its foobar</p>",
				Subject: "Foobar Yah",
			}
			theModel := theTemplate.ToModel()

			Expect(theModel).To(BeAssignableToTypeOf(models.Template{}))
			Expect(theModel.Name).To(Equal("The Foo to the Bar"))
			Expect(theModel.Text).To(Equal("its foobar of course"))
			Expect(theModel.HTML).To(Equal("<p>its foobar</p>"))
			Expect(theModel.Overridden).To(BeTrue())
			Expect(theModel.CreatedAt).To(BeZero())
			Expect(theModel.UpdatedAt).To(BeZero())
		})
	})
})
