-- +migrate Up
ALTER TABLE `jobs` MODIFY payload longtext;

-- +migrate Down
ALTER TABLE `jobs` MODIFY payload text;
