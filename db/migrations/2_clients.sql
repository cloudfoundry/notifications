-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `clients` (
      `primary` int(11) NOT NULL AUTO_INCREMENT,
      `id` varchar(255) DEFAULT NULL,
      `description` varchar(255) DEFAULT NULL,
      `created_at` datetime DEFAULT NULL,
      PRIMARY KEY (`primary`),
      UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `clients`;
