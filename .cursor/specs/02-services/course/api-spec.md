# course 服务 HTTP API 接口规格

> 来源：`docs/tjxt.openapi.json` (聚合文档) | 最后同步：2026-08-05

## 接口清单

| 方法 | 路径 | 摘要 | 认证 | 请求体 | 响应体 | 权限标签 |
|------|------|------|------|--------|--------|----------|
| GET | /catalogues/batchQuery | 根据章节目录批量查询基础信息 | 否 | - | R{data: R%C2%ABList%C2%ABCataSimpleInfoVO%C2%BB%C2%BB} |  |
| GET | /catalogues/querySectionInfoById/{id} | 获取小节信息 | 否 | - | R{data: R%C2%ABCataSimpleInfoVO%C2%BB} |  |
| POST | /categorys/add | 新增课程分类 | 否 | CategoryAddDTO | R{data: R} |  |
| GET | /categorys/all | 获取所有的课程分类信息，只包含id,名称，课程分类关系 | 否 | - | R{data: R%C2%ABList%C2%ABSimpleCategoryVO%C2%BB%C2%BB} |  |
| PUT | /categorys/disableOrEnable | 课程分类停用或启用 | 否 | CategoryDisableOrEnableDTO | R{data: R} |  |
| GET | /categorys/getAllOfOneLevel | 获取所有的课程分类，不分层 | 否 | - | R{data: R%C2%ABList%C2%ABCategoryVO%C2%BB%C2%BB} |  |
| GET | /categorys/list | 查询课程分类信息 | 否 | - | R{data: R%C2%ABList%C2%ABCategoryVO%C2%BB%C2%BB} |  |
| PUT | /categorys/update | 更新课程分类 | 否 | CategoryUpdateDTO | R{data: R} |  |
| DELETE | /categorys/{id} | 删除分类信息 | 否 | - | R{data: R} |  |
| GET | /categorys/{id} | 获取课程分类信息 | 否 | - | R{data: R%C2%ABCategoryInfoVO%C2%BB} |  |
| GET | /course/infoByTeacherIds | 通过老师id获取老师负责的课程和出的题目数量 | 否 | - | R{data: R%C2%ABList%C2%AB%E8%80%81%E5%B8%88id%E5%92%8C%E8%80%81%E5%B8%88%E5%AF%B9%E5%BA%94%E7%9A%84%E8%AF%BE%E7%A8%8B%E6%95%B0%EF%BC%8C%E5%87%BA%E9%A2%98%E6%95%B0%C2%BB%C2%BB} | internal |
| GET | /course/media/useInfo | mediaUserInfo | 否 | - | R{data: R%C2%ABList%C2%AB%E5%AA%92%E8%B5%84%E8%A2%AB%E5%BC%95%E7%94%A8%E6%83%85%E5%86%B5%C2%BB%C2%BB} | internal |
| GET | /course/name | queryCoursesIdByName | 否 | - | R{data: R%C2%ABList%C2%ABlong%C2%BB%C2%BB} | internal |
| GET | /course/section/{id} | sectionInfo | 否 | - | R{data: R%C2%AB%E5%B0%8F%E8%8A%82%E4%BF%A1%E6%81%AF%EF%BC%8C%E5%8C%85%E5%90%AB%E8%AF%BE%E7%A8%8Bid%E5%92%8C%E5%AA%92%E8%B5%84id%C2%BB} | internal |
| GET | /course/{id} | 获取课程信息 | 否 | - | R{data: R%C2%ABCourseFullInfoDTO%C2%BB} | internal |
| GET | /course/{id}/searchInfo | 课程上架时，需要查询课程信息，加入索引库 | 否 | - | R{data: R%C2%AB%E8%AF%BE%E7%A8%8B%E4%BF%A1%E6%81%AF%C2%BB} | internal |
| POST | /courses/baseInfo/save | 保存课程基本信息 | 否 | CourseBaseInfoSaveDTO | R{data: R%C2%ABCourseSaveVO%C2%BB} |  |
| GET | /courses/baseInfo/{id} | 获取课程基础信息 | 否 | - | R{data: R%C2%ABCourseBaseInfoVO%C2%BB} |  |
| GET | /courses/catas/index/list/{id} | 根据课程id，查询所有章节的序号 | 否 | - | R{data: R%C2%ABList%C2%ABCataSimpleInfoVO%C2%BB%C2%BB} |  |
| POST | /courses/catas/save/{id}/{step} | 保存章节 | 否 | JSON | R{data: R} |  |
| GET | /courses/catas/{id} | 获取课程的章节 | 否 | - | R{data: R%C2%ABList%C2%AB%E8%AF%BE%E7%A8%8B%E7%9B%AE%E5%BD%95%C2%BB%C2%BB} |  |
| GET | /courses/checkBeforeUpShelf/{id} | 课程上架前校验 | 否 | - | R{data: R} |  |
| GET | /courses/checkName | 校验课程名称是否已经存在 | 否 | - | R{data: R%C2%ABNameExistVO%C2%BB} |  |
| DELETE | /courses/delete/{id} | 课程删除 | 否 | - | R{data: R} |  |
| POST | /courses/down | 处理指定课程下架失败的问题 | 否 | - | R{data: R} |  |
| POST | /courses/downShelf | 课程下架 | 否 | CourseIdDTO | R{data: R} |  |
| GET | /courses/generator | 生成练习id | 否 | - | R{data: R%C2%ABCourseCataIdVO%C2%BB} |  |
| POST | /courses/media/save/{id} | 课程视频 | 否 | JSON | R{data: R} |  |
| GET | /courses/page | 管理端课程搜索接口 | 否 | - | R{data: R%C2%ABPageDTO%C2%ABCoursePageVO%C2%BB%C2%BB} |  |
| GET | /courses/portal | 用户端课程搜索接口 | 否 | - | R{data: R%C2%ABPageDTO%C2%ABCourseVO%C2%BB%C2%BB} |  |
| GET | /courses/simpleInfo/list | 根根条件列表获取课程信息 | 否 | - | R{data: R%C2%ABList%C2%ABCourseSimpleInfoDTO%C2%BB%C2%BB} |  |
| GET | /courses/subjects/get/{id} | 获取小节或练习中的题目（用于编辑） | 否 | - | R{data: R%C2%ABList%C2%ABCataSimpleSubjectVO%C2%BB%C2%BB} |  |
| POST | /courses/subjects/save/{id} | 保存小节或练习中的题目 | 否 | JSON | R{data: R} |  |
| POST | /courses/teachers/save | 保存老师信息 | 否 | %E8%AF%BE%E7%A8%8B%E8%80%81%E5%B8%88%E5%85%B3%E7%B3%BB%E6%A8%A1%E5%9E%8B | R{data: R} |  |
| GET | /courses/teachers/{id} | 查询课程相关的老师信息 | 否 | - | R{data: R%C2%ABList%C2%AB%E8%80%81%E5%B8%88%E8%AF%BE%E7%A8%8B%E4%BF%A1%E6%81%AF%C2%BB%C2%BB} |  |
| POST | /courses/up | 处理指定课程上架失败的问题 | 否 | - | R{data: R} |  |
| POST | /courses/upShelf | 课程上架 | 否 | CourseIdDTO | R{data: R} |  |
| GET | /courses/{id}/catalogs | 查询课程基本信息、目录、学习进度 | 否 | - | R{data: R%C2%ABCourseAndSectionVO%C2%BB} |  |

## 统一约定（引用全局规范）
- 响应格式：`pkg/response.R{Code,Msg,RequestId,Data any}`
- 分页：`PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`
- 错误码：见 [共享错误码表](../03-shared/error-codes.md)