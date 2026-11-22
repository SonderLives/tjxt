-- tj_search 库 seed 迁移（pressly/goose 管理）。
-- 源自原 sql/migration/tj_search.sql 的 INSERT 初始化数据，goose 化改造：Up 灌数据，Down 按表清空。

-- +goose Up

-- 演示基线数据（按源文件顺序）

INSERT INTO `interests` VALUES (2, '2002,2054,2001', '2022-07-22 16:29:22', '2022-12-10 10:34:47');
INSERT INTO `interests` VALUES (129, '2002,2054,2004', '2022-10-26 10:16:58', '2022-11-29 15:24:25');

-- +goose Down

-- 按表清空（教学演示数据，无需精确逆序）

DELETE FROM `interests`;
