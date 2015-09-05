-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `messages` ADD `campaign_id` varchar(255);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `messages` DROP COLUMN `campaign_id`;
