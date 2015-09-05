-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `campaigns` ADD `status` varchar(255);
ALTER TABLE `campaigns` ADD `total_messages` integer;
ALTER TABLE `campaigns` ADD `sent_messages` integer;
ALTER TABLE `campaigns` ADD `retry_messages` integer;
ALTER TABLE `campaigns` ADD `failed_messages` integer;
ALTER TABLE `campaigns` ADD `start_time` datetime DEFAULT NULL;
ALTER TABLE `campaigns` ADD `completed_time` datetime DEFAULT NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `campaigns` DROP COLUMN `status`;
ALTER TABLE `campaigns` DROP COLUMN `total_messages`;
ALTER TABLE `campaigns` DROP COLUMN `sent_messages`;
ALTER TABLE `campaigns` DROP COLUMN `retry_messages`;
ALTER TABLE `campaigns` DROP COLUMN `failed_messages`;
ALTER TABLE `campaigns` DROP COLUMN `start_time`;
ALTER TABLE `campaigns` DROP COLUMN `completed_time`;
