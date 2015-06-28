package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AssignClientTemplate", func() {
	var (
		handler          handlers.AssignClientTemplate
		templateAssigner *fakes.TemplateAssigner
		errorWriter      *fakes.ErrorWriter
		context          stack.Context
		database         *fakes.Database
	)

	BeforeEach(func() {
		templateAssigner = fakes.NewTemplateAssigner()
		errorWriter = fakes.NewErrorWriter()
		handler = handlers.NewAssignClientTemplate(templateAssigner, errorWriter)
		database = fakes.NewDatabase()
		context = stack.NewContext()
		context.Set("database", database)
	})

	It("associates a template with a client", func() {
		body, err := json.Marshal(map[string]string{
			"template": "my-template",
		})
		if err != nil {
			panic(err)
		}

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/template", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}

		handler.ServeHTTP(w, request, context)

		Expect(w.Code).To(Equal(http.StatusNoContent))

		Expect(templateAssigner.AssignToClientArguments).To(Equal([]interface{}{database, "my-client", "my-template"}))
	})

	It("delegates to the error writer when the assigner errors", func() {
		templateAssigner.AssignToClientError = errors.New("banana")
		body, err := json.Marshal(map[string]string{
			"template": "my-template",
		})
		if err != nil {
			panic(err)
		}

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/template", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}

		handler.ServeHTTP(w, request, context)
		Expect(errorWriter.Error).To(Equal(errors.New("banana")))
	})

	It("writes a ParseError to the error writer when request body is invalid", func() {
		body := []byte(`{ "this is" : not-valid-json }`)

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/template", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}

		handler.ServeHTTP(w, request, context)
		Expect(errorWriter.Error).To(BeAssignableToTypeOf(handlers.ParseError{}))
	})
})
