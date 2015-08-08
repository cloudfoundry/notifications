package notifications_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notifications"
	"github.com/cloudfoundry-incubator/notifications/web/webutil"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AssignTemplateHandler", func() {
	var (
		handler          notifications.AssignTemplateHandler
		templateAssigner *fakes.TemplateAssigner
		errorWriter      *fakes.ErrorWriter
		context          stack.Context
		database         *fakes.Database
	)

	BeforeEach(func() {
		templateAssigner = fakes.NewTemplateAssigner()
		errorWriter = fakes.NewErrorWriter()
		database = fakes.NewDatabase()
		context = stack.NewContext()
		context.Set("database", database)

		handler = notifications.NewAssignTemplateHandler(templateAssigner, errorWriter)
	})

	It("associates a template with a notification", func() {
		body, err := json.Marshal(map[string]string{
			"template": "my-template",
		})
		if err != nil {
			panic(err)
		}

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/notifications/my-notification/template", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}

		handler.ServeHTTP(w, request, context)

		Expect(w.Code).To(Equal(http.StatusNoContent))
		Expect(templateAssigner.AssignToNotificationArguments).To(Equal([]interface{}{database, "my-client", "my-notification", "my-template"}))
	})

	It("delegates to the error writer when the assigner errors", func() {
		templateAssigner.AssignToNotificationError = errors.New("banana")
		body, err := json.Marshal(map[string]string{
			"template": "my-template",
		})
		if err != nil {
			panic(err)
		}

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/notifications/my-notification/template", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}

		handler.ServeHTTP(w, request, context)
		Expect(errorWriter.WriteCall.Receives.Error).To(Equal(errors.New("banana")))
	})

	It("writes a ParseError to the error writer when request body is invalid", func() {
		body := []byte(`{ "this is" : not-valid-json }`)

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/notifications/my-notification/template", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}

		handler.ServeHTTP(w, request, context)
		Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.ParseError{}))
	})
})
