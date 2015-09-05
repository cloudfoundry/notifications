package models

import "time"

type statusCount struct {
	Status string `db:"status"`
	Count  int    `db:"count"`
}

type MessageCounts struct {
	Total     int
	Failed    int
	Delivered int
}

type Message struct {
	ID         string    `db:"id"`
	CampaignID string    `db:"campaign_id"`
	Status     string    `db:"status"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type clock interface {
	Now() time.Time
}

type MessagesRepository struct {
	clock clock
}

func NewMessagesRepository(clock clock) MessagesRepository {
	return MessagesRepository{
		clock: clock,
	}
}

func (mr MessagesRepository) CountByStatus(conn ConnectionInterface, campaignID string) (MessageCounts, error) {
	var counts []statusCount
	var messageCounts MessageCounts

	_, err := conn.Select(&counts, "SELECT `status`, COUNT(id) AS `count` FROM `messages` WHERE `campaign_id` = ? GROUP BY `status`", campaignID)
	if err != nil {
		return messageCounts, err
	}

	for _, count := range counts {
		switch count.Status {
		case "delivered":
			messageCounts.Delivered = count.Count
		case "failed":
			messageCounts.Failed = count.Count
		}
		messageCounts.Total += count.Count
	}

	return messageCounts, nil
}

func (mr MessagesRepository) MostRecentlyUpdatedByCampaignID(conn ConnectionInterface, campaignID string) (Message, error) {
	var message Message
	err := conn.SelectOne(&message, "SELECT * FROM `messages` WHERE `campaign_id` = ? ORDER BY `updated_at` DESC LIMIT 1", campaignID)
	if err != nil {
		return message, err
	}

	return message, nil
}

func (mr MessagesRepository) Insert(conn ConnectionInterface, message Message) (Message, error) {
	err := conn.Insert(&message)
	if err != nil {
		return Message{}, err
	}

	return message, nil
}

func (mr MessagesRepository) Update(conn ConnectionInterface, message Message) (Message, error) {
	message.UpdatedAt = mr.clock.Now()

	_, err := conn.Update(&message)
	if err != nil {
		return Message{}, err
	}

	return message, nil
}
