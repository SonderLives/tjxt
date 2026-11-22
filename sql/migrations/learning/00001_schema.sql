-- tj_learning 库 schema 迁移（pressly/goose 管理）。
-- 源自原 sql/ddl/tj_learning.sql的建表部分（含索引），goose 化改造：Up 建表，Down 逆序 DROP。
-- 注：原 sql/migration/tj_learning.sql 为空库（仅建库语句），本 schema 以 goctl 生成源
-- sql/ddl/tj_learning.sql 为准（仅 learning_lesson 一张表）。

-- +goose Up

-- Table structure（源文件顺序）

CREATE TABLE IF NOT EXISTS learning_lesson (
  id BIGINT NOT NULL COMMENT '雪花主键',
  user_id BIGINT NOT NULL COMMENT '学员ID',
  course_id BIGINT NOT NULL COMMENT '课程ID',
  status TINYINT NOT NULL DEFAULT 0 COMMENT '0未开始，1学习中，2完成，3失效',
  week_freq INT NULL COMMENT '每周学习章节数',
  plan_status TINYINT NOT NULL DEFAULT 0 COMMENT '0无计划，1计划中',
  learned_sections INT NOT NULL DEFAULT 0 COMMENT '已学习章节数',
  latest_section_id BIGINT NULL COMMENT '最近学习小节',
  latest_learn_time DATETIME NULL COMMENT '最近学习时间',
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expire_time DATETIME NULL COMMENT '课程失效时间',
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_user_course (user_id, course_id),
  KEY idx_course_status (course_id, status),
  KEY idx_user_status_created (user_id, status, create_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='学员课程表';

-- +goose Down

-- 逆序 DROP（IF EXISTS 保证幂等）

DROP TABLE IF EXISTS `learning_lesson`;
