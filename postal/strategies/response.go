package strategies

type Response struct {
    Status         string `json:"status"`
    Recipient      string `json:"recipient"`
    NotificationID string `json:"notification_id"`
    Email          string `json:"email"`
}
