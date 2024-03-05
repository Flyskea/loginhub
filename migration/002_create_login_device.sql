-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `login_device` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `user_id` bigint NOT NULL,
    `device_id` binary(16) NOT NULL,
    `browser` varchar(50) NOT NULL,
    `os` varchar(20) NOT NULL,
    `ip` varchar(40) NOT NULL,
    `created_at` datetime(3) NOT NULL,
    `updated_at` datetime(3) NOT NULL,
    `deleted_at` bigint NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    INDEX `idx_user_id`(`user_id`) USING BTREE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE `login_device`;

-- +goose StatementEnd