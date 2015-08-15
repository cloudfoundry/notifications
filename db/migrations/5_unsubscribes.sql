-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `unsubscribes` (
      `primary` int(11) NOT NULL AUTO_INCREMENT,
      `user_id` varchar(255) DEFAULT NULL,
      `client_id` varchar(255) DEFAULT NULL,
      `kind_id` varchar(255) DEFAULT NULL,
      `created_at` datetime DEFAULT NULL,
      PRIMARY KEY (`primary`),
      UNIQUE KEY `user_id` (`user_id`,`client_id`,`kind_id`),
      UNIQUE KEY `user_id_2` (`user_id`,`client_id`,`kind_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `unsubscribes`;
