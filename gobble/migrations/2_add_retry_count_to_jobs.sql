-- +migrate Up
SET @preparedStatement = (SELECT IF(
    (SELECT COUNT(*)
        FROM INFORMATION_SCHEMA.COLUMNS
        WHERE  table_name = 'jobs'
        AND table_schema = DATABASE()
        AND column_name = 'retry_count'
    ) > 0,
    "SELECT 1",
    "ALTER TABLE `jobs` ADD `retry_count` INT(11) NOT NULL DEFAULT '0';"
));

PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

SET @preparedStatement = (SELECT IF(
    (SELECT COUNT(*)
        FROM INFORMATION_SCHEMA.COLUMNS
        WHERE  table_name = 'jobs'
        AND table_schema = DATABASE()
        AND column_name = 'active_at'
    ) > 0,
    "SELECT 1",
    "ALTER TABLE `jobs` ADD `active_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;"
));

PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

-- +migrate Down
ALTER TABLE `jobs` DROP COLUMN retry_count;
ALTER TABLE `jobs` DROP COLUMN active_at;
