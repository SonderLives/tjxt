-- Pure DDL for goctl model generation (extracted from migration/tj_media.sql)
-- ----------------------------
-- Table structure for file
-- ----------------------------
DROP TABLE IF EXISTS `file`;
CREATE TABLE `file`  (
  `id` bigint NOT NULL COMMENT '主键，文件id',
  `key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文件在云端的唯一标示，例如：aaa.jpg',
  `filename` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文件上传时的名称',
  `request_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '请求id',
  `status` tinyint NOT NULL COMMENT '状态：1-待上传 2-已上传,未使用 3-已使用',
  `platform` tinyint NULL DEFAULT 1 COMMENT '平台：1-腾讯，2-阿里',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` bigint NOT NULL DEFAULT 0 COMMENT '创建者',
  `updater` bigint NOT NULL DEFAULT 0 COMMENT '更新者',
  `dep_id` bigint NOT NULL DEFAULT 0 COMMENT '部门id',
  `deleted` tinyint NOT NULL DEFAULT 0 COMMENT '逻辑删除，默认0',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '文件表，可以是普通文件、图片等' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of file
-- ----------------------------

-- ----------------------------
-- Table structure for media
-- ----------------------------
DROP TABLE IF EXISTS `media`;
CREATE TABLE `media`  (
  `id` bigint NOT NULL COMMENT '主键',
  `file_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文件在云端的唯一标示，例如：387702302659783576',
  `filename` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '文件名称',
  `media_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '媒体播放地址',
  `cover_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '媒体封面地址',
  `duration` double NOT NULL DEFAULT 0 COMMENT '视频时长，单位秒',
  `size` bigint NOT NULL DEFAULT 0 COMMENT '视频大小，单位是字节',
  `request_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '请求id',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态：1-上传中，2-已上传',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` bigint NOT NULL DEFAULT 0 COMMENT '创建者',
  `updater` bigint NOT NULL DEFAULT 0 COMMENT '更新者',
  `dep_id` bigint NOT NULL DEFAULT 0 COMMENT '部门id',
  `deleted` tinyint NOT NULL DEFAULT 0 COMMENT '逻辑删除，默认0',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '媒资表，主要是视频文件' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of media
-- ----------------------------
