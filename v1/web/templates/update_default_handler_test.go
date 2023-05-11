package templates_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/web/templates"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/cloudfoundry-incubator/notifications/valiant"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateDefaultHandler", func() {
	var (
		handler     templates.UpdateDefaultHandler
		writer      *httptest.ResponseRecorder
		request     *http.Request
		context     stack.Context
		updater     *mocks.TemplateUpdater
		errorWriter *mocks.ErrorWriter
		database    *mocks.Database
	)

	BeforeEach(func() {
		var err error
		updater = mocks.NewTemplateUpdater()
		errorWriter = mocks.NewErrorWriter()
		writer = httptest.NewRecorder()
		request, err = http.NewRequest("PUT", "/default_template", strings.NewReader(`{
			"name": "Defaultish Template",
			"subject": "{{.Subject}}",
			"html": "<p>something</p>",
			"text": "something",
			"metadata": {"hello": true}
		}`))
		Expect(err).NotTo(HaveOccurred())

		database = mocks.NewDatabase()
		context = stack.NewContext()
		context.Set("database", database)

		handler = templates.NewUpdateDefaultHandler(updater, errorWriter)
	})

	It("updates the default template", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusNoContent))
		Expect(updater.UpdateCall.Receives.Database).To(Equal(database))
		Expect(updater.UpdateCall.Receives.TemplateID).To(Equal(models.DefaultTemplateID))
		Expect(updater.UpdateCall.Receives.Template).To(Equal(models.Template{
			Name:     "Defaultish Template",
			Subject:  "{{.Subject}}",
			HTML:     "<p>something</p>",
			Text:     "something",
			Metadata: `{"hello": true}`,
		}))
	})

	Context("when the request is not valid", func() {
		It("indicates that fields are missing", func() {
			body := `{
				"name": "Defaultish Template",
				"subject": "{{.Subject}}",
				"metadata": {}
			}`
			request, err := http.NewRequest("PUT", "/default_template", strings.NewReader(body))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)

			Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(webutil.ValidationError{Err: valiant.RequiredFieldError{ErrorMessage: "Missing required field 'html'"}}))
		})
	})

	Context("when the updater errors", func() {
		It("delegates the error handling to the error writer", func() {
			updater.UpdateCall.Returns.Error = errors.New("updating default template error")

			handler.ServeHTTP(writer, request, context)

			Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(errors.New("updating default template error")))
		})
	})
})
