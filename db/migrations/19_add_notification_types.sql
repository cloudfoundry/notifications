-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `notification_types` (
      `id` varchar(36) NOT NULL,
      `name` varchar(255) DEFAULT NULL,
      `description` varchar(255) DEFAULT NULL,
      `critical` bool DEFAULT FALSE,
      `template_id` varchar(255) DEFAULT NULL,
      `sender_id` varchar(255) DEFAULT NULL,
      PRIMARY KEY (`id`),
      UNIQUE KEY `name_sender_id` (`name`, `sender_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE notification_types;
