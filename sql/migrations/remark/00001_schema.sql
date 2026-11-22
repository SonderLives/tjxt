-- tj_remark 库 schema 迁移（pressly/goose 管理）。
-- 源自原 sql/ddl/tj_remark.sql的建表部分（含索引），goose 化改造：Up 建表，Down 逆序 DROP。
-- 注：原 sql/migration/tj_remark.sql 无种子数据，故无 seed 迁移。

-- +goose Up

-- Table structure（源文件顺序）

CREATE TABLE `like_record` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '点赞记录id',
  `user_id` bigint NOT NULL COMMENT '点赞人id',
  `biz_id` bigint NOT NULL COMMENT '被点赞业务id',
  `biz_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '点赞业务类型，例如 reply / note / question',
  `liked` tinyint NOT NULL DEFAULT 1 COMMENT '1:点赞 0:取消',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_user_biz` (`user_id` ASC, `biz_id` ASC, `biz_type` ASC) USING BTREE,
  INDEX `idx_biz` (`biz_type` ASC, `biz_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '点赞记录' ROW_FORMAT = Dynamic;

-- +goose Down

-- 逆序 DROP（IF EXISTS 保证幂等）

DROP TABLE IF EXISTS `like_record`;
