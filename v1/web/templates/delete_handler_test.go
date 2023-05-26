package templates_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/web/templates"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteHandler", func() {
	var (
		handler     templates.DeleteHandler
		errorWriter *mocks.ErrorWriter
		writer      *httptest.ResponseRecorder
		request     *http.Request
		context     stack.Context
		deleter     *mocks.TemplateDeleter
		err         error
		connection  *mocks.Connection
	)

	Describe("ServeHTTP", func() {
		BeforeEach(func() {
			deleter = mocks.NewTemplateDeleter()
			errorWriter = mocks.NewErrorWriter()
			writer = httptest.NewRecorder()
			request, err = http.NewRequest("DELETE", "/templates/template-id-123", bytes.NewBuffer([]byte{}))
			Expect(err).NotTo(HaveOccurred())

			connection = mocks.NewConnection()
			database := mocks.NewDatabase()
			database.ConnectionCall.Returns.Connection = connection

			context = stack.NewContext()
			context.Set("database", database)

			handler = templates.NewDeleteHandler(deleter, errorWriter)
		})

		It("calls delete on the repo", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNoContent))

			Expect(deleter.DeleteCall.Receives.Connection).To(Equal(connection))
			Expect(deleter.DeleteCall.Receives.TemplateID).To(Equal("template-id-123"))
		})

		Context("When the deleter errors", func() {
			It("writes the error to the errorWriter", func() {
				deleter.DeleteCall.Returns.Error = errors.New("BOOM!")
				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.WriteCall.Receives.Error).To(Equal(errors.New("BOOM!")))
			})
		})
	})
})
