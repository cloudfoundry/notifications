-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
RENAME TABLE notification_types to campaign_types;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
RENAME TABLE campaign_types to notification_types;