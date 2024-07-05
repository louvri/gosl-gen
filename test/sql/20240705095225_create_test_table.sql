-- +goose Up
-- +goose StatementBegin
CREATE TABLE `test_table` (
     `id` bigint(20) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `field_value` VARCHAR(256) COLLATE utf8mb4_bin NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `test_table`;
-- +goose StatementEnd