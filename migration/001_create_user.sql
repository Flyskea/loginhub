-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `name` varchar(20) NOT NULL,
  `password` varchar(255) NOT NULL,
  `avatar` varchar(255) NOT NULL,
  `mobile` varchar(20),
  `email` varchar(255) NOT NULL,
  `last_login_at` bigint NOT NULL COMMENT '最后登陆时间',
  `ip` varchar(45) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` bigint NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_user_id`(`user_id`) USING BTREE,
  UNIQUE INDEX `idx_email`(`email`) USING BTREE,
  UNIQUE INDEX `idx_mobile`(`mobile`) USING BTREE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE `user`;

-- +goose StatementEnd