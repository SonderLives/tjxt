/*
  Navicat Premium Dump SQL - DDL Only for goctl model generation
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '课程分类id',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '分类名称',
  `parent_id` bigint NOT NULL DEFAULT 0 COMMENT '父分类id，一级分类父id为0',
  `level` int NOT NULL COMMENT '分类级别，1,2,3：代表一级分类，二级分类，三级分类',
  `priority` int NOT NULL DEFAULT 1 COMMENT '同级目录优先级，数字越小优先级越高，可以重复',
  `status` tinyint NOT NULL COMMENT '课程分类状态，1：正常，2：禁用',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` bigint NOT NULL DEFAULT 0 COMMENT '创建者',
  `updater` bigint NOT NULL DEFAULT 0 COMMENT '更新者',
  `deleted` tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3656 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '课程分类' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for course
-- ----------------------------
DROP TABLE IF EXISTS `course`;
CREATE TABLE `course`  (
  `id` bigint NOT NULL COMMENT '课程草稿id，对应正式草稿id',
  `name` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '课程名称',
  `course_type` tinyint NOT NULL DEFAULT 2 COMMENT '课程类型，1：直播课，2：录播课',
  `cover_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '封面链接',
  `first_cate_id` bigint NOT NULL COMMENT '一级课程分类id',
  `second_cate_id` bigint NOT NULL DEFAULT 0 COMMENT '二级课程分类id',
  `third_cate_id` bigint NOT NULL DEFAULT 0 COMMENT '三级课程分类id',
  `free` tinyint NOT NULL DEFAULT 0 COMMENT '售卖方式0付费，1：免费',
  `price` int NOT NULL COMMENT '课程价格，单位为分',
  `template_type` tinyint NOT NULL DEFAULT 1 COMMENT '模板类型，1：固定模板，2：自定义模板',
  `template_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '自定义模板的连接',
  `status` tinyint NOT NULL COMMENT '课程状态，1：待上架，2：已上架，3：下架，4：已完结',
  `purchase_start_time` datetime NULL DEFAULT NULL COMMENT '课程购买有效期开始时间',
  `purchase_end_time` datetime NOT NULL COMMENT '课程购买有效期结束时间',
  `step` tinyint NOT NULL COMMENT '信息填写进度',
  `score` int NULL DEFAULT 0 COMMENT '课程评价得分，45代表4.5星',
  `media_duration` int(10) UNSIGNED ZEROFILL NULL DEFAULT NULL COMMENT '课程总时长',
  `valid_duration` int NOT NULL COMMENT '课程有效期，单位月',
  `section_num` int NULL DEFAULT NULL COMMENT '课程总节数，包括练习',
  `dep_id` bigint NOT NULL COMMENT '部门id',
  `publish_times` int NULL DEFAULT 1 COMMENT '发布次数',
  `publish_time` datetime NULL DEFAULT NULL COMMENT '最近一次发布时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` bigint NOT NULL COMMENT '创建人',
  `updater` bigint NOT NULL COMMENT '更新人',
  `deleted` tinyint NOT NULL DEFAULT 0 COMMENT '逻辑删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '草稿课程' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for course_cata_subject_draft
-- ----------------------------
DROP TABLE IF EXISTS `course_cata_subject_draft`;
CREATE TABLE `course_cata_subject_draft`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '小节题目关系id',
  `course_id` bigint NULL DEFAULT NULL,
  `cata_id` bigint NOT NULL COMMENT '小节id',
  `subject_id` bigint NOT NULL COMMENT '题目id',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '小节题目关系草稿' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for course_catalogue
-- ----------------------------
DROP TABLE IF EXISTS `course_catalogue`;
CREATE TABLE `course_catalogue`  (
  `id` bigint NOT NULL COMMENT '目录id',
  `course_id` bigint NOT NULL COMMENT '课程id',
  `parent_id` bigint NOT NULL DEFAULT 0 COMMENT '父id，章的父id为0',
  `index` int NOT NULL COMMENT '序号',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '目录名称',
  `type` tinyint NOT NULL COMMENT '目录类型，1：章，2：节，3：练习',
  `media_duration` int NULL DEFAULT 0 COMMENT '媒资时长，单位秒',
  `trailer` tinyint NOT NULL DEFAULT 0 COMMENT '是否试看，0否，1是',
  `media_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '媒资名称',
  `media_id` bigint NULL DEFAULT 0 COMMENT '媒资id',
  `subject_num` int NULL DEFAULT 0 COMMENT '题目数量',
  `total_score` int NULL DEFAULT 0 COMMENT '总分',
  `can_update` tinyint NOT NULL DEFAULT 1 COMMENT '是否可更新',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_course_id`(`course_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '课程目录正式表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for course_catalogue_draft
-- ----------------------------
DROP TABLE IF EXISTS `course_catalogue_draft`;
CREATE TABLE `course_catalogue_draft`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '目录id',
  `course_id` bigint NOT NULL COMMENT '课程id',
  `parent_id` bigint NOT NULL DEFAULT 0 COMMENT '父id，章的父id为0',
  `index` int NOT NULL COMMENT '序号',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '目录名称',
  `type` tinyint NOT NULL COMMENT '目录类型，1：章，2：节，3：练习',
  `media_duration` int NULL DEFAULT 0 COMMENT '媒资时长，单位秒',
  `trailer` tinyint NOT NULL DEFAULT 0 COMMENT '是否试看，0否，1是',
  `media_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '媒资名称',
  `media_id` bigint NULL DEFAULT 0 COMMENT '媒资id',
  `subject_num` int NULL DEFAULT 0 COMMENT '题目数量',
  `total_score` int NULL DEFAULT 0 COMMENT '总分',
  `can_update` tinyint NOT NULL DEFAULT 1 COMMENT '是否可更新',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_course_id`(`course_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '课程目录草稿表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for course_content
-- ----------------------------
DROP TABLE IF EXISTS `course_content`;
CREATE TABLE `course_content`  (
  `id` bigint NOT NULL COMMENT '课程内容id，对应课程id',
  `course_introduce` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '课程简介',
  `course_detail` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '课程详情',
  `use_people` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '适用人群',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` bigint NOT NULL COMMENT '创建人',
  `updater` bigint NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '课程内容正式表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for course_content_draft
-- ----------------------------
DROP TABLE IF EXISTS `course_content_draft`;
CREATE TABLE `course_content_draft`  (
  `id` bigint NOT NULL COMMENT '课程内容草稿id，对应课程草稿id',
  `course_introduce` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '课程简介',
  `course_detail` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '课程详情',
  `use_people` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '适用人群',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` bigint NOT NULL COMMENT '创建人',
  `updater` bigint NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '课程内容草稿表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for course_draft
-- ----------------------------
DROP TABLE IF EXISTS `course_draft`;
CREATE TABLE `course_draft`  (
  `id` bigint NOT NULL COMMENT '课程草稿id',
  `name` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '课程名称',
  `course_type` tinyint NOT NULL DEFAULT 2 COMMENT '课程类型，1：直播课，2：录播课',
  `cover_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '封面链接',
  `first_cate_id` bigint NOT NULL COMMENT '一级课程分类id',
  `second_cate_id` bigint NOT NULL DEFAULT 0 COMMENT '二级课程分类id',
  `third_cate_id` bigint NOT NULL DEFAULT 0 COMMENT '三级课程分类id',
  `free` tinyint NOT NULL DEFAULT 0 COMMENT '售卖方式0付费，1：免费',
  `price` int NOT NULL COMMENT '课程价格，单位为分',
  `template_type` tinyint NOT NULL DEFAULT 1 COMMENT '模板类型，1：固定模板，2：自定义模板',
  `template_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '自定义模板的连接',
  `status` tinyint NOT NULL COMMENT '课程状态，1：待上架，2：已上架，3：下架，4：已完结',
  `purchase_start_time` datetime NULL DEFAULT NULL COMMENT '课程购买有效期开始时间',
  `purchase_end_time` datetime NOT NULL COMMENT '课程购买有效期结束时间',
  `step` tinyint NOT NULL COMMENT '信息填写进度',
  `score` int NULL DEFAULT 0 COMMENT '课程评价得分，45代表4.5星',
  `media_duration` int(10) UNSIGNED ZEROFILL NULL DEFAULT NULL COMMENT '课程总时长',
  `valid_duration` int NOT NULL COMMENT '课程有效期，单位月',
  `section_num` int NULL DEFAULT NULL COMMENT '课程总节数，包括练习',
  `dep_id` bigint NOT NULL COMMENT '部门id',
  `publish_times` int NULL DEFAULT 1 COMMENT '发布次数',
  `publish_time` datetime NULL DEFAULT NULL COMMENT '最近一次发布时间',
  `can_update` tinyint NOT NULL DEFAULT 1 COMMENT '是否可更新',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` bigint NOT NULL COMMENT '创建人',
  `updater` bigint NOT NULL COMMENT '更新人',
  `deleted` tinyint NOT NULL DEFAULT 0 COMMENT '逻辑删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '草稿课程' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for course_subject
-- ----------------------------
DROP TABLE IF EXISTS `course_subject`;
CREATE TABLE `course_subject`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '课程题目关系id',
  `course_id` bigint NOT NULL COMMENT '课程id',
  `cata_id` bigint NOT NULL COMMENT '目录id',
  `subject_id` bigint NOT NULL COMMENT '题目id',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_course_id`(`course_id` ASC) USING BTREE,
  INDEX `idx_cata_id`(`cata_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '课程题目关系正式表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for course_teacher
-- ----------------------------
DROP TABLE IF EXISTS `course_teacher`;
CREATE TABLE `course_teacher`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '课程老师关系id',
  `course_id` bigint NOT NULL COMMENT '课程id',
  `teacher_id` bigint NOT NULL COMMENT '老师id',
  `is_show` tinyint NOT NULL DEFAULT 1 COMMENT '是否前端展示',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_course_id`(`course_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '课程老师关系正式表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for course_teacher_draft
-- ----------------------------
DROP TABLE IF EXISTS `course_teacher_draft`;
CREATE TABLE `course_teacher_draft`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '课程老师关系草稿id',
  `course_id` bigint NOT NULL COMMENT '课程id',
  `teacher_id` bigint NOT NULL COMMENT '老师id',
  `is_show` tinyint NOT NULL DEFAULT 1 COMMENT '是否前端展示',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted` tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_course_id`(`course_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '课程老师关系草稿表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for subject
-- ----------------------------
DROP TABLE IF EXISTS `subject`;
CREATE TABLE `subject`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '题目id',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '题目内容',
  `type` tinyint NOT NULL COMMENT '题目类型，1单选，2多选，3判断',
  `difficulty` tinyint NOT NULL COMMENT '难度，1简单，2中等，3困难',
  `analysis` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '解析',
  `answer` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '答案',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` bigint NOT NULL COMMENT '创建人',
  `updater` bigint NOT NULL COMMENT '更新人',
  `deleted` tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '题目表' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;