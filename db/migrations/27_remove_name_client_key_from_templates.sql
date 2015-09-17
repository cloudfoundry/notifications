-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `v2_templates` DROP KEY `name_client_id`;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `v2_templates` ADD UNIQUE KEY `name_client_id` (`name`, `client_id`);
