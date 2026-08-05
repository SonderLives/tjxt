# 数据域边界与数据库表归属

> 版本：v1.0 | 更新：2026-08-05

## 数据库实例与库分布

| 数据库实例 | 库名 | 归属服务 | 备注 |
|------------|------|----------|------|
| MySQL 127.0.0.1:3306 | tj_auth | auth | 用户、角色、权限、菜单、登录记录 |
| MySQL 127.0.0.1:3306 | tj_course | course | 课程、章节、资源、分类、评价 |
| MySQL 127.0.0.1:3306 | tj_trade | trade | 订单、订单项、支付记录、退款、分账 |
| MySQL 127.0.0.1:3306 | tj_learning | learning | 学习进度、笔记、收藏、证书、打卡 |
| MySQL 127.0.0.1:3306 | tj_pay | pay | 支付渠道、流水、对账单、退款单 |
| MySQL 127.0.0.1:3306 | tj_media | media | 文件元数据、存储桶、签名记录 |
| MySQL 127.0.0.1:3306 | tj_promotion | promotion | 优惠券、券码、用户券、活动规则 |
| MySQL 127.0.0.1:3306 | tj_message | message | 站内信、通知、模板、短信平台、发送任务 |
| MySQL 127.0.0.1:3306 | tj_exam | exam | 题库、试卷、考试、答卷、阅卷 |
| MySQL 127.0.0.1:3306 | tj_search | search | 兴趣标签、用户画像、ES 同步记录 |
| MySQL 127.0.0.1:3306 | tj_user | user | 用户详情、教师/学生档案、后台管理 |

## 核心表清单（各服务核心实体）

### auth - tj_auth
| 表名 | 核心字段 | 说明 |
|------|----------|------|
| `user` | id, username, password_hash, phone, email, status, role_id | 用户主表 |
| `user_detail` | user_id, real_name, avatar, gender, birthday | 用户详情 |
| `role` | id, name, code, description, status | 角色表 |
| `privilege` | id, name, code, type, parent_id, path, component | 权限/菜单树 |
| `menu` | id, name, path, component, icon, sort, parent_id | 菜单表 |
| `login_record` | id, user_id, ip, device, location, created_at | 登录日志 |

### course - tj_course
| 表名 | 核心字段 | 说明 |
|------|----------|------|
| `course` | id, title, cover, teacher_id, category_id, price, status | 课程主表 |
| `chapter` | id, course_id, title, sort, free_preview | 章节 |
| `resource` | id, chapter_id, type, url, duration, size | 资源(视频/文档/题目) |
| `category` | id, name, parent_id, level, sort | 分类树 |
| `course_evaluate` | id, course_id, user_id, score, content | 评价 |

### trade - tj_trade
| 表名 | 核心字段 | 说明 |
|------|----------|------|
| `order` | id, user_id, course_id, amount, status, pay_type, paid_at | 订单主表 |
| `order_item` | id, order_id, course_id, price, quantity | 订单明细 |
| `payment` | id, order_id, channel, transaction_id, amount, status | 支付流水 |
| `refund` | id, order_id, reason, amount, status, processed_at | 退款单 |
| `split_account` | id, order_id, receiver_id, amount, status | 分账记录 |

### learning - tj_learning
| 表名 | 核心字段 | 说明 |
|------|----------|------|
| `study_record` | id, user_id, course_id, chapter_id, progress, duration | 学习进度 |
| `note` | id, user_id, course_id, chapter_id, content, timestamp | 笔记 |
| `favorite` | id, user_id, target_type, target_id | 收藏 |
| `certificate` | id, user_id, course_id, certificate_no, issued_at | 结业证书 |
| `checkin` | id, user_id, course_id, date, streak_days | 打卡记录 |

### pay - tj_pay
| 表名 | 核心字段 | 说明 |
|------|----------|------|
| `pay_channel` | id, name, code, config, status | 支付渠道配置 |
| `pay_flow` | id, out_trade_no, channel, amount, status, callback_at | 支付流水 |
| `reconciliation` | id, date, channel, total_amount, diff_amount, status | 对账单 |
| `refund_flow` | id, pay_flow_id, out_refund_no, amount, status | 退款流水 |

### media - tj_media
| 表名 | 核心字段 | 说明 |
|------|----------|------|
| `file` | id, user_id, bucket, object_key, mime, size, status | 文件元数据 |
| `bucket` | id, name, endpoint, region, access_key, secret_key | 存储桶配置 |
| `signature` | id, file_id, policy, signature, expire_at | 上传签名记录 |

### promotion - tj_promotion
| 表名 | 核心字段 | 说明 |
|------|----------|------|
| `coupon` | id, name, type, discount, threshold, valid_days, total, status | 优惠券模板 |
| `coupon_code` | id, coupon_id, code, status, user_id, used_at | 券码池 |
| `user_coupon` | id, user_id, coupon_id, code, status, expire_at, used_at | 用户持有券 |
| `activity_rule` | id, name, condition, benefit, start_at, end_at | 活动规则 |

### message - tj_message
| 表名 | 核心字段 | 说明 |
|------|----------|------|
| `inbox` | id, user_id, title, content, type, read, read_at | 站内信收件箱 |
| `notice` | id, title, content, target_type, target_ids, sent_at | 系统通知 |
| `message_template` | id, code, channel, subject, content, variables | 消息模板 |
| `notice_task` | id, notice_id, status, execute_at, retry_count | 定时发送任务 |
| `sms_platform` | id, name, provider, config, status | 短信平台配置 |

### exam - tj_exam
| 表名 | 核心字段 | 说明 |
|------|----------|------|
| `question` | id, type, stem, options, answer, analysis, difficulty | 题目 |
| `question_biz` | id, question_id, course_id, chapter_id, tags | 题目业务关联 |
| `paper` | id, title, course_id, total_score, pass_score, duration | 试卷 |
| `paper_question` | id, paper_id, question_id, score, sort | 试卷题目 |
| `exam` | id, paper_id, user_id, score, status, started_at, submitted_at | 考试记录 |
| `answer_sheet` | id, exam_id, question_id, user_answer, score | 答卷 |

### search - tj_search
| 表名 | 核心字段 | 说明 |
|------|----------|------|
| `interests` | id, user_id, tags, weight, updated_at | 用户兴趣标签 |

### user - tj_user
| 表名 | 核心字段 | 说明 |
|------|----------|------|
| `user` | id, username, phone, email, status, role | 用户基础信息 |
| `user_detail` | user_id, real_name, avatar, gender, birthday, bio | 用户详情 |
| `teacher_profile` | user_id, title, introduction, certificates | 教师档案 |
| `student_profile` | user_id, grade, school, major | 学生档案 |

### data - 无独立数据库（聚合查询其他服务）

## 跨服务数据访问规则

| 场景 | 正确做法 | 禁止做法 |
|------|----------|----------|
| trade 需要课程价格 | 调用 course-rpc `GetCourse` | 直接查 `tj_course.course` |
| learning 需要用户信息 | 调用 user-rpc `GetUserById` | 直接查 `tj_user.user` |
| data 大屏需要交易汇总 | 调用 trade-rpc `GetTradeStats` | 直接查 `tj_trade.order` |
| promotion 校验用户身份 | 调用 auth-rpc `VerifyToken` | 共享 JWT secret 自己解析 |

## 迁移与版本管理

- **纯 DDL**：`sql/ddl/tj_<domain>.sql`（goctl model 专用，仅建表/索引）
- **含数据迁移**：`sql/migration/tj_<domain>.sql`（上线执行，包含初始数据、结构变更）
- **版本控制**：每次发布前必须同步更新对应 DDL 文件，CI 校验 model 生成无报错