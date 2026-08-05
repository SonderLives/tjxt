-- Pure DDL for goctl model generation
-- ----------------------------
-- Table structure for like_record
-- ----------------------------
DROP TABLE IF EXISTS `like_record`;
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