-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `templates` ADD `metadata` longtext;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `templates` DROP COLUMN `metadata`;
