-- +goose Up
ALTER TABLE `jobs` ADD retry_count int(11) NOT NULL DEFAULT 0;
ALTER TABLE `jobs` ADD active_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- +goose Down
ALTER TABLE `jobs` DROP COLUMN retry_count;
ALTER TABLE `jobs` DROP COLUMN active_at;
