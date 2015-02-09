package handlers_test

import (
	"errors"
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
	var handler handlers.UpdateDefaultTemplate
	var writer *httptest.ResponseRecorder
	var request *http.Request
	var context stack.Context
	var updater *fakes.TemplateUpdater
	var errorWriter *fakes.ErrorWriter

	BeforeEach(func() {
		var err error
		updater = fakes.NewTemplateUpdater()
		errorWriter = fakes.NewErrorWriter()
		handler = handlers.NewUpdateDefaultTemplate(updater, errorWriter)
		writer = httptest.NewRecorder()
		request, err = http.NewRequest("PUT", "/default_template", strings.NewReader(`{
			"name": "Defaultish Template",
			"subject": "{{.Subject}}",
			"html": "<p>something</p>",
			"text": "something",
			"metadata": {"hello": true}
		}`))
		if err != nil {
			panic(err)
		}
	})

	It("updates the default template", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(updater.UpdateArgumentID).To(Equal(models.DefaultTemplateID))
		Expect(updater.UpdateArgumentBody).To(Equal(models.Template{
			Name:     "Defaultish Template",
			Subject:  "{{.Subject}}",
			HTML:     "<p>something</p>",
			Text:     "something",
			Metadata: `{"hello": true}`,
		}))
		Expect(writer.Code).To(Equal(http.StatusNoContent))
	})

	Context("when the request is not valid", func() {
		It("indicates that fields are missing", func() {
			var err error
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

			Expect(errorWriter.Error).To(BeAssignableToTypeOf(params.ValidationError([]string{})))
		})
	})

	Context("when the updater errors", func() {
		It("delegates the error handling to the error writer", func() {
			updater.UpdateError = errors.New("updating default template error")

			handler.ServeHTTP(writer, request, context)

			Expect(errorWriter.Error).To(MatchError(errors.New("updating default template error")))
		})
	})
})
