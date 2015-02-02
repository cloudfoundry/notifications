-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
UPDATE `templates` SET `metadata` = "{}" WHERE `metadata` = "" OR `metadata` IS NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
