-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `templates` ADD `id` varchar(255) DEFAULT NULL;
ALTER TABLE `templates` ADD UNIQUE KEY `id` (`id`);
ALTER TABLE `templates` DROP INDEX `name`;
ALTER TABLE `templates` ADD `subject` varchar(255) DEFAULT NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `templates` DROP COLUMN `id`;
ALTER TABLE `templates` DROP INDEX `id`;
ALTER TABLE `templates` ADD UNIQUE KEY `name` (`name`);
ALTER TABLE `templates` DROP COLUMN `subject`;

