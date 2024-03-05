-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `oauth2_provider` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `type` varchar(20) NOT NULL COMMENT 'OAUTH2 provider type, e.g. github, wechat, etc.',
  `client_id` varchar(255) NOT NULL,
  `client_secret` varchar(255) NOT NULL,
  `redirect_url` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` bigint NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_type`(`type`) USING BTREE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE `oauth2_provider`;

-- +goose StatementEnd