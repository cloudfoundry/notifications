package support

type Preference struct {
	ClientID                string
	NotificationID          string
	Count                   int
	Email                   bool
	NotificationDescription string
	SourceDescription       string
}

type Preferences struct {
	GlobalUnsubscribe       bool
	NotificationPreferences []Preference
}

type ClientPreferences map[string]NotificationPreference

func (response PreferenceDocument) Preferences() Preferences {
	var preferences []Preference

	for clientID, clientPreferences := range response.Clients {
		for notificationID, notificationPreference := range clientPreferences {
			preferences = append(preferences, Preference{
				ClientID:       clientID,
				NotificationID: notificationID,
				Count:          notificationPreference.Count,
				Email:          notificationPreference.Email,
				NotificationDescription: notificationPreference.NotificationDescription,
				SourceDescription:       notificationPreference.SourceDescription,
			})
		}
	}

	return Preferences{
		GlobalUnsubscribe:       response.GlobalUnsubscribe,
		NotificationPreferences: preferences,
	}
}
