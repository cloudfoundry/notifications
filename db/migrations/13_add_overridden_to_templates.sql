-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `templates` ADD `overridden` bool DEFAULT false;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `templates` DROP COLUMN `overridden`;
