-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `unsubscribers` (
      `id` varchar(36) NOT NULL,
      `campaign_type_id` varchar(36) DEFAULT NULL,
      `user_guid` varchar(255) DEFAULT NULL,
      PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE unsubscribers;
