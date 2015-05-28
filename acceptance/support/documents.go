package support

type NotificationPreference struct {
	Count                   int    `json:"count"`
	Email                   bool   `json:"email"`
	NotificationDescription string `json:"kind_description"`
	SourceDescription       string `json:"source_description"`
}

type PreferenceDocument struct {
	GlobalUnsubscribe bool                         `json:"global_unsubscribe"`
	Clients           map[string]ClientPreferences `json:"clients,omitempty"`
}

type Template struct {
	Name     string                 `json:"name"`
	Subject  string                 `json:"subject"`
	Text     string                 `json:"text"`
	HTML     string                 `json:"html"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type TemplateListItem struct {
	ID   string
	Name string `json:"name"`
}

type TemplateAssociations struct {
	Associations []TemplateAssociation `json:"associations"`
}

type TemplateAssociation struct {
	ClientID       string `json:"client"`
	NotificationID string `json:"notification"`
}

type notifyRequest struct {
	To      string `json:"to,omitempty"`
	Role    string `json:"role,omitempty"`
	Subject string `json:"subject"`
	HTML    string `json:"html,omitempty"`
	Text    string `json:"text,omitempty"`
	KindID  string `json:"kind_id,omitempty"`
	ReplyTo string `json:"reply_to,omitempty"`
}

type NotifyResponse struct {
	Status         string `json:"status"`
	Recipient      string `json:"recipient"`
	NotificationID string `json:"notification_id"`
	VCAPRequestID  string `json:"vcap_request_id"`
}

type Message struct {
	Status string `json:"status"`
}

type RegisterClient struct {
	SourceName    string                          `json:"source_name"`
	Notifications map[string]RegisterNotification `json:"notifications,omitempty"`
}

type RegisterNotification struct {
	Description string `json:"description"`
	Critical    bool   `json:"critical"`
}

type NotificationsList map[string]NotificationClient

type NotificationClient struct {
	Name          string                  `json:"name"`
	Template      string                  `json:"template"`
	Notifications map[string]Notification `json:"notifications"`
}

type Notification struct {
	Description string `json:"description"`
	Template    string `json:"template"`
	Critical    bool   `json:"critical"`
}
