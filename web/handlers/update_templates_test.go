package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateTemplates", func() {
	var err error
	var handler handlers.UpdateTemplates
	var writer *httptest.ResponseRecorder
	var request *http.Request
	var context stack.Context
	var updater *fakes.TemplateUpdater
	var errorWriter *fakes.ErrorWriter

	Describe("ServeHTTP", func() {
		BeforeEach(func() {
			updater = fakes.NewTemplateUpdater()
			errorWriter = fakes.NewErrorWriter()
			handler = handlers.NewUpdateTemplates(updater, errorWriter)
			writer = httptest.NewRecorder()
			body := []byte(`{"name": "An Interesting Template", "subject": "very interesting", "text": "{{turkey}}", "html": "<p>{{turkey}} gobble</p>"}`)
			request, err = http.NewRequest("PUT", "/templates/a-template-id", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}
		})

		It("calls update on its updater with appropriate arguments", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(updater.UpdateArgumentID).To(Equal("a-template-id"))
			Expect(updater.UpdateArgumentBody).To(Equal(models.Template{
				Name:     "An Interesting Template",
				Subject:  "very interesting",
				Text:     "{{turkey}}",
				HTML:     "<p>{{turkey}} gobble</p>",
				Metadata: "{}",
			}))
			Expect(writer.Code).To(Equal(http.StatusNoContent))
		})

		It("can update a template without a subject field", func() {
			body := []byte(`{"name": "my template name", "html": "<p>gobble</p>", "text": "my awesome text"}`)
			request, err = http.NewRequest("PUT", "/templates/a-template-id.", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNoContent))
		})

		It("can update a template without a text field", func() {
			body := []byte(`{"name": "a temlate name", "subject": "my subject", "html": "<p>my html</p>"}`)
			request, err = http.NewRequest("PUT", "/templates/a-template-id", bytes.NewBuffer(body))
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNoContent))
		})

		Context("an error occurs", func() {
			Describe("when the name is missing from a template JSON body", func() {
				It("returns a validation error indicating the name is missing", func() {
					body := []byte(`{"subject": "my awesome subject", "html": "<p>gobble</p>", "text":"my awesome text"}`)
					request, err = http.NewRequest("PUT", "/templates/my-template-id", bytes.NewBuffer(body))
					if err != nil {
						panic(err)
					}

					handler.ServeHTTP(writer, request, context)
					Expect(errorWriter.Error).To(Equal(params.ValidationError([]string{
						"Request is missing the required field: name",
					})))
				})
			})

			Describe("when the html is missing from a template JSON body", func() {
				It("returns a validation error indicating the html is missing", func() {
					body := []byte(`{"subject": "my awesome subject", "name": "my awesome name", "text":"my awesome text"}`)
					request, err = http.NewRequest("PUT", "/templates/my-template-id", bytes.NewBuffer(body))
					if err != nil {
						panic(err)
					}

					handler.ServeHTTP(writer, request, context)
					Expect(errorWriter.Error).To(Equal(params.ValidationError([]string{
						"Request is missing the required field: html",
					})))
				})
			})

			Describe("when the update returns an error", func() {
				It("returns the error", func() {
					updater.UpdateError = models.TemplateUpdateError{Message: "My New Error"}
					body := []byte(`{"name": "a temlate name", "html": "<p>my html</p>"}`)
					request, err = http.NewRequest("PUT", "/templates/a-template-id", bytes.NewBuffer(body))
					if err != nil {
						panic(err)
					}

					handler.ServeHTTP(writer, request, context)
					Expect(errorWriter.Error).To(MatchError(models.TemplateUpdateError{Message: "My New Error"}))
				})
			})
		})
	})
})
