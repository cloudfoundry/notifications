package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UnsetTemplates", func() {
	var handler handlers.UnsetTemplates
	var errorWriter *fakes.ErrorWriter
	var writer *httptest.ResponseRecorder
	var request *http.Request
	var context stack.Context
	var deleter *fakes.TemplateDeleter
	var err error

	Describe("ServeHTTP", func() {
		BeforeEach(func() {
			deleter = fakes.NewTemplateDeleter()
			errorWriter = fakes.NewErrorWriter()
			handler = handlers.NewUnsetTemplates(deleter, errorWriter)
			writer = httptest.NewRecorder()
			request, err = http.NewRequest("DELETE", "/templates/login.forgot."+models.UserBodyTemplateName, bytes.NewBuffer([]byte{}))
			if err != nil {
				panic(err)
			}
		})

		It("calls delete on the repo", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(deleter.DeleteArgument).To(Equal("login.forgot." + models.UserBodyTemplateName))
			Expect(writer.Code).To(Equal(http.StatusNoContent))
		})

		Context("When the deleter errors", func() {
			It("writes the error to the errorWriter", func() {
				deleter.DeleteError = errors.New("BOOM!")
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(Equal(deleter.DeleteError))
			})
		})
	})
})
