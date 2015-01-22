-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
UPDATE `clients` SET `template_id` = "default" WHERE `template_id` = "" OR `template_id` IS NULL;
ALTER TABLE `clients` MODIFY COLUMN `template_id` varchar(255) NOT NULL DEFAULT "default";

UPDATE `kinds` SET `template_id` = "default" WHERE `template_id` = "" OR `template_id` IS NULL;
ALTER TABLE `kinds` MODIFY COLUMN `template_id` varchar(255) NOT NULL DEFAULT "default";

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `clients` MODIFY COLUMN `template_id` varchar(255) DEFAULT "";
ALTER TABLE `kinds` MODIFY COLUMN `template_id` varchar(255) DEFAULT "";
