-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `v2_templates` (
      `id` varchar(36) NOT NULL,
      `name` varchar(255) DEFAULT NULL,
      `html` longtext DEFAULT NULL,
      `text` longtext DEFAULT NULL,
      `subject` varchar(255) DEFAULT NULL,
      `metadata` longtext DEFAULT NULL,
      `client_id` varchar(255) DEFAULT NULL,
      PRIMARY KEY (`id`),
      UNIQUE KEY `name_client_id` (`name`, `client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE v2_templates;
