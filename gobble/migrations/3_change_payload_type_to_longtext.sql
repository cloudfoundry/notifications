-- +goose Up
ALTER TABLE `jobs` MODIFY payload longtext;

-- +goose Down
ALTER TABLE `jobs` MODIFY payload text;
