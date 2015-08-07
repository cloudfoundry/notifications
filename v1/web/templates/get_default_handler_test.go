package templates_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/v1/web/templates"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetDefaultHandler", func() {
	var (
		handler        templates.GetDefaultHandler
		templateFinder *fakes.TemplateFinder
		errorWriter    *fakes.ErrorWriter
		database       *fakes.Database
		context        stack.Context
	)

	BeforeEach(func() {
		errorWriter = fakes.NewErrorWriter()
		templateFinder = fakes.NewTemplateFinder()
		templateFinder.Templates[models.DefaultTemplateID] = models.Template{
			ID:       models.DefaultTemplateID,
			Name:     "Default Template",
			Subject:  "CF Notification: {{.Subject}}",
			Text:     "Default Template {{.Text}}",
			HTML:     "<p>Default Template</p> {{.HTML}}",
			Metadata: "{}",
		}
		database = fakes.NewDatabase()
		context = stack.NewContext()
		context.Set("database", database)

		handler = templates.NewGetDefaultHandler(templateFinder, errorWriter)
	})

	It("responds with a 200 status code and JSON representation of the template", func() {
		request, err := http.NewRequest("GET", "/default_template", nil)
		if err != nil {
			panic(err)
		}
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

		Expect(templateFinder.FindByIDCall.Arguments).To(ConsistOf([]interface{}{database, models.DefaultTemplateID}))
	})

	It("delegates error handling to the error writer", func() {
		templateFinder.FindByIDCall.Error = errors.New("BANANA!!!")

		request, err := http.NewRequest("GET", "/default_template", nil)
		if err != nil {
			panic(err)
		}
		writer := httptest.NewRecorder()

		handler.ServeHTTP(writer, request, context)

		Expect(errorWriter.Error).To(MatchError(errors.New("BANANA!!!")))
	})
})
