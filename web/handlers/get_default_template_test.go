package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetDefaultTemplate", func() {
	var handler handlers.GetDefaultTemplate
	var templateFinder *fakes.TemplateFinder
	var errorWriter *fakes.ErrorWriter

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

		handler = handlers.NewGetDefaultTemplate(templateFinder, errorWriter)
	})

	It("responds with a 200 status code and JSON representation of the template", func() {
		request, err := http.NewRequest("GET", "/default_template", nil)
		if err != nil {
			panic(err)
		}
		writer := httptest.NewRecorder()

		handler.ServeHTTP(writer, request, nil)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body).To(MatchJSON(`{
			"name": "Default Template",
			"subject": "CF Notification: {{.Subject}}",
			"text": "Default Template {{.Text}}",
			"html": "<p>Default Template</p> {{.HTML}}",
			"metadata": {}
		}`))
	})

	It("delegates error handling to the error writer", func() {
		templateFinder.FindError = errors.New("BANANA!!!")

		request, err := http.NewRequest("GET", "/default_template", nil)
		if err != nil {
			panic(err)
		}
		writer := httptest.NewRecorder()

		handler.ServeHTTP(writer, request, nil)

		Expect(errorWriter.Error).To(MatchError(errors.New("BANANA!!!")))
	})
})
