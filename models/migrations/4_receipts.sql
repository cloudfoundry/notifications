-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `receipts` (
      `primary` int(11) NOT NULL AUTO_INCREMENT,
      `user_guid` varchar(255) DEFAULT NULL,
      `client_id` varchar(255) DEFAULT NULL,
      `kind_id` varchar(255) DEFAULT NULL,
      `count` int(11) DEFAULT NULL,
      `created_at` datetime DEFAULT NULL,
      PRIMARY KEY (`primary`),
      UNIQUE KEY `user_guid` (`user_guid`,`client_id`,`kind_id`),
      UNIQUE KEY `user_guid_2` (`user_guid`,`client_id`,`kind_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `receipts`;
