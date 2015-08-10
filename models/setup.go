package models

func Setup(database *DB) {
	database.connection.AddTableWithName(Client{}, "clients").SetKeys(true, "Primary").ColMap("ID").SetUnique(true)
	database.connection.AddTableWithName(Kind{}, "kinds").SetKeys(true, "Primary").SetUniqueTogether("id", "client_id")
	database.connection.AddTableWithName(Receipt{}, "receipts").SetKeys(true, "Primary").SetUniqueTogether("user_guid", "client_id", "kind_id")
	database.connection.AddTableWithName(Unsubscribe{}, "unsubscribes").SetKeys(true, "Primary").SetUniqueTogether("user_id", "client_id", "kind_id")
	database.connection.AddTableWithName(GlobalUnsubscribe{}, "global_unsubscribes").SetKeys(true, "Primary").ColMap("UserID").SetUnique(true)
	database.connection.AddTableWithName(Template{}, "templates").SetKeys(true, "Primary").ColMap("Name").SetUnique(true)
	database.connection.AddTableWithName(Message{}, "messages").SetKeys(false, "ID")
	database.connection.AddTableWithName(Sender{}, "senders").SetKeys(false, "ID").SetUniqueTogether("name", "client_id")
	database.connection.AddTableWithName(CampaignType{}, "campaign_types").SetKeys(false, "ID").SetUniqueTogether("name", "sender_id")
}
