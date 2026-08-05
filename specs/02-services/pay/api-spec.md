# pay 服务 HTTP API 接口规格

> 来源：`docs/tjxt.openapi.json` (聚合文档) | 最后同步：2026-08-05

## 接口清单

| 方法 | 路径 | 摘要 | 认证 | 请求体 | 响应体 | 权限标签 |
|------|------|------|------|--------|--------|----------|
| POST | /pay-channels | 添加支付渠道 | 否 | PayChannelDTO | R{data: R%C2%ABlong%C2%BB} |  |
| GET | /pay-channels/list | 查询支付渠道列表 | 否 | - | R{data: R%C2%ABList%C2%ABPayChannelDTO%C2%BB%C2%BB} |  |
| PUT | /pay-channels/{id} | 修改支付渠道 | 否 | PayChannelDTO | R{data: R} |  |
| POST | /pay-orders | 扫码支付申请支付单，返回支付url地址，用于生产二维码 | 否 | PayApplyDTO | R{data: R%C2%ABstring%C2%BB} |  |
| GET | /pay-orders/{bizOrderId}/status | 根据业务端订单id查询支付结果 | 否 | - | R{data: R%C2%ABPayResultDTO%C2%BB} |  |
| GET | /pay/channels | 获取支付渠道列表接口 | 否 | - | R{data: R%C2%ABList%C2%ABPayChannelVO%C2%BB%C2%BB} |  |
| POST | /pay/order | 支付申请,返回支付二维码url | 否 | PayApplyFormDTO | R{data: R%C2%ABstring%C2%BB} |  |

## 统一约定（引用全局规范）
- 响应格式：`pkg/response.R{Code,Msg,RequestId,Data any}`
- 分页：`PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`
- 错误码：见 [共享错误码表](../03-shared/error-codes.md)