-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
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

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE campaigns;
