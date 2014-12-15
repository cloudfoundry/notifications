package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"

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
			notificationsFinder.Clients = map[string]models.Client{
				"client-123": {
					ID:          "client-123",
					Description: "Jurassic Park",
				},
				"client-456": {
					ID:          "client-456",
					Description: "Jurassic Park Ride",
				},
			}

			notificationsFinder.Kinds = map[string]models.Kind{
				"perimeter-breach": {
					ID:          "perimeter-breach",
					Description: "very bad",
					Critical:    true,
					ClientID:    "client-123",
				},
				"fence-broken": {
					ID:          "fence-broken",
					Description: "even worse",
					Critical:    true,
					ClientID:    "client-123",
				},
				"perimeter-is-good": {
					ID:          "perimeter-is-good",
					Description: "very good",
					Critical:    false,
					ClientID:    "client-456",
				},
				"fence-works": {
					ID:          "fence-works",
					Description: "even better",
					Critical:    true,
					ClientID:    "client-456",
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
				notificationsFinder.AllClientsAndNotificationsError = errors.New("BANANA!!!")

				handler.ServeHTTP(writer, request, nil)

				Expect(errorWriter.Error).To(Equal(errors.New("BANANA!!!")))
			})
		})
	})
})
