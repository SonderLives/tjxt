# learning 服务 HTTP API 接口规格

> 来源：`docs/tjxt.openapi.json` (聚合文档) | 最后同步：2026-08-05

## 接口清单

| 方法 | 路径 | 摘要 | 认证 | 请求体 | 响应体 | 权限标签 |
|------|------|------|------|--------|--------|----------|
| GET | /boards | 分页查询指定赛季的积分排行榜 | 否 | - | R{data: R%C2%ABPointsBoardVO%C2%BB} |  |
| GET | /boards/seasons/list | 查询赛季信息列表 | 否 | - | R{data: R%C2%ABList%C2%ABPointsBoardSeasonVO%C2%BB%C2%BB} |  |
| POST | /learning-records | 提交学习记录 | 否 | LearningRecordFormDTO | R{data: R} |  |
| GET | /learning-records/course/{courseId} | 查询指定课程的学习记录 | 否 | - | R{data: R%C2%ABLearningLessonDTO%C2%BB} |  |
| GET | /lessons/now | 查询我正在学习的课程 | 否 | - | R{data: R%C2%ABLearningLessonVO%C2%BB} |  |
| GET | /lessons/page | 分页查询我的课表 | 否 | - | R{data: R%C2%ABPageDTO%C2%ABLearningLessonVO%C2%BB%C2%BB} |  |
| GET | /lessons/plans | 查询我的学习计划 | 否 | - | R{data: R%C2%ABLearningPlanPageVO%C2%BB} |  |
| POST | /lessons/plans | 创建学习计划 | 否 | LearningPlanDTO | R{data: R} |  |
| DELETE | /lessons/{courseId} | 删除指定课程信息 | 否 | - | R{data: R} |  |
| GET | /lessons/{courseId} | 查询指定课程信息 | 否 | - | R{data: R%C2%ABLearningLessonVO%C2%BB} |  |
| GET | /lessons/{courseId}/count | 统计课程学习人数 | 否 | - | R{data: R%C2%ABint%C2%BB} |  |
| GET | /lessons/{courseId}/valid | 校验当前课程是否已经报名 | 否 | - | R{data: R%C2%ABlong%C2%BB} |  |
| POST | /notes | 新增笔记 | 否 | NoteFormDTO | R{data: R} |  |
| DELETE | /notes/gathers/{id} | 取消采集笔记 | 否 | - | R{data: R} |  |
| POST | /notes/gathers/{id} | 采集笔记 | 否 | - | R{data: R} |  |
| GET | /notes/page | 用户端分页查询笔记 | 否 | - | R{data: R%C2%ABPageDTO%C2%ABNoteVO%C2%BB%C2%BB} |  |
| DELETE | /notes/{id} | 删除我的笔记 | 否 | - | R{data: R} |  |
| PUT | /notes/{id} | 更新笔记 | 否 | NoteFormDTO | R{data: R} |  |
| GET | /points/today | 查询我的今日积分 | 否 | - | R{data: R%C2%ABList%C2%ABPointsStatisticsVO%C2%BB%C2%BB} |  |
| POST | /replies | 新增回答或评论 | 否 | ReplyDTO | R{data: R} |  |
| GET | /replies/page | 分页查询回答或评论 | 否 | - | R{data: R%C2%ABPageDTO%C2%ABReplyVO%C2%BB%C2%BB} |  |
| GET | /sign-records | 查询签到记录 | 否 | - | R{data: R%C2%ABArray%C2%ABbyte%C2%BB%C2%BB} |  |
| POST | /sign-records | 签到功能接口 | 否 | - | R{data: R%C2%ABSignResultVO%C2%BB} |  |

## 统一约定（引用全局规范）
- 响应格式：`pkg/response.R{Code,Msg,RequestId,Data any}`
- 分页：`PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`
- 错误码：见 [共享错误码表](../03-shared/error-codes.md)