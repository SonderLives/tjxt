# trade 服务 HTTP API 接口规格

> 来源：`docs/tjxt.openapi.json` (聚合文档) | 最后同步：2026-08-05

## 接口清单

| 方法 | 路径 | 摘要 | 认证 | 请求体 | 响应体 | 权限标签 |
|------|------|------|------|--------|--------|----------|
| DELETE | /carts | 批量删除购物车条目 | 否 | - | R{data: R} |  |
| GET | /carts | 获取购物车中的课程 | 否 | - | R{data: R%C2%ABList%C2%ABCartVO%C2%BB%C2%BB} |  |
| POST | /carts | 添加课程到购物车 | 否 | CartsAddDTO | R{data: R} |  |
| DELETE | /carts/{id} | 删除指定的购物车条目 | 否 | - | R{data: R} |  |
| POST | /orders/freeCourse/{courseId} | 免费课立刻报名接口 | 否 | - | R{data: R%C2%ABPlaceOrderResultVO%C2%BB} |  |
| GET | /orders/page | 分页查询我的订单 | 否 | - | R{data: R%C2%ABPageDTO%C2%ABOrderPageVO%C2%BB%C2%BB} |  |
| POST | /orders/placeOrder | 下单接口 | 否 | PlaceOrderDTO | R{data: R%C2%ABPlaceOrderResultVO%C2%BB} |  |
| GET | /orders/prePlaceOrder | 预下单接口，生成订单id，确认订单可用优惠券信息 | 否 | - | R{data: R%C2%ABOrderConfirmVO%C2%BB} |  |
| DELETE | /orders/{id} | 删除订单接口 | 否 | - | R{data: R} |  |
| GET | /orders/{id} | 根据id查询订单详细信息 | 否 | - | R{data: R%C2%ABOrderVO%C2%BB} |  |
| PUT | /orders/{id}/cancel | 取消订单接口 | 否 | - | R{data: R} |  |
| GET | /orders/{id}/status | 查询订单支付状态 | 否 | - | R{data: R%C2%ABPlaceOrderResultVO%C2%BB} |  |

## 统一约定（引用全局规范）
- 响应格式：`pkg/response.R{Code,Msg,RequestId,Data any}`
- 分页：`PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`
- 错误码：见 [共享错误码表](../03-shared/error-codes.md)