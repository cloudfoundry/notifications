-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `global_unsubscribes` (
      `primary` int(11) NOT NULL AUTO_INCREMENT,
      `user_id` varchar(255) DEFAULT NULL,
      `created_at` datetime DEFAULT NULL,
      PRIMARY KEY (`primary`),
      UNIQUE KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `global_unsubscribes`;
