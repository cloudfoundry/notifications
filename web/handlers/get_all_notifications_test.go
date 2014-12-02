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
					Name: "Jurassic Park",
					Notifications: map[string]services.Notification{
						"perimeter-breach": {
							Description: "very bad",
							Critical:    true,
						},
						"fence-broken": {
							Description: "even worse",
							Critical:    true,
						},
					},
				},
				"client-456": services.ClientWithNotifications{
					Name: "Jurassic Park Ride",
					Notifications: map[string]services.Notification{
						"perimeter-is-good": {
							Description: "very good",
							Critical:    false,
						},
						"fence-works": {
							Description: "even better",
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
					"notifications": {
						"perimeter-breach": {
							"description": "very bad",
							"critical": true
						},
						"fence-broken": {
							"description": "even worse",
							"critical": true
						}
					}
				},
				"client-456": {
					"name": "Jurassic Park Ride",
					"notifications": {
						"perimeter-is-good": {
							"description": "very good",
							"critical": false
						},
						"fence-works": {
							"description": "even better",
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
