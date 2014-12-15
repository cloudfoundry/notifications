package services

import "github.com/cloudfoundry-incubator/notifications/models"

type ClientWithNotifications struct {
	Name          string                  `json:"name"`
	Template      string                  `json:"template"`
	Notifications map[string]Notification `json:"notifications"`
}

type Notification struct {
	Description string `json:"description"`
	Template    string `json:"template"`
	Critical    bool   `json:"critical"`
}

type NotificationsFinder struct {
	clientsRepo models.ClientsRepoInterface
	kindsRepo   models.KindsRepoInterface
	database    models.DatabaseInterface
}

type NotificationsFinderInterface interface {
	AllClientNotifications() (map[string]ClientWithNotifications, error)
	ClientAndKind(string, string) (models.Client, models.Kind, error)
}

func NewNotificationsFinder(clientsRepo models.ClientsRepoInterface, kindsRepo models.KindsRepoInterface, database models.DatabaseInterface) NotificationsFinder {
	return NotificationsFinder{
		clientsRepo: clientsRepo,
		kindsRepo:   kindsRepo,
		database:    database,
	}
}

func (finder NotificationsFinder) AllClientNotifications() (map[string]ClientWithNotifications, error) {
	clients, err := finder.allClients()
	if err != nil {
		return map[string]ClientWithNotifications{}, err
	}

	allClientNotifications := map[string]ClientWithNotifications{}
	for _, client := range clients {
		clientWithNotifications, err := finder.getAllNotificationsForClient(client)
		if err != nil {
			return map[string]ClientWithNotifications{}, err
		}
		allClientNotifications[client.ID] = clientWithNotifications
	}

	return allClientNotifications, nil
}

func (finder NotificationsFinder) getAllNotificationsForClient(client models.Client) (ClientWithNotifications, error) {
	kinds, err := finder.allKindsForClient(client.ID)
	if err != nil {
		return ClientWithNotifications{}, err
	}

	notifications := finder.kindsToNotifications(kinds)
	clientWithNotifications := ClientWithNotifications{
		Name:          client.Description,
		Template:      "default",
		Notifications: notifications,
	}
	return clientWithNotifications, nil
}

func (finder NotificationsFinder) kindsToNotifications(kinds []models.Kind) map[string]Notification {
	notifications := map[string]Notification{}
	for _, kind := range kinds {
		notification := Notification{
			Description: kind.Description,
			Template:    "default",
			Critical:    kind.Critical,
		}
		notifications[kind.ID] = notification
	}

	return notifications
}

func (finder NotificationsFinder) allClients() ([]models.Client, error) {
	clientsRepo := finder.clientsRepo
	clients, err := clientsRepo.FindAll(finder.database.Connection())
	if err != nil {
		return []models.Client{}, err
	}

	return clients, nil
}

func (finder NotificationsFinder) allKindsForClient(clientID string) ([]models.Kind, error) {
	kinds, err := finder.kindsRepo.FindByClient(finder.database.Connection(), clientID)
	if err != nil {
		return []models.Kind{}, err
	}

	return kinds, nil
}

func (finder NotificationsFinder) ClientAndKind(clientID, kindID string) (models.Client, models.Kind, error) {
	client, err := finder.client(clientID)
	if err != nil {
		return models.Client{}, models.Kind{}, err
	}

	kind, err := finder.kind(clientID, kindID)
	if err != nil {
		return client, models.Kind{}, err
	}

	return client, kind, nil
}

func (finder NotificationsFinder) client(clientID string) (models.Client, error) {
	client, err := finder.clientsRepo.Find(finder.database.Connection(), clientID)
	if err != nil {
		if _, ok := err.(models.ErrRecordNotFound); ok {
			return models.Client{ID: clientID}, nil
		} else {
			return models.Client{}, err
		}
	}
	return client, nil
}

func (finder NotificationsFinder) kind(clientID, kindID string) (models.Kind, error) {
	kind, err := finder.kindsRepo.Find(finder.database.Connection(), kindID, clientID)
	if err != nil {
		if _, ok := err.(models.ErrRecordNotFound); ok {
			return models.Kind{ID: kindID, ClientID: clientID}, nil
		} else {
			return models.Kind{}, err
		}
	}
	return kind, nil
}
