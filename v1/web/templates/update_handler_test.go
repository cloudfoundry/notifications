package templates_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/web/templates"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/cloudfoundry-incubator/notifications/valiant"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateHandler", func() {
	var (
		err         error
		handler     templates.UpdateHandler
		writer      *httptest.ResponseRecorder
		request     *http.Request
		context     stack.Context
		updater     *mocks.TemplateUpdater
		errorWriter *mocks.ErrorWriter
		database    *mocks.Database
	)

	Describe("ServeHTTP", func() {
		BeforeEach(func() {
			updater = mocks.NewTemplateUpdater()
			errorWriter = mocks.NewErrorWriter()
			writer = httptest.NewRecorder()
			body := []byte(`{"name":"An Interesting Template", "subject":"very interesting subject", "text":"Here's the msg {{.Text}}", "html":"<p>turkey gobble</p>"}`)
			request, err = http.NewRequest("PUT", "/templates/a-template-id", bytes.NewBuffer(body))
			Expect(err).NotTo(HaveOccurred())

			database = mocks.NewDatabase()
			context = stack.NewContext()
			context.Set("database", database)

			handler = templates.NewUpdateHandler(updater, errorWriter)
		})

		It("calls update on its updater with appropriate arguments", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNoContent))

			Expect(updater.UpdateCall.Receives.Database).To(Equal(database))
			Expect(updater.UpdateCall.Receives.TemplateID).To(Equal("a-template-id"))
			Expect(updater.UpdateCall.Receives.Template).To(Equal(models.Template{
				Name:     "An Interesting Template",
				Subject:  "very interesting subject",
				Text:     "Here's the msg {{.Text}}",
				HTML:     "<p>turkey gobble</p>",
				Metadata: "{}",
			}))
		})

		It("can update a template without a subject field", func() {
			body := []byte(`{"name": "my template name", "html": "<p>gobble</p>", "text": "my awesome text"}`)
			request, err = http.NewRequest("PUT", "/templates/a-template-id.", bytes.NewBuffer(body))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNoContent))
		})

		It("can update a template without a text field", func() {
			body := []byte(`{"name": "a temlate name", "subject": "my subject", "html": "<p>my html</p>"}`)
			request, err = http.NewRequest("PUT", "/templates/a-template-id", bytes.NewBuffer(body))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNoContent))
		})

		Context("an error occurs", func() {
			Describe("when the name is missing from a template JSON body", func() {
				It("returns a validation error indicating the name is missing", func() {
					body := []byte(`{"subject": "my awesome subject", "html": "<p>gobble</p>", "text":"my awesome text"}`)
					request, err = http.NewRequest("PUT", "/templates/my-template-id", bytes.NewBuffer(body))
					Expect(err).NotTo(HaveOccurred())

					handler.ServeHTTP(writer, request, context)
					Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(webutil.ValidationError{Err: valiant.RequiredFieldError{ErrorMessage: "Missing required field 'name'"}}))
				})
			})

			Describe("when the html is missing from a template JSON body", func() {
				It("returns a validation error indicating the html is missing", func() {
					body := []byte(`{"subject": "my awesome subject", "name": "my awesome name", "text":"my awesome text"}`)
					request, err = http.NewRequest("PUT", "/templates/my-template-id", bytes.NewBuffer(body))
					Expect(err).NotTo(HaveOccurred())

					handler.ServeHTTP(writer, request, context)
					Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(webutil.ValidationError{Err: valiant.RequiredFieldError{ErrorMessage: "Missing required field 'html'"}}))
				})
			})

			Describe("when the update returns an error", func() {
				It("returns the error", func() {
					updater.UpdateCall.Returns.Error = models.TemplateUpdateError{Err: errors.New("some error")}
					body := []byte(`{"name": "a temlate name", "html": "<p>my html</p>"}`)
					request, err = http.NewRequest("PUT", "/templates/a-template-id", bytes.NewBuffer(body))
					Expect(err).NotTo(HaveOccurred())

					handler.ServeHTTP(writer, request, context)
					Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(models.TemplateUpdateError{Err: errors.New("some error")}))
				})
			})
		})
	})
})
