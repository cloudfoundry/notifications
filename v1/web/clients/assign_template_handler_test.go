package clients_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/web/clients"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AssignTemplateHandler", func() {
	var (
		handler          clients.AssignTemplateHandler
		templateAssigner *mocks.TemplateAssigner
		errorWriter      *mocks.ErrorWriter
		context          stack.Context
		database         *mocks.Database
		connection       *mocks.Connection
	)

	BeforeEach(func() {
		templateAssigner = mocks.NewTemplateAssigner()
		errorWriter = mocks.NewErrorWriter()
		connection = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = connection
		context = stack.NewContext()
		context.Set("database", database)

		handler = clients.NewAssignTemplateHandler(templateAssigner, errorWriter)
	})

	It("associates a template with a client", func() {
		body, err := json.Marshal(map[string]string{
			"template": "my-template",
		})
		Expect(err).NotTo(HaveOccurred())

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/template", bytes.NewBuffer(body))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(w, request, context)

		Expect(w.Code).To(Equal(http.StatusNoContent))
		Expect(templateAssigner.AssignToClientCall.Receives.Connection).To(Equal(connection))
		Expect(templateAssigner.AssignToClientCall.Receives.ClientID).To(Equal("my-client"))
		Expect(templateAssigner.AssignToClientCall.Receives.TemplateID).To(Equal("my-template"))
	})

	It("delegates to the error writer when the assigner errors", func() {
		templateAssigner.AssignToClientCall.Returns.Error = errors.New("banana")
		body, err := json.Marshal(map[string]string{
			"template": "my-template",
		})
		Expect(err).NotTo(HaveOccurred())

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/template", bytes.NewBuffer(body))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(w, request, context)
		Expect(errorWriter.WriteCall.Receives.Error).To(Equal(errors.New("banana")))
	})

	It("writes a ParseError to the error writer when request body is invalid", func() {
		body := []byte(`{ "this is" : not-valid-json }`)

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/template", bytes.NewBuffer(body))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(w, request, context)
		Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.ParseError{}))
	})
})
