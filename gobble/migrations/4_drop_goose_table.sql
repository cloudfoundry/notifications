-- +migrate Up
DROP TABLE IF EXISTS `goose_db_version`;

-- +migrate Down
CREATE TABLE IF NOT EXISTS `goose_db_version`;
