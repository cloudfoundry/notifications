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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AssignTemplateHandler", func() {
	var (
		handler          notifications.AssignTemplateHandler
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

		handler = notifications.NewAssignTemplateHandler(templateAssigner, errorWriter)
	})

	It("associates a template with a notification", func() {
		body, err := json.Marshal(map[string]string{
			"template": "my-template",
		})
		Expect(err).NotTo(HaveOccurred())

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/notifications/my-notification/template", bytes.NewBuffer(body))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(w, request, context)

		Expect(w.Code).To(Equal(http.StatusNoContent))
		Expect(templateAssigner.AssignToNotificationCall.Receives.Connection).To(Equal(connection))
		Expect(templateAssigner.AssignToNotificationCall.Receives.ClientID).To(Equal("my-client"))
		Expect(templateAssigner.AssignToNotificationCall.Receives.NotificationID).To(Equal("my-notification"))
		Expect(templateAssigner.AssignToNotificationCall.Receives.TemplateID).To(Equal("my-template"))
	})

	It("delegates to the error writer when the assigner errors", func() {
		templateAssigner.AssignToNotificationCall.Returns.Error = errors.New("banana")
		body, err := json.Marshal(map[string]string{
			"template": "my-template",
		})
		Expect(err).NotTo(HaveOccurred())

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/notifications/my-notification/template", bytes.NewBuffer(body))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(w, request, context)
		Expect(errorWriter.WriteCall.Receives.Error).To(Equal(errors.New("banana")))
	})

	It("writes a ParseError to the error writer when request body is invalid", func() {
		body := []byte(`{ "this is" : not-valid-json }`)

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/notifications/my-notification/template", bytes.NewBuffer(body))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(w, request, context)
		Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.ParseError{}))
	})
})
