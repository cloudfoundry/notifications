-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `unsubscribers` ADD UNIQUE KEY `user_guid_campaign_type_id` (`user_guid`, `campaign_type_id`);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `unsubscribers` DROP KEY `user_guid_campaign_type_id`;
