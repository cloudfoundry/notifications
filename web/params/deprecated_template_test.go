package params_test

import (
	"bytes"
	"encoding/json"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/params"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeprecatedTemplate", func() {
	Describe("NewDeprecatedTemplate", func() {
		It("contructs parameters from a reader", func() {
			templateName := models.UserBodyTemplateName
			body, err := json.Marshal(map[string]interface{}{
				"text": `its foobar of course`,
				"html": `<p>its foobar</p>`,
			})
			if err != nil {
				panic(err)
			}

			parameters, err := params.NewDeprecatedTemplate(templateName, bytes.NewBuffer(body))

			Expect(parameters).To(BeAssignableToTypeOf(params.DeprecatedTemplate{}))
			Expect(parameters.Name).To(Equal(models.UserBodyTemplateName))
			Expect(parameters.Text).To(Equal("its foobar of course"))
			Expect(parameters.HTML).To(Equal("<p>its foobar</p>"))
		})
	})

	Describe("Validate", func() {
		Context("when the name is valid", func() {
			It("returns no error", func() {
				bad_endings := []string{models.UserBodyTemplateName, "my.silly." + models.SpaceBodyTemplateName, "this.special." + models.EmailBodyTemplateName, "emergency.email." + models.SubjectMissingTemplateName,
					models.SubjectProvidedTemplateName, "my.client." + models.UserBodyTemplateName, "client." + models.SpaceBodyTemplateName}

				for _, ending := range bad_endings {
					theTemplate := params.DeprecatedTemplate{
						Name: ending,
						Text: "its foobar of course",
						HTML: "<p>its foobar</p>",
					}
					err := theTemplate.Validate()
					Expect(err).ToNot(HaveOccurred())
				}
			})
		})

		Context("when the name is invalid", func() {
			It("returns an invalid name error", func() {
				bad_endings := []string{"user.body", "something_body", "subject.something", "still.missing.something",
					"client.kind.otherkind." + models.UserBodyTemplateName, "stupid.stuff.subject.uh.oh.damn." + models.EmailBodyTemplateName, "foo%.space_body"}

				for _, ending := range bad_endings {
					theTemplate := params.DeprecatedTemplate{
						Name: ending,
						Text: "its foobar of course",
						HTML: "<p>its foobar</p>",
					}
					err := theTemplate.Validate()
					Expect(err).To(HaveOccurred())
				}
			})
		})
	})

	Describe("ToModel", func() {
		It("turns a params.Template into a models.Template", func() {
			theTemplate := params.DeprecatedTemplate{
				Name: models.UserBodyTemplateName,
				Text: "its foobar of course",
				HTML: "<p>its foobar</p>",
			}
			theModel := theTemplate.ToModel()

			Expect(theModel).To(BeAssignableToTypeOf(models.Template{}))
			Expect(theModel.Name).To(Equal(models.UserBodyTemplateName))
			Expect(theModel.Text).To(Equal("its foobar of course"))
			Expect(theModel.HTML).To(Equal("<p>its foobar</p>"))
			Expect(theModel.CreatedAt).To(BeZero())
		})
	})
})
