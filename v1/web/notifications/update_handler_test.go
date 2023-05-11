package notifications_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notifications"
	"github.com/cloudfoundry-incubator/notifications/v1/web/webutil"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateHandler", func() {
	var (
		err         error
		handler     notifications.UpdateHandler
		writer      *httptest.ResponseRecorder
		request     *http.Request
		context     stack.Context
		updater     *mocks.NotificationUpdater
		errorWriter *mocks.ErrorWriter
		database    *mocks.Database
	)

	Describe("ServeHTTP", func() {
		BeforeEach(func() {
			updater = &mocks.NotificationUpdater{}
			errorWriter = mocks.NewErrorWriter()
			writer = httptest.NewRecorder()
			body := []byte(`{"description": "test kind", "critical": false, "template": "template-name"}`)
			request, err = http.NewRequest("PUT", "/clients/this-client/notifications/this-kind", bytes.NewBuffer(body))
			Expect(err).NotTo(HaveOccurred())

			database = mocks.NewDatabase()
			context = stack.NewContext()
			context.Set("database", database)

			handler = notifications.NewUpdateHandler(updater, errorWriter)
		})

		It("calls update on its updater with appropriate arguments", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNoContent))

			Expect(updater.UpdateCall.Receives.Database).To(Equal(database))
			Expect(updater.UpdateCall.Receives.Notification).To(Equal(models.Kind{
				Description: "test kind",
				Critical:    false,
				TemplateID:  "template-name",
				ClientID:    "this-client",
				ID:          "this-kind",
			}))
		})

		Context("when an error occurs", func() {
			It("propagates the error returned from the updater into the error writer", func() {
				updater.UpdateCall.Returns.Error = errors.New("error occurred while updating notification")
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(errors.New("error occurred while updating notification")))
			})

			It("writes a params validation error when the request is semantically invalid", func() {
				body := []byte(`{"description": "test kind", "template": "template-name"}`)
				request, err = http.NewRequest("PUT", "/clients/this-client/notifications/this-kind", bytes.NewBuffer(body))
				Expect(err).NotTo(HaveOccurred())

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(BeAssignableToTypeOf(webutil.ValidationError{}))
			})
		})
	})
})
