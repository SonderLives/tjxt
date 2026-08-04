-- Pure DDL for goctl model generation (extracted from migration/tj_exam.sql)
-- ----------------------------
-- Table structure for question
-- ----------------------------
DROP TABLE IF EXISTS `question`;
CREATE TABLE `question`  (
  `id` bigint NOT NULL COMMENT '题目id',
  `name` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '题干',
  `type` tinyint NOT NULL COMMENT '题目类型，1：单选题，2：多选题，3：不定向选择题，4：判断题，5：主观题',
  `cate_id1` bigint NOT NULL COMMENT '1级课程分类id',
  `cate_id2` bigint NOT NULL COMMENT '2级课程分类id',
  `cate_id3` bigint NOT NULL COMMENT '3级课程分类id',
  `difficulty` tinyint NOT NULL COMMENT '难易度，1：简单，2：中等，3：困难',
  `answer_times` int NOT NULL DEFAULT 0 COMMENT '回答次数',
  `correct_times` int NOT NULL DEFAULT 0 COMMENT '回答正确次数',
  `score` int NOT NULL COMMENT '分值',
  `dep_id` bigint NULL DEFAULT NULL COMMENT '部门id',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` bigint NOT NULL DEFAULT 1 COMMENT '创建人',
  `updater` bigint NOT NULL DEFAULT 1 COMMENT '更新人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '题目' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of question
-- ----------------------------

-- ----------------------------
-- Table structure for question_biz
-- ----------------------------
DROP TABLE IF EXISTS `question_biz`;
CREATE TABLE `question_biz`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `biz_id` bigint NULL DEFAULT NULL COMMENT '业务id，要关联问题的某业务id，例如小节id',
  `question_id` bigint NULL DEFAULT NULL COMMENT '问题id',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `biz_id`(`biz_id` ASC, `question_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 147 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '问题和业务关联表，例如把小节id和问题id关联，一个小节下可以有多个问题' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of question_biz
-- ----------------------------

-- ----------------------------
-- Table structure for question_detail
-- ----------------------------
DROP TABLE IF EXISTS `question_detail`;
CREATE TABLE `question_detail`  (
  `id` bigint NOT NULL COMMENT '题目id',
  `options` json NULL COMMENT '选择题选项，json数组格式',
  `answer` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '选择题正确答案1到10，如果有多个答案，中间使用逗号隔开，如果是判断题，1：代表正确，其他代表错误',
  `analysis` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '答案解析',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '题目' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of question_detail
-- ----------------------------

-- ----------------------------
-- Table structure for undo_log
-- ----------------------------
DROP TABLE IF EXISTS `undo_log`;
CREATE TABLE `undo_log`  (
  `branch_id` bigint NOT NULL COMMENT 'branch transaction id',
  `xid` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'global transaction id',
  `context` varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'undo_log context,such as serialization',
  `rollback_info` longblob NOT NULL COMMENT 'rollback info',
  `log_status` int NOT NULL COMMENT '0:normal status,1:defense status',
  `log_created` datetime(6) NOT NULL COMMENT 'create datetime',
  `log_modified` datetime(6) NOT NULL COMMENT 'modify datetime',
  UNIQUE INDEX `ux_undo_log`(`xid` ASC, `branch_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = 'AT transaction mode undo table' ROW_FORMAT = COMPACT;

-- ----------------------------
-- Records of undo_log
-- ----------------------------
