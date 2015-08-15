-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `kinds` ADD `updated_at` datetime DEFAULT '2000-01-01 00:00:00';

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `kinds` DROP COLUMN `updated_at`;
