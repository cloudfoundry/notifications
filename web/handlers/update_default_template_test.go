package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateDefaultTemplate", func() {
	var err error
	var handler handlers.UpdateDefaultTemplate
	var writer *httptest.ResponseRecorder
	var request *http.Request
	var context stack.Context
	var updater *fakes.TemplateUpdater
	var errorWriter *fakes.ErrorWriter

	BeforeEach(func() {
		updater = fakes.NewTemplateUpdater()
		errorWriter = fakes.NewErrorWriter()
		handler = handlers.NewUpdateDefaultTemplate(updater, errorWriter)
		writer = httptest.NewRecorder()
	})

	It("updates the default template", func() {
		body := `{
			"name": "Defaultish Template",
			"subject": "{{.Subject}}",
			"html": "<p>something</p>",
			"text": "something",
			"metadata": {
				"hello": true
			}
		}`
		request, err = http.NewRequest("PUT", "/default_template", strings.NewReader(body))
		if err != nil {
			panic(err)
		}

		handler.ServeHTTP(writer, request, context)

		Expect(updater.UpdateArgumentID).To(Equal("default"))
		Expect(updater.UpdateArgumentBody).To(Equal(models.Template{
			Name:     "Defaultish Template",
			Subject:  "{{.Subject}}",
			HTML:     "<p>something</p>",
			Text:     "something",
			Metadata: `{"hello":true}`,
		}))
		Expect(writer.Code).To(Equal(http.StatusNoContent))
	})

	Context("when the request is not valid", func() {
		It("indicates that fields are missing", func() {
			body := `{
				"name": "Defaultish Template",
				"subject": "{{.Subject}}",
				"metadata": {}
			}`
			request, err = http.NewRequest("PUT", "/default_template", strings.NewReader(body))
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request, context)

			Expect(errorWriter.Error).To(MatchError(params.ValidationError([]string{
				"Request is missing the required field: html",
			})))
		})
	})
})
