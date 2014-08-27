package handlers

import "github.com/cloudfoundry-incubator/notifications/models"

type NotificationPreferences map[string]map[string]map[string]bool

func NewNotificationPreferences() NotificationPreferences {
    return map[string]map[string]map[string]bool{}
}

func (pref NotificationPreferences) Add(client string, kind string, emails bool) {
    if clientMap, ok := pref[client]; ok {
        clientMap[kind] = map[string]bool{
            "email": emails,
        }
    } else {
        pref[client] = map[string]map[string]bool{
            kind: map[string]bool{
                "email": emails,
            },
        }

    }

}
