-- --------------------------------------------------------
-- 主机:                           192.168.150.101
-- 服务器版本:                        8.0.29 - MySQL Community Server - GPL
-- 服务器操作系统:                      Linux
-- HeidiSQL 版本:                  12.2.0.6576
-- --------------------------------------------------------
 
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
 
-- 导出  表 tj_learning.learning_lesson 结构
DROP TABLE IF EXISTS `learning_lesson`;
CREATE TABLE IF NOT EXISTS `learning_lesson` (
  `id` bigint NOT NULL COMMENT '主键',
  `user_id` bigint NOT NULL COMMENT '学员id',
  `course_id` bigint NOT NULL COMMENT '课程id',
  `status` tinyint DEFAULT '0' COMMENT '课程状态，0-未学习，1-学习中，2-已学完，3-已失效',
  `week_freq` tinyint DEFAULT NULL COMMENT '每周学习频率，例如每周学习6小节，则频率为6',
  `plan_status` tinyint NOT NULL DEFAULT '0' COMMENT '学习计划状态，0-没有计划，1-计划进行中',
  `learned_sections` int NOT NULL DEFAULT '0' COMMENT '已学习小节数量',
  `latest_section_id` bigint DEFAULT NULL COMMENT '最近一次学习的小节id',
  `latest_learn_time` datetime DEFAULT NULL COMMENT '最近一次学习的时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `expire_time` datetime DEFAULT NULL COMMENT '过期时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_user_id` (`user_id`,`course_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='学生课程表';
 
-- 正在导出表  tj_learning.learning_lesson 的数据：~4 rows (大约)
DELETE FROM `learning_lesson`;
INSERT INTO `learning_lesson` (`id`, `user_id`, `course_id`, `status`, `week_freq`, `plan_status`, `learned_sections`, `latest_section_id`, `latest_learn_time`, `create_time`, `expire_time`, `update_time`) VALUES
    (1, 2, 2, 2, 6, 1, 12, 16, '2023-04-11 22:34:45', '2022-08-05 20:02:50', '2023-08-05 20:02:29', '2023-04-19 10:29:29'),
    (2, 2, 3, 1, 4, 1, 3, 31, '2023-04-19 11:42:50', '2022-08-06 15:16:48', '2023-08-06 15:16:37', '2023-04-19 11:42:50'),
    (1585170299127607297, 129, 2, 0, NULL, 0, 0, 16, '2023-04-11 22:37:05', '2022-12-05 23:00:29', '2023-10-26 15:14:54', '2023-04-11 22:37:05'),
    (1601061367207464961, 2, 1549025085494521857, 1, 3, 1, 4, 1550383240983875589, '2023-04-11 16:34:44', '2022-12-09 11:49:11', '2023-12-09 11:49:11', '2023-04-11 16:34:43');
 
/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;






-- --------------------------------------------------------
-- 主机:                           192.168.150.101
-- 服务器版本:                        8.0.29 - MySQL Community Server - GPL
-- 服务器操作系统:                      Linux
-- HeidiSQL 版本:                  12.2.0.6576
-- --------------------------------------------------------
 
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
 
 
-- 导出 tj_learning 的数据库结构
CREATE DATABASE IF NOT EXISTS `tj_learning` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `tj_learning`;
 
-- 导出  表 tj_learning.learning_record 结构
CREATE TABLE IF NOT EXISTS `learning_record` (
  `id` bigint NOT NULL COMMENT '学习记录的id',
  `lesson_id` bigint NOT NULL COMMENT '对应课表的id',
  `section_id` bigint NOT NULL COMMENT '对应小节的id',
  `user_id` bigint NOT NULL COMMENT '用户id',
  `moment` int DEFAULT '0' COMMENT '视频的当前观看时间点，单位秒',
  `finished` bit(1) NOT NULL DEFAULT b'0' COMMENT '是否完成学习，默认false',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '第一次观看时间',
  `finish_time` datetime DEFAULT NULL COMMENT '完成学习的时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间（最近一次观看时间）',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_update_time` (`update_time`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE,
  KEY `idx_lesson_id` (`lesson_id`,`section_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='学习记录表';
 
-- 正在导出表  tj_learning.learning_record 的数据：~14 rows (大约)
DELETE FROM `learning_record`;
INSERT INTO `learning_record` (`id`, `lesson_id`, `section_id`, `user_id`, `moment`, `finished`, `create_time`, `finish_time`, `update_time`) VALUES
    (1582555272977592322, 1, 16, 2, 239, b'1', '2022-10-19 10:12:34', '2022-10-19 10:45:55', '2022-10-19 15:50:12'),
    (1582565729037791233, 1, 17, 2, 148, b'1', '2022-10-19 10:54:07', '2022-10-19 10:57:47', '2022-10-19 15:50:05'),
    (1582572518466760706, 1, 18, 2, 14, b'1', '2022-10-19 11:21:06', '2022-10-19 14:52:50', '2022-10-27 12:05:14'),
    (1582572534866489346, 1, 19, 2, 250, b'1', '2022-10-19 11:21:10', '2022-10-19 17:53:55', '2022-10-19 17:56:01'),
    (1582572535164284930, 1, 20, 2, 135, b'1', '2022-10-19 11:21:10', '2022-10-19 17:58:09', '2022-10-27 12:07:12'),
    (1582572537429209089, 1, 21, 2, 253, b'1', '2022-10-19 11:21:10', '2022-10-19 18:02:34', '2022-10-19 18:04:54'),
    (1582674101338656770, 1, 22, 2, 257, b'1', '2022-10-19 18:04:45', '2022-10-19 18:07:15', '2022-10-19 18:09:35'),
    (1582675262523326466, 1, 23, 2, 254, b'1', '2022-10-19 18:09:22', '2022-10-19 18:11:37', '2022-10-19 18:13:57'),
    (1582676374013886465, 1, 24, 2, 250, b'1', '2022-10-19 18:13:47', '2022-10-19 18:16:02', '2022-10-20 11:36:50'),
    (1582938844335001602, 1, 25, 2, 80, b'1', '2022-10-20 11:36:44', '2022-10-20 11:38:59', '2022-10-20 11:58:00'),
    (1583012729738776577, 1, 26, 2, 262, b'1', '2022-10-20 16:30:20', '2022-10-20 16:30:21', '2022-10-20 16:38:39'),
    (1586757474101342209, 1, 27, 2, 0, b'1', '2022-10-31 00:30:37', '2022-10-31 00:30:36', '2022-10-31 00:30:37'),
    (1586757474101342309, 2, 29, 2, 10, b'0', '2022-10-31 00:30:37', NULL, '2022-10-31 00:30:37'),
    (1599780755855228929, 1585170299127607297, 16, 129, 37, b'0', '2022-12-05 23:00:29', NULL, '2022-12-05 23:01:34');
 
/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;