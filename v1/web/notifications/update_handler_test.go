package notifications_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notifications"
	"github.com/cloudfoundry-incubator/notifications/web/webutil"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateHandler", func() {
	var (
		err         error
		handler     notifications.UpdateHandler
		writer      *httptest.ResponseRecorder
		request     *http.Request
		context     stack.Context
		updater     *fakes.NotificationUpdater
		errorWriter *fakes.ErrorWriter
		database    *fakes.Database
	)

	Describe("ServeHTTP", func() {
		BeforeEach(func() {
			updater = &fakes.NotificationUpdater{}
			errorWriter = fakes.NewErrorWriter()
			writer = httptest.NewRecorder()
			body := []byte(`{"description": "test kind", "critical": false, "template": "template-name"}`)
			request, err = http.NewRequest("PUT", "/clients/this-client/notifications/this-kind", bytes.NewBuffer(body))
			Expect(err).NotTo(HaveOccurred())

			database = fakes.NewDatabase()
			context = stack.NewContext()
			context.Set("database", database)

			handler = notifications.NewUpdateHandler(updater, errorWriter)
		})

		It("calls update on its updater with appropriate arguments", func() {
			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNoContent))

			Expect(updater.UpdateCall.Arguments).To(ConsistOf([]interface{}{database, models.Kind{
				Description: "test kind",
				Critical:    false,
				TemplateID:  "template-name",
				ClientID:    "this-client",
				ID:          "this-kind",
			}}))
		})

		Context("when an error occurs", func() {
			It("propagates the error returned from the updater into the error writer", func() {
				updater.UpdateCall.Error = errors.New("error occurred while updating notification")
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(MatchError(errors.New("error occurred while updating notification")))
			})

			It("writes a params validation error when the request is semantically invalid", func() {
				body := []byte(`{"description": "test kind", "template": "template-name"}`)
				request, err = http.NewRequest("PUT", "/clients/this-client/notifications/this-kind", bytes.NewBuffer(body))
				Expect(err).NotTo(HaveOccurred())

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(BeAssignableToTypeOf(webutil.ValidationError{}))
			})
		})
	})
})
