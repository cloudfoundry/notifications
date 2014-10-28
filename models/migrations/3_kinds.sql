-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `kinds` (
      `primary` int(11) NOT NULL AUTO_INCREMENT,
      `id` varchar(255) DEFAULT NULL,
      `description` varchar(255) DEFAULT NULL,
      `critical` tinyint(1) DEFAULT NULL,
      `client_id` varchar(255) DEFAULT NULL,
      `created_at` datetime DEFAULT NULL,
      PRIMARY KEY (`primary`),
      UNIQUE KEY `id` (`id`,`client_id`),
      UNIQUE KEY `id_2` (`id`,`client_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `kinds`;
