package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	"github.com/ryanmoran/stack"
)

type NotificationsByClient map[string]Client

type Client struct {
	Name          string                  `json:"name"`
	Template      string                  `json:"template"`
	Notifications map[string]Notification `json:"notifications"`
}

type Notification struct {
	Description string `json:"description"`
	Template    string `json:"template"`
	Critical    bool   `json:"critical"`
}

type GetAllNotifications struct {
	finder      services.NotificationsFinderInterface
	errorWriter ErrorWriterInterface
}

func NewGetAllNotifications(notificationsFinder services.NotificationsFinderInterface, errorWriter ErrorWriterInterface) GetAllNotifications {
	return GetAllNotifications{
		finder:      notificationsFinder,
		errorWriter: errorWriter,
	}
}

func (handler GetAllNotifications) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	clients, notifications, err := handler.finder.AllClientsAndNotifications()
	if err != nil {
		handler.errorWriter.Write(w, err)
		return
	}

	notificationsByClient := handler.constructNotifications(clients, notifications)

	response, err := json.Marshal(notificationsByClient)
	if err != nil {
		panic(err)
	}

	w.Write(response)
}

func (handler GetAllNotifications) constructNotifications(clients []models.Client, notifications []models.Kind) NotificationsByClient {
	notificationsByClient := NotificationsByClient{}

	for _, client := range clients {
		clientWithNotifications := Client{
			Name:     client.Description,
			Template: client.TemplateToUse(),
		}

		clientNotifications := make(map[string]Notification)
		for _, notification := range notifications {
			if notification.ClientID == client.ID {
				clientNotifications[notification.ID] = Notification{
					Description: notification.Description,
					Template:    "default",
					Critical:    notification.Critical,
				}
			}
		}

		clientWithNotifications.Notifications = clientNotifications
		notificationsByClient[client.ID] = clientWithNotifications
	}

	return notificationsByClient
}
