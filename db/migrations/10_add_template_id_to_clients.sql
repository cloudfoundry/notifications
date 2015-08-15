-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `clients` ADD `template_id` varchar(255) DEFAULT "";

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `clients` DROP COLUMN `template_id`;
