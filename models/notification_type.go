package models

type NotificationType struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Critical    bool   `db:"critical"`
	TemplateID  string `db:"template_id"`
	SenderID    string `db:"sender_id"`
}
