CREATE TABLE IF NOT EXISTS `zyy_es_users` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '姓名',
    `phone` char(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '手机号码',
    `password` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '密码',
    `interests` JSON DEFAULT NULL COMMENT '兴趣 [{"name":"","other_info":""}]',
    `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
    `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `phone_unique` (`phone`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
