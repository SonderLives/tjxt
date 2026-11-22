-- tj_search 库 schema 迁移（pressly/goose 管理）。
-- 源自原 sql/ddl/tj_search.sql的建表部分（含索引），goose 化改造：Up 建表，Down 逆序 DROP。

-- +goose Up

-- Table structure（源文件顺序）

CREATE TABLE `interests`  (
  `id` bigint NOT NULL COMMENT '主键，对应用户id',
  `interests` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '感兴趣的二级分类id，以逗号分隔，例如：120,220,330',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户兴趣表，保存感兴趣的二级分类id' ROW_FORMAT = Dynamic;

-- +goose Down

-- 逆序 DROP（IF EXISTS 保证幂等）

DROP TABLE IF EXISTS `interests`;
