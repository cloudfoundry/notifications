-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `templates` ADD `updated_at` datetime DEFAULT NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `templates` DROP COLUMN `updated_at`;
