package services

type Response struct {
	Status         string `json:"status"`
	Recipient      string `json:"recipient"`
	NotificationID string `json:"notification_id"`
	VCAPRequestID  string `json:"vcap_request_id"`
}
