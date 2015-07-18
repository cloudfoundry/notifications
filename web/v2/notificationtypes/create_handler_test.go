package notificationtypes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/v2/notificationtypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ryanmoran/stack"
)

var _ = Describe("CreateHandler", func() {
	var (
		handler                     notificationtypes.CreateHandler
		notificationTypesCollection *fakes.NotificationTypesCollection
		context                     stack.Context
		writer                      *httptest.ResponseRecorder
		request                     *http.Request
		database                    *fakes.Database
	)

	BeforeEach(func() {
		database = fakes.NewDatabase()
		context = stack.NewContext()
		context.Set("database", database)

		writer = httptest.NewRecorder()
		notificationTypesCollection = fakes.NewNotificationTypesCollection()
		notificationTypesCollection.AddCall.ReturnNotificationType = collections.NotificationType{
			ID:          "some-notification-type-id",
			Name:        "some-notification-type",
			Description: "some-notification-type-description",
			Critical:    false,
			TemplateID:  "some-template-id",
		}

		requestBody, err := json.Marshal(map[string]interface{}{
			"name":        "some-notification-type",
			"description": "some-notification-type-description",
			"critical":    false,
			"template_id": "some-template-id",
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("POST", "/senders/some-sender-id/notification_types", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		handler = notificationtypes.NewCreateHandler(notificationTypesCollection)
	})

	It("creates a notification type", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(notificationTypesCollection.AddCall.NotificationType).To(Equal(collections.NotificationType{
			Name:        "some-notification-type",
			Description: "some-notification-type-description",
			Critical:    false,
			TemplateID:  "some-template-id",
		}))
		Expect(notificationTypesCollection.AddCall.Conn).To(Equal(database.Conn))
		Expect(database.ConnectionWasCalled).To(BeTrue())

		Expect(writer.Code).To(Equal(http.StatusCreated))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-notification-type-id",
			"name": "some-notification-type",
			"description": "some-notification-type-description",
			"critical": false,
			"template_id": "some-template-id"
		}`))
	})

	It("allows the critical field to be omitted", func() {
		requestBody, err := json.Marshal(map[string]string{
			"name":        "some-notification-type",
			"description": "some-notification-type-description",
			"template_id": "some-template-id",
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("POST", "/senders/some-sender-id/notification_types", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(notificationTypesCollection.AddCall.NotificationType).To(Equal(collections.NotificationType{
			Name:        "some-notification-type",
			Description: "some-notification-type-description",
			Critical:    false,
			TemplateID:  "some-template-id",
		}))

		Expect(writer.Code).To(Equal(http.StatusCreated))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-notification-type-id",
			"name": "some-notification-type",
			"description": "some-notification-type-description",
			"critical": false,
			"template_id": "some-template-id"
		}`))
	})

	It("allows the template_id field to be omitted", func() {
		notificationTypesCollection.AddCall.ReturnNotificationType = collections.NotificationType{
			ID:          "some-notification-type-id",
			Name:        "some-notification-type",
			Description: "some-notification-type-description",
			Critical:    false,
			TemplateID:  "",
		}

		requestBody, err := json.Marshal(map[string]interface{}{
			"name":        "some-notification-type",
			"description": "some-notification-type-description",
			"critical":    false,
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("POST", "/senders/some-sender-id/notification_types", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(notificationTypesCollection.AddCall.NotificationType).To(Equal(collections.NotificationType{
			Name:        "some-notification-type",
			Description: "some-notification-type-description",
			Critical:    false,
			TemplateID:  "",
		}))

		Expect(writer.Code).To(Equal(http.StatusCreated))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-notification-type-id",
			"name": "some-notification-type",
			"description": "some-notification-type-description",
			"critical": false,
			"template_id": ""
		}`))
	})

	PIt("requires critical_notifications.write to create a critical notification type", func() {})

	Context("failure cases", func() {
		PIt("returns a 400 when the JSON request body cannot be decoded", func() {})

		PIt("returns a 422 when the request does not contain all the required fields", func() {})

		PIt("returns a 500 when there is a persistence error", func() {})

		PIt("returns a 422 when the template_id is not valid for the given client", func() {})
	})
})
