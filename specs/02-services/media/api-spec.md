# media 服务 HTTP API 接口规格

> 来源：`docs/tjxt.openapi.json` (聚合文档) | 最后同步：2026-08-05

## 接口清单

| 方法 | 路径 | 摘要 | 认证 | 请求体 | 响应体 | 权限标签 |
|------|------|------|------|--------|--------|----------|
| POST | /files | 上传文件 | 否 | JSON | R{data: R%C2%ABFileDTO%C2%BB} |  |
| DELETE | /files/{id} | 删除文件 | 否 | - | R{data: R} |  |
| GET | /files/{id} | 获取文件信息 | 否 | - | R{data: R%C2%ABFileDTO%C2%BB} |  |
| DELETE | /medias | 批量删除媒资视频 | 否 | - | R{data: R} |  |
| GET | /medias | 分页搜索已上传媒资信息 | 否 | - | R{data: R%C2%ABPageDTO%C2%ABMediaVO%C2%BB%C2%BB} |  |
| POST | /medias | 上传视频后保存媒资信息 | 否 | MediaUploadResultDTO | R{data: R%C2%ABMediaDTO%C2%BB} |  |
| GET | /medias/signature/play | 获取播放视频的授权签名 | 否 | - | R{data: R%C2%ABVideoPlayVO%C2%BB} |  |
| GET | /medias/signature/preview | 管理端获取预览视频的授权签名 | 否 | - | R{data: R%C2%ABVideoPlayVO%C2%BB} |  |
| GET | /medias/signature/upload | 获取上传视频的授权签名 | 否 | - | R{data: R%C2%ABstring%C2%BB} |  |
| DELETE | /medias/{mediaId} | 删除媒资视频 | 否 | - | R{data: R} |  |

## 统一约定（引用全局规范）
- 响应格式：`pkg/response.R{Code,Msg,RequestId,Data any}`
- 分页：`PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`
- 错误码：见 [共享错误码表](../03-shared/error-codes.md)