package notifications_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notifications"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListHandler", func() {
	var (
		handler             notifications.ListHandler
		writer              *httptest.ResponseRecorder
		request             *http.Request
		errorWriter         *mocks.ErrorWriter
		notificationsFinder *mocks.NotificationsFinder
		err                 error
		database            *mocks.Database
		context             stack.Context
	)

	BeforeEach(func() {
		errorWriter = mocks.NewErrorWriter()
		writer = httptest.NewRecorder()
		database = mocks.NewDatabase()
		context = stack.NewContext()
		context.Set("database", database)

		request, err = http.NewRequest("GET", "/notifications", nil)
		Expect(err).NotTo(HaveOccurred())

		notificationsFinder = mocks.NewNotificationsFinder()
		handler = notifications.NewListHandler(notificationsFinder, errorWriter)
	})

	Describe("ServeHTTP", func() {
		It("receives the clients/notifications from the finder", func() {
			notificationsFinder.AllClientsAndNotificationsCall.Returns.Clients = []models.Client{
				{
					ID:          "client-123",
					Description: "Jurassic Park",
				},
				{
					ID:          "client-456",
					Description: "Jurassic Park Ride",
				},
			}

			notificationsFinder.AllClientsAndNotificationsCall.Returns.Kinds = []models.Kind{
				{
					ID:          "perimeter-breach",
					Description: "very bad",
					Critical:    true,
					ClientID:    "client-123",
				},
				{
					ID:          "fence-broken",
					Description: "even worse",
					Critical:    true,
					ClientID:    "client-123",
				},
				{
					ID:          "perimeter-is-good",
					Description: "very good",
					Critical:    false,
					ClientID:    "client-456",
				},
				{
					ID:          "fence-works",
					Description: "even better",
					Critical:    true,
					ClientID:    "client-456",
				},
			}

			handler.ServeHTTP(writer, request, context)

			Expect(errorWriter.WriteCall.Receives.Error).To(BeNil())
			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.Body.Bytes()).To(MatchJSON(`{
				"client-123": {
					"name": "Jurassic Park",
					"template": "default",
					"notifications": {
						"perimeter-breach": {
							"description": "very bad",
							"template": "default",
							"critical": true
						},
						"fence-broken": {
							"description": "even worse",
							"template": "default",
							"critical": true
						}
					}
				},
				"client-456": {
					"name": "Jurassic Park Ride",
					"template": "default",
					"notifications": {
						"perimeter-is-good": {
							"description": "very good",
							"template": "default",
							"critical": false
						},
						"fence-works": {
							"description": "even better",
							"template": "default",
							"critical": true
						}
					}
				}
			}`))

			Expect(notificationsFinder.AllClientsAndNotificationsCall.Receives.Database).To(Equal(database))
		})

		Context("when the notifications finder errors", func() {
			It("delegates to the error writer", func() {
				notificationsFinder.AllClientsAndNotificationsCall.Returns.Error = errors.New("BANANA!!!")

				handler.ServeHTTP(writer, request, context)

				Expect(errorWriter.WriteCall.Receives.Error).To(Equal(errors.New("BANANA!!!")))
			})
		})
	})
})
