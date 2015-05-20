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
	var (
		handler     handlers.UpdateDefaultTemplate
		writer      *httptest.ResponseRecorder
		request     *http.Request
		context     stack.Context
		updater     *fakes.TemplateUpdater
		errorWriter *fakes.ErrorWriter
		database    *fakes.Database
	)

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
		Expect(err).NotTo(HaveOccurred())

		database = fakes.NewDatabase()
		context = stack.NewContext()
		context.Set("database", database)
	})

	It("updates the default template", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusNoContent))
		Expect(updater.UpdateCall.Arguments).To(ConsistOf([]interface{}{database, models.DefaultTemplateID, models.Template{
			Name:     "Defaultish Template",
			Subject:  "{{.Subject}}",
			HTML:     "<p>something</p>",
			Text:     "something",
			Metadata: `{"hello": true}`,
		}}))
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

			Expect(errorWriter.Error).To(BeAssignableToTypeOf(params.ValidationError([]string{})))
		})
	})

	Context("when the updater errors", func() {
		It("delegates the error handling to the error writer", func() {
			updater.UpdateCall.Error = errors.New("updating default template error")

			handler.ServeHTTP(writer, request, context)

			Expect(errorWriter.Error).To(MatchError(errors.New("updating default template error")))
		})
	})
})
