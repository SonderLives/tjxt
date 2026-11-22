-- tj_message 库 seed 迁移（pressly/goose 管理）。
-- 源自原 sql/migration/tj_message.sql 的 INSERT 初始化数据，goose 化改造：Up 灌数据，Down 按表清空。

-- +goose Up

-- 演示基线数据（按源文件顺序）

INSERT INTO `message_template` VALUES (1561895814765240322, '短信验证码', 'aliYun', '天机学堂', 'SMS_251685292', '验证码:{code},您在使用天机学堂短信验证功能，仅限本人使用，请勿向他人泄露验证码信息', 1561895814668771330, 0, 0, 0, '2022-08-23 01:59:16', '2022-09-24 11:23:32');
INSERT INTO `message_template` VALUES (1561897355307925505, '短信验证码', 'tencent', '天机学堂', 'SMS_185822702', '您的验证码为xxx,验证码5分钟内有效，请勿泄露给他人！', 1561895814668771330, 0, 0, 0, '2022-08-23 02:05:24', '2022-08-23 02:05:24');
INSERT INTO `notice_template` VALUES (1561895814668771330, '短信验证码', 'VERIFY_CODE', 4, 1, NULL, '您的验证码为xxx,验证码5分钟内有效，请勿泄露给他人！', b'1', 0, 0, '2022-08-23 01:59:16', '2022-09-24 11:31:08');
INSERT INTO `sms_third_platform` VALUES (1, '阿里云', 'aliYun', 0, 1, 0, 0, '2022-08-23 01:04:57', '2022-08-23 01:04:57');
INSERT INTO `sms_third_platform` VALUES (2, '腾讯云', 'tencent', 1, 1, 0, 0, '2022-08-23 01:05:18', '2022-08-23 01:05:18');
INSERT INTO `sms_third_platform` VALUES (3, '优刻得', 'uCloud', 3, 0, 0, 0, '2022-08-23 01:06:14', '2022-08-23 01:07:19');

-- +goose Down

-- 按表清空（教学演示数据，无需精确逆序）

DELETE FROM `sms_third_platform`;

DELETE FROM `notice_template`;

DELETE FROM `message_template`;
