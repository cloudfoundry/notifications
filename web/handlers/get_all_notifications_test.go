package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetAllNotifications", func() {
	var handler handlers.GetAllNotifications
	var writer *httptest.ResponseRecorder
	var request *http.Request
	var errorWriter *fakes.ErrorWriter
	var notificationsFinder *fakes.NotificationsFinder
	var err error

	BeforeEach(func() {
		errorWriter = fakes.NewErrorWriter()
		writer = httptest.NewRecorder()

		request, err = http.NewRequest("GET", "/notifications", nil)
		if err != nil {
			panic(err)
		}

		notificationsFinder = fakes.NewNotificationsFinder()
		handler = handlers.NewGetAllNotifications(notificationsFinder, errorWriter)
	})

	Describe("ServeHTTP", func() {
		It("receives the clients/notifications from the finder", func() {
			notificationsFinder.ClientsWithNotifications = map[string]services.ClientWithNotifications{
				"client-123": services.ClientWithNotifications{
					Name:     "Jurassic Park",
					Template: "default",
					Notifications: map[string]services.Notification{
						"perimeter-breach": {
							Description: "very bad",
							Template:    "default",
							Critical:    true,
						},
						"fence-broken": {
							Description: "even worse",
							Template:    "default",
							Critical:    true,
						},
					},
				},
				"client-456": services.ClientWithNotifications{
					Name:     "Jurassic Park Ride",
					Template: "default",
					Notifications: map[string]services.Notification{
						"perimeter-is-good": {
							Description: "very good",
							Template:    "default",
							Critical:    false,
						},
						"fence-works": {
							Description: "even better",
							Template:    "default",
							Critical:    true,
						},
					},
				},
			}

			handler.ServeHTTP(writer, request, nil)

			Expect(errorWriter.Error).To(BeNil())
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
		})

		Context("when the notifications finder errors", func() {
			It("delegates to the error writer", func() {
				notificationsFinder.AllClientNotificationsError = errors.New("BANANA!!!")

				handler.ServeHTTP(writer, request, nil)

				Expect(errorWriter.Error).To(Equal(errors.New("BANANA!!!")))
			})
		})
	})
})
