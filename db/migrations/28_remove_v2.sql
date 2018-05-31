-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
DROP TABLE v2_templates;
DROP TABLE campaigns;
ALTER TABLE `messages` DROP COLUMN `campaign_id`;
DROP TABLE unsubscribers;
DROP TABLE campaign_types;
DROP TABLE senders;

-- +migrate Down
-- SQL in section 'Down' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `senders` (
  `id` varchar(36) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `client_id` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_client_id` (`name`, `client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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

RENAME TABLE notification_types to campaign_types;

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

CREATE TABLE IF NOT EXISTS `campaigns` (
  `id` varchar(36) NOT NULL,
  `campaign_type_id` varchar(36) DEFAULT NULL,
  `template_id` varchar(36) DEFAULT NULL,
  `sender_id` varchar(36) DEFAULT NULL,
  `send_to` longtext DEFAULT NULL,
  `text` longtext DEFAULT NULL,
  `html` longtext DEFAULT NULL,
  `subject` varchar(255) DEFAULT NULL,
  `reply_to` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `campaigns` ADD `status` varchar(255);
ALTER TABLE `campaigns` ADD `total_messages` integer;
ALTER TABLE `campaigns` ADD `sent_messages` integer;
ALTER TABLE `campaigns` ADD `retry_messages` integer;
ALTER TABLE `campaigns` ADD `failed_messages` integer;
ALTER TABLE `campaigns` ADD `start_time` datetime DEFAULT NULL;
ALTER TABLE `campaigns` ADD `completed_time` datetime DEFAULT NULL;
ALTER TABLE `messages` ADD `campaign_id` varchar(255);

CREATE TABLE IF NOT EXISTS `unsubscribers` (
  `id` varchar(36) NOT NULL,
  `campaign_type_id` varchar(36) DEFAULT NULL,
  `user_guid` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `unsubscribers` ADD UNIQUE KEY `user_guid_campaign_type_id` (`user_guid`, `campaign_type_id`);
ALTER TABLE `v2_templates` DROP KEY `name_client_id`;
