package templates_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/web/templates"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetDefaultHandler", func() {
	var (
		handler        templates.GetDefaultHandler
		templateFinder *mocks.TemplateFinder
		errorWriter    *mocks.ErrorWriter
		database       *mocks.Database
		context        stack.Context
	)

	BeforeEach(func() {
		errorWriter = mocks.NewErrorWriter()
		templateFinder = mocks.NewTemplateFinder()
		templateFinder.FindByIDCall.Returns.Template = models.Template{
			ID:       models.DefaultTemplateID,
			Name:     "Default Template",
			Subject:  "CF Notification: {{.Subject}}",
			Text:     "Default Template {{.Text}}",
			HTML:     "<p>Default Template</p> {{.HTML}}",
			Metadata: "{}",
		}

		database = mocks.NewDatabase()
		context = stack.NewContext()
		context.Set("database", database)

		handler = templates.NewGetDefaultHandler(templateFinder, errorWriter)
	})

	It("responds with a 200 status code and JSON representation of the template", func() {
		request, err := http.NewRequest("GET", "/default_template", nil)
		Expect(err).NotTo(HaveOccurred())
		writer := httptest.NewRecorder()

		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body).To(MatchJSON(`{
			"name": "Default Template",
			"subject": "CF Notification: {{.Subject}}",
			"text": "Default Template {{.Text}}",
			"html": "<p>Default Template</p> {{.HTML}}",
			"metadata": {}
		}`))

		Expect(templateFinder.FindByIDCall.Receives.Database).To(Equal(database))
		Expect(templateFinder.FindByIDCall.Receives.TemplateID).To(Equal(models.DefaultTemplateID))
	})

	It("delegates error handling to the error writer", func() {
		templateFinder.FindByIDCall.Returns.Error = errors.New("BANANA!!!")

		request, err := http.NewRequest("GET", "/default_template", nil)
		Expect(err).NotTo(HaveOccurred())
		writer := httptest.NewRecorder()

		handler.ServeHTTP(writer, request, context)

		Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(errors.New("BANANA!!!")))
	})
})
