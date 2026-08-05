# promotion 服务 HTTP API 接口规格

> 来源：`docs/tjxt.openapi.json` (聚合文档) | 最后同步：2026-08-05

## 接口清单

| 方法 | 路径 | 摘要 | 认证 | 请求体 | 响应体 | 权限标签 |
|------|------|------|------|--------|--------|----------|
| POST | /coupons | 新增优惠券接口 | 否 | CouponFormDTO | R{data: R} |  |
| GET | /coupons/list | 查询发放中的优惠券列表 | 否 | - | R{data: R%C2%ABList%C2%ABCouponVO%C2%BB%C2%BB} |  |
| GET | /coupons/page | 分页查询优惠券接口 | 否 | - | R{data: R%C2%ABPageDTO%C2%ABCouponPageVO%C2%BB%C2%BB} |  |
| DELETE | /coupons/{id} | 删除优惠券 | 否 | - | R{data: R} |  |
| GET | /coupons/{id} | 根据id查询优惠券接口 | 否 | - | R{data: R%C2%ABCouponDetailVO%C2%BB} |  |
| PUT | /coupons/{id}/issue | 发放优惠券接口 | 否 | CouponIssueFormDTO | R{data: R} |  |
| PUT | /coupons/{id}/pause | 暂停发放优惠券接口 | 否 | - | R{data: R} |  |
| POST | /user-coupons/available | 查询我的优惠券可用方案 | 否 | JSON | R{data: R%C2%ABList%C2%ABCouponDiscountDTO%C2%BB%C2%BB} |  |
| POST | /user-coupons/discount | 根据券方案计算订单优惠明细 | 否 | OrderCouponDTO | R{data: R%C2%ABCouponDiscountDTO%C2%BB} |  |
| GET | /user-coupons/page | 分页查询我的优惠券接口 | 否 | - | R{data: R%C2%ABPageDTO%C2%ABCouponVO%C2%BB%C2%BB} |  |
| PUT | /user-coupons/refund | 退还指定优惠券 | 否 | - | R{data: R} |  |
| GET | /user-coupons/rules | 分页查询我的优惠券接口 | 否 | - | R{data: R%C2%ABList%C2%ABstring%C2%BB%C2%BB} |  |
| PUT | /user-coupons/use | 核销指定优惠券 | 否 | - | R{data: R} |  |
| POST | /user-coupons/{code}/exchange | 兑换码兑换优惠券接口 | 否 | - | R{data: R} |  |
| POST | /user-coupons/{couponId}/receive | 领取优惠券接口 | 否 | - | R{data: R} |  |

## 统一约定（引用全局规范）
- 响应格式：`pkg/response.R{Code,Msg,RequestId,Data any}`
- 分页：`PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`
- 错误码：见 [共享错误码表](../03-shared/error-codes.md)