package notificationtypes_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/v2/notificationtypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ryanmoran/stack"
)

var _ = Describe("ListHandler", func() {
	var (
		handler                     notificationtypes.ListHandler
		notificationTypesCollection *fakes.NotificationTypesCollection
		context                     stack.Context
		writer                      *httptest.ResponseRecorder
		request                     *http.Request
		database                    *fakes.Database
	)

	BeforeEach(func() {
		context = stack.NewContext()

		context.Set("client_id", "some-client-id")

		database = fakes.NewDatabase()
		context.Set("database", database)

		writer = httptest.NewRecorder()

		notificationTypesCollection = fakes.NewNotificationTypesCollection()

		handler = notificationtypes.NewListHandler(notificationTypesCollection)
	})

	It("returns a list of notification types", func() {
		notificationTypesCollection.ListCall.ReturnNotificationTypeList = []collections.NotificationType{
			{
				ID:          "notification-type-id-one",
				Name:        "first-notification-type",
				Description: "first-notification-type-description",
				Critical:    false,
				TemplateID:  "",
				SenderID:    "some-sender-id",
			},
			{
				ID:          "notification-type-id-two",
				Name:        "second-notification-type",
				Description: "second-notification-type-description",
				Critical:    true,
				TemplateID:  "",
				SenderID:    "some-sender-id",
			},
		}

		var err error
		request, err = http.NewRequest("GET", "/senders/some-sender-id/notification_types", nil)
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"notification_types": [
				{
					"id": "notification-type-id-one",
					"name": "first-notification-type",
					"description": "first-notification-type-description",
					"critical": false,
					"template_id": ""
				},
				{
					"id": "notification-type-id-two",
					"name": "second-notification-type",
					"description": "second-notification-type-description",
					"critical": true,
					"template_id": ""
				}
			]
		}`))
	})

	It("returns an empty list of notification types if the table has no records", func() {
		var err error
		request, err = http.NewRequest("GET", "/senders/some-sender-id/notification_types", nil)
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"notification_types": []
		}`))
	})

	Context("failure cases", func() {
		It("returns a 404 when the sender does not exist", func() {
			notificationTypesCollection.ListCall.Err = collections.NotFoundError{
				Err: errors.New("sender not found"),
			}

			var err error
			request, err = http.NewRequest("GET", "/senders/non-existent-sender-id/notification_types", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "sender not found"
			}`))
		})

		It("returns a 500 when the collection indicates a system error", func() {
			notificationTypesCollection.ListCall.Err = errors.New("BOOM!")

			var err error
			request, err = http.NewRequest("GET", "/senders/some-sender-id/notification_types", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "BOOM!"
			}`))
		})
	})
})
