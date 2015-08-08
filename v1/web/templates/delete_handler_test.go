package templates_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/web/templates"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteHandler", func() {
	var (
		handler     templates.DeleteHandler
		errorWriter *fakes.ErrorWriter
		writer      *httptest.ResponseRecorder
		request     *http.Request
		context     stack.Context
		deleter     *fakes.TemplateDeleter
		err         error
		database    *fakes.Database
	)

	Describe("ServeHTTP", func() {
		BeforeEach(func() {
			deleter = fakes.NewTemplateDeleter()
			errorWriter = fakes.NewErrorWriter()
			database = fakes.NewDatabase()
			writer = httptest.NewRecorder()
			request, err = http.NewRequest("DELETE", "/templates/template-id-123", bytes.NewBuffer([]byte{}))
			Expect(err).NotTo(HaveOccurred())

			context = stack.NewContext()
			context.Set("database", database)

			handler = templates.NewDeleteHandler(deleter, errorWriter)
		})

		It("calls delete on the repo", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(deleter.DeleteCall.Arguments).To(ConsistOf([]interface{}{database, "template-id-123"}))
			Expect(writer.Code).To(Equal(http.StatusNoContent))
		})

		Context("When the deleter errors", func() {
			It("writes the error to the errorWriter", func() {
				deleter.DeleteCall.Error = errors.New("BOOM!")
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(Equal(deleter.DeleteCall.Error))
			})
		})
	})
})
