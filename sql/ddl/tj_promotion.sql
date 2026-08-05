-- 优惠券 / 促销中心数据库 DDL
-- 金额单位：分；时间：datetime
-- 风格对齐 tj_trade.order：InnoDB / utf8mb4 / 逻辑删除 deleted / creater+updater

CREATE TABLE IF NOT EXISTS `coupon` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '优惠券id',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '优惠券名称',
  `discount_type` varchar(20) NOT NULL DEFAULT '' COMMENT '优惠类型：reduce-满减 discount-折扣 no_threshold-无门槛',
  `discount_value` int NOT NULL DEFAULT '0' COMMENT '优惠值：满减金额/折扣百分比，单位分',
  `max_discount_amount` int NOT NULL DEFAULT '0' COMMENT '最高优惠金额，单位分',
  `threshold_amount` int NOT NULL DEFAULT '0' COMMENT '使用门槛金额，单位分，0表示无门槛',
  `obtain_way` varchar(20) NOT NULL DEFAULT '' COMMENT '获取方式：receive-领取 exchange-兑换码 assign-发放',
  `specific` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否限定适用范围：0-不限 1-限定',
  `scopes` json DEFAULT NULL COMMENT '适用范围id列表（课程/分类），specific=1时有效',
  `total_num` int NOT NULL DEFAULT '0' COMMENT '发行总量',
  `issue_num` int NOT NULL DEFAULT '0' COMMENT '已领取数量',
  `used_num` int NOT NULL DEFAULT '0' COMMENT '已使用数量',
  `user_limit` int NOT NULL DEFAULT '0' COMMENT '每人限领数量，0-不限',
  `issue_begin_time` datetime DEFAULT NULL COMMENT '发放开始时间',
  `issue_end_time` datetime DEFAULT NULL COMMENT '发放结束时间',
  `term_begin_time` datetime DEFAULT NULL COMMENT '使用开始时间',
  `term_end_time` datetime DEFAULT NULL COMMENT '使用结束时间',
  `term_days` int NOT NULL DEFAULT '0' COMMENT '有效期天数（相对领取日），与term_begin/end二选一',
  `status` varchar(20) NOT NULL DEFAULT 'draft' COMMENT '状态：draft-草稿 issued-已发放 paused-暂停 ended-结束',
  `remark` varchar(255) DEFAULT '' COMMENT '备注',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` bigint NOT NULL DEFAULT '0' COMMENT '创建人',
  `updater` bigint NOT NULL DEFAULT '0' COMMENT '更新人',
  `deleted` tinyint NOT NULL DEFAULT '0' COMMENT '逻辑删除',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_status` (`status`),
  KEY `idx_deleted` (`deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='优惠券';

CREATE TABLE IF NOT EXISTS `user_coupon` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户优惠券id',
  `user_id` bigint NOT NULL DEFAULT '0' COMMENT '用户id',
  `coupon_id` bigint NOT NULL DEFAULT '0' COMMENT '关联优惠券id',
  `status` varchar(20) NOT NULL DEFAULT 'unused' COMMENT '状态：unused-未使用 used-已使用 expired-已过期 refunded-已退',
  `obtain_time` datetime DEFAULT NULL COMMENT '领取时间',
  `use_time` datetime DEFAULT NULL COMMENT '使用时间',
  `expire_time` datetime DEFAULT NULL COMMENT '过期时间',
  `order_id` bigint DEFAULT NULL COMMENT '关联订单id（使用时记录）',
  `code` varchar(50) DEFAULT '' COMMENT '兑换码（兑换方式获取）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` bigint NOT NULL DEFAULT '0' COMMENT '创建人',
  `updater` bigint NOT NULL DEFAULT '0' COMMENT '更新人',
  `deleted` tinyint NOT NULL DEFAULT '0' COMMENT '逻辑删除',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_user` (`user_id`),
  KEY `idx_coupon` (`coupon_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户优惠券';

CREATE TABLE IF NOT EXISTS `coupon_code` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '兑换码id',
  `coupon_id` bigint NOT NULL DEFAULT '0' COMMENT '关联优惠券id',
  `code` varchar(50) NOT NULL DEFAULT '' COMMENT '兑换码',
  `status` varchar(20) NOT NULL DEFAULT 'unused' COMMENT '状态：unused-未使用 used-已使用',
  `user_id` bigint DEFAULT NULL COMMENT '兑换用户id',
  `expire_time` datetime DEFAULT NULL COMMENT '过期时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `creater` bigint NOT NULL DEFAULT '0' COMMENT '创建人',
  `updater` bigint NOT NULL DEFAULT '0' COMMENT '更新人',
  `deleted` tinyint NOT NULL DEFAULT '0' COMMENT '逻辑删除',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_coupon` (`coupon_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='优惠券兑换码';
