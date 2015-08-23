package notifications_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notifications"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AssignTemplateHandler", func() {
	var (
		handler          notifications.AssignTemplateHandler
		templateAssigner *mocks.TemplateAssigner
		errorWriter      *mocks.ErrorWriter
		context          stack.Context
		database         *mocks.Database
	)

	BeforeEach(func() {
		templateAssigner = mocks.NewTemplateAssigner()
		errorWriter = mocks.NewErrorWriter()
		database = mocks.NewDatabase()
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
		Expect(templateAssigner.AssignToNotificationCall.Receives.Database).To(Equal(database))
		Expect(templateAssigner.AssignToNotificationCall.Receives.ClientID).To(Equal("my-client"))
		Expect(templateAssigner.AssignToNotificationCall.Receives.NotificationID).To(Equal("my-notification"))
		Expect(templateAssigner.AssignToNotificationCall.Receives.TemplateID).To(Equal("my-template"))
	})

	It("delegates to the error writer when the assigner errors", func() {
		templateAssigner.AssignToNotificationCall.Returns.Error = errors.New("banana")
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
