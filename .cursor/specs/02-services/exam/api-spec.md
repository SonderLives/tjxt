# exam 服务 HTTP API 接口规格

> 来源：`docs/tjxt.openapi.json` (聚合文档) | 最后同步：2026-08-05

## 接口清单

| 方法 | 路径 | 摘要 | 认证 | 请求体 | 响应体 | 权限标签 |
|------|------|------|------|--------|--------|----------|
| GET | /admin/notes/page | 管理端分页查询笔记 | 否 | - | R{data: R%C2%ABPageDTO%C2%ABNoteAdminVO%C2%BB%C2%BB} |  |
| GET | /admin/notes/{id} | 管理端查询笔记详情 | 否 | - | R{data: R%C2%ABNoteAdminDetailVO%C2%BB} |  |
| PUT | /admin/notes/{id}/hidden/{hidden} | 隐藏指定笔记 | 否 | - | R{data: R} |  |
| GET | /admin/questions/page | 管理端分页查询互动问题 | 否 | - | R{data: R%C2%ABPageDTO%C2%ABQuestionAdminVO%C2%BB%C2%BB} |  |
| GET | /admin/questions/{id} | 管理端根据id查询互动问题 | 否 | - | R{data: R%C2%ABQuestionAdminVO%C2%BB} |  |
| PUT | /admin/questions/{id}/hidden/{hidden} | 隐藏或显示问题 | 否 | - | R{data: R} |  |
| GET | /admin/replies/page | 分页查询回答或评论 | 否 | - | R{data: R%C2%ABPageDTO%C2%ABReplyVO%C2%BB%C2%BB} |  |
| GET | /admin/replies/{id} | 根据id查询回答或评论 | 否 | - | R{data: R%C2%ABReplyVO%C2%BB} |  |
| PUT | /admin/replies/{id}/hidden/{hidden} | 隐藏或显示评论 | 否 | - | R{data: R} |  |
| GET | /question-biz/biz/list | 批量查询与业务有关的题目id | 否 | - | R{data: R%C2%ABList%C2%ABQuestionBizDTO%C2%BB%C2%BB} |  |
| GET | /question-biz/biz/{id} | 查询与业务有关的题目id | 否 | - | R{data: R%C2%ABList%C2%ABQuestionBizDTO%C2%BB%C2%BB} |  |
| POST | /question-biz/list | 批量保存题目和业务关系 | 否 | JSON | R{data: R} |  |
| GET | /question-biz/scores | 查询业务下的题目分数和 | 否 | - | R{data: R%C2%ABMap%C2%ABlong%2Cint%C2%BB%C2%BB} |  |
| POST | /questions | 新增互动问题 | 否 | QuestionFormDTO | R{data: R} |  |
| GET | /questions/checkName | 校验名称是否有效，存在则无效返回false，不存在返回true | 否 | - | R{data: R%C2%ABboolean%C2%BB} |  |
| GET | /questions/list | 查询题目列表 | 否 | - | R{data: R%C2%ABList%C2%ABQuestionDTO%C2%BB%C2%BB} |  |
| GET | /questions/listOfBiz | 查询业务关联的题目列表 | 否 | - | R{data: R%C2%ABList%C2%ABQuestionDTO%C2%BB%C2%BB} |  |
| GET | /questions/numOfTeacher | 查询老师出题数量 | 否 | - | R{data: R%C2%ABMap%C2%ABlong%2Cint%C2%BB%C2%BB} |  |
| GET | /questions/page | 分页查询互动问题 | 否 | - | R{data: R%C2%ABPageDTO%C2%ABQuestionVO%C2%BB%C2%BB} |  |
| GET | /questions/scores | 查询题目分值 | 否 | - | R{data: R%C2%ABMap%C2%ABlong%2Cint%C2%BB%C2%BB} |  |
| DELETE | /questions/{id} | 根据id删除当前用户问题 | 否 | - | R{data: R} |  |
| GET | /questions/{id} | 根据id查询互动问题 | 否 | - | R{data: R%C2%ABQuestionVO%C2%BB} |  |
| PUT | /questions/{id} | 修改提问 | 否 | QuestionFormDTO | R{data: R} |  |

## 统一约定（引用全局规范）
- 响应格式：`pkg/response.R{Code,Msg,RequestId,Data any}`
- 分页：`PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`
- 错误码：见 [共享错误码表](../03-shared/error-codes.md)