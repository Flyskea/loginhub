-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `user_third_auth`  (
  `id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `type` varchar(50) NOT NULL COMMENT 'OAUTH2 provider type, e.g. github, wechat, etc.',
  `auth_id` varchar(50) NOT NULL COMMENT '第三方 uid 、openid 等',
  `union_id` varchar(50) NOT NULL COMMENT 'QQ / 微信同一主体下 Unionid 相同',
  `credential` varchar(255) NOT NULL COMMENT 'access_token',
  `refresh_token` varchar(255) NOT NULL COMMENT 'refresh_token',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_user_id`(`user_id`) USING BTREE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE `oauth2_provider`;

-- +goose StatementEnd