# trade 服务 logic 实现契约（单一事实源）

> 目标：把 `apps/trade/rpc/internal/logic/` 的 37 个 + `apps/trade/api/internal/logic/` 的 36 个 goctl 占位 stub 实现为真实业务逻辑，并 `go build ./...` 通过。
> 子代理只需**改写函数体**（替换 `// todo: add your logic here and delete this line` + 保留 `return`/构造返回），**不要改动函数签名、struct、NewXxx、import 块里除新增外的已有 import**。

## §0 总则（强制）

1. 保留原始 `func (l *XxxLogic) Xxx(...)` 签名与返回值形式（命名返回 `resp` 也要保留）。
2. 子代理之间**不运行 `go build`**（并发冲突）。最后由主控统一编译修复。
3. 每个文件只改自己负责的文件，互不覆盖。
4. 错误统一用 `tjxt/pkg/xerr`：`xerr.BadRequestf`、`xerr.NotFound`、`xerr.New(code, msg)`、`xerr.Wrap(err, msg)`。
   - 导入：`"tjxt/pkg/xerr"`。
   - 常见 code：`xerr.CodeBadRequest`、`xerr.CodeNotFound`、`xerr.CodeConflict`、`xerr.CodeInternal`。
5. 金额、id、status、数量一律用 **int64**；proto 里 `int32` 字段在构造 pb 时显式 `int32(x)`；API types 里都是 int64，无需转换。
6. 自定义 Model 内嵌 `CachedConn`，**统一使用** `m.QueryRowsNoCacheCtx` / `m.QueryRowNoCacheCtx` / `m.ExecNoCacheCtx`（不要直接用 `QueryRowsCtx`，被遮蔽会编译失败）。
7. 不要引入 `interface{}`，用 `any`。
8. 不要在 logic 里直接连 DB 写裸 SQL；一律走 `l.svcCtx.XxxModel` 的自定义方法（见 §2）。

## §1 RPC common.go 已有 helper（rpc/internal/logic/common.go）

状态常量：
- `OrderStatusPending=1 / Paid=2 / Closed=3 / Finished=4 / Enrolled=5 / Refunding=6`
- `DetailStatusPending=1 / Paid=2 / Closed=3 / Finished=4 / Enrolled=5`
- `RefundStatusPending=1 / Cancel=2 / Approve=3 / Reject=4 / Success=5 / Failed=6`

描述函数：`orderStatusDesc(int64) string`、`detailStatusDesc(int64) string`、`refundStatusDesc(int64) string`
工具：`nextID() int64`、`now() time.Time`、`formatTime(t time.Time) string`、`formatNullTime(sql.NullTime) string`、`nullInt64Value(sql.NullInt64) int64`、`nullStringValue(sql.NullString) string`、`calcPages(total,pageSize int64) int64`
课程回填：`fetchCourseMap(ctx, svcCtx, ids) map[int64]*courseclient.CourseSimpleInfoItem`、`courseName(m,id)`、`courseCover(m,id)`、`coursePrice(m,id)`
VO 构建：`toCartVO(*model.Cart) *pb.CartVO`、`toOrderDetailItemVO(*model.OrderDetail) *pb.OrderDetailItemVO`、`toOrderVO(*model.Order,[]*model.OrderDetail) *pb.OrderVO`、`toOrderDetailAdminVO(*model.OrderDetail,*model.Order,*model.RefundApply) *pb.OrderDetailAdminVO`、`toRefundApplyVO(*model.RefundApply,*model.Order,*model.OrderDetail) *pb.RefundApplyVO`、`toRefundApplyPageVO(*model.RefundApply) *pb.RefundApplyPageVO`

> RPC 层 **没有 CourseRpc 客户端**（svc 里只有 PayRpc / CourseRpc 在 RPC svc 已接，但 RPC logic 只用 Model + PayRpc）。课程名/封面/价格由 RPC 层直接读 `model.OrderDetail.Name/CoverUrl/Price` 字段（下单时已落库），**不需要**调 CourseRpc。简单购物车 CartAdd 才用 CourseRpc 校验课程是否存在。

## §2 Model 方法（已就绪，直接调用）

`l.svcCtx.CartModel / OrderModel / OrderDetailModel / RefundApplyModel`

- CartModel: `Insert(ctx,*Cart) (sql.Result,error)`、`FindOne(ctx,id) (*Cart,error)`、`Delete(ctx,id) error`、`Update(ctx,*Cart) error`、`ListByUserId(ctx,userId) ([]*Cart,error)`、`FindByUserIdAndCourseId(ctx,userId,courseId) (*Cart,error)`
- OrderModel: `Insert(ctx,*Order) (sql.Result,error)`、`FindOne(ctx,id) (*Order,error)`、`Delete(ctx,id) error`、`Update(ctx,*Order) error`、`PageQueryByUser(ctx,userId,pageNo,pageSize,status,noNo int64, sortBy string, isAsc bool) ([]*Order,int64,error)`、`UpdateStatus(ctx,id,status int64,message string) error`、`MarkPaid(ctx,id,payOrderNo int64,payChannel string,payTime time.Time,realAmount int64) error`、`MarkClosed(ctx,id,message string) error`
- OrderDetailModel: `Insert(ctx,*OrderDetail) (sql.Result,error)`、`FindOne(ctx,id) (*OrderDetail,error)`、`Delete(ctx,id) error`、`Update(ctx,*OrderDetail) error`、`ListByOrderId(ctx,orderId) ([]*OrderDetail,error)`、`ListByCourseId(ctx,courseId) ([]*OrderDetail,error)`、`ListByCourseIds(ctx,[]int64) ([]*OrderDetail,error)`、`CountPaidByCourseIds(ctx,[]int64) (map[int64]int64,error)`、`StatByCourseId(ctx,courseId) (enrollNum,realPayAmount,refundNum int64,err error)`、`CountPaidByUserIds(ctx,[]int64) (map[int64]int64,error)`、`UpdateStatus(ctx,id,status int64) error`、`UpdateRefundStatus(ctx,id,refundStatus int64) error`、`PageQuery(ctx,OrderDetailPageFilter) ([]*OrderDetail,int64,error)`  (filter 字段见 model 文件)
- RefundApplyModel: `Insert(ctx,*RefundApply) (sql.Result,error)`、`FindOne(ctx,id) (*RefundApply,error)`、`Delete(ctx,id) error`、`Update(ctx,*RefundApply) error`、`ListByUserId(ctx,userId) ([]*RefundApply,error)`、`FindByOrderDetailId(ctx,orderDetailId) (*RefundApply,error)`、`ListByOrderId(ctx,orderId) ([]*RefundApply,error)`、`FindNextPending(ctx) (*RefundApply,error)`、`PageQuery(ctx,RefundApplyPageFilter) ([]*RefundApply,int64,error)`、`UpdateStatus(ctx,id,status int64) error`、`UpdateApprove(ctx,id,status,approver int64,opinion,remark string,approveTime sql.NullTime,refundOrderNo,payOrderNo int64) error`

> `model.ErrNotFound` 是"未找到"错误（import `"github.com/zeromicro/go-zero/core/stores/sqlc"` 的 `sqlc.ErrNotFound` 也可用于判定：在自定义 Model 方法里已转成 `model.ErrNotFound`）。

## §3 svc 字段

RPC (`rpc/internal/svc`): `l.svcCtx.CartModel / OrderModel / OrderDetailModel / RefundApplyModel / PayRpc (payclient.Pay) / CourseRpc (courseclient.Course) / MQProducer (*mq.Producer, 可能 nil)`
API (`api/internal/svc`): `l.svcCtx.TradeRpc (tradeclient.Trade)` —— API 层**只调 TradeRpc**，不碰 Model。

## §4 API 层约定（api/internal/logic）

1. user_id 从 JWT 取：`userId, err := auth.UserIdFromCtx(l.ctx)`（import `"tjxt/pkg/auth"`）。需要鉴权的接口先取；公开统计类（如 PurchaseInfo 用 courseId 查询）可不取。
2. 调用 RPC：`reply, err := l.svcCtx.TradeRpc.Xxx(l.ctx, &pb.XxxRequest{...})`，pb 类型用别名 `tradeclient.XxxRequest`（已在 `tradeclient` 包里 `= pb.XxxRequest`）。
   - 导入：`"tjxt/apps/trade/rpc/pb"` 以及 `"tjxt/apps/trade/rpc/trade"`（别名 `tradeclient`）。实际上 handler 已 import，logic 只要 import pb 和 tradeclient。
3. 把 `pb.Reply` 映射到 `*types.XxxResp` / `[]types.Xxx` 后返回 `(resp, err)`。**handler 已经 `result.Write(w,r,resp,err)`**，logic 只管返回。
4. 无请求体的接口（CartList / PayChannels / RefundApplyNext / PayChannelList）：方法签名**没有 req 参数**，直接 `l.svcCtx.TradeRpc.Xxx(l.ctx, &pb.Empty{})`。
5. 返回 `NamePlaceVO` 的接口（增删改）：成功返回 `&types.NamePlaceVO{Existed:true, Message:"ok"}`，失败 `return nil, err`。
6. pb↔types 字段同名对应（见 `types.go` 与 `trade.proto`）。嵌套类型需逐层映射（如 OrderVO.Details 是 `[]OrderDetailItemVO`，OrderVO.ProgressNodes 是 `[]OrderProgressNodeVO`）。务必对照两个文件的字段名精确赋值，int64 直接传，string 直接传。

## §5 已知缺口（logic 层无需处理，按"本地填默认值"实现即可）

- `PayRpc` 真实实现依赖 pay 服务与第三方（微信/支付宝）。logic 里调用 `l.svcCtx.PayRpc.Xxx(...)` 即可；若 pay 服务未启动，运行期会报错，但**编译**不受影响（接口已存在）。
- **支付渠道无本地表**：trade 没有 `pay_channel` 表，所有 `PayChannel*` / `PayChannels` RPC **委托给 PayRpc**（已接，类型别名 `payclient`，import `"tjxt/apps/pay/rpc/pay"`）：
  - `PayChannelAdd(in *pb.PayChannelDTO)` → `PayRpc.AddPayChannel(ctx, &payclient.PayChannelRequest{Name, ChannelCode, ChannelIcon, ChannelPriority})` → `&pb.IdReply{Id: resp.Id}`
  - `PayChannelList(in *pb.Empty)` → `PayRpc.ListPayChannels(ctx, &payclient.ListPayChannelsRequest{})` → 遍历 `resp.List` 映射为 `[]*pb.PayChannelDTO{Id,Name,ChannelCode,ChannelIcon,ChannelPriority,Status}`
  - `PayChannelGet(in *pb.IdRequest)` → `PayRpc.ListPayChannels(ctx, ...)` 然后按 `id` 找到匹配项返回；找不到 `return nil, xerr.NotFound("支付渠道不存在")`
  - `PayChannelDelete(in *pb.IdRequest)` → `PayRpc.UpdatePayChannelStatus(ctx, &payclient.UpdatePayChannelStatusRequest{Id: in.Id, Status: 2})`（2=禁用，等价删除）
  - `PayChannels(in *pb.Empty)`（学员侧）→ `PayRpc.ListPayChannels` → 映射为 `[]*pb.PayChannelVO{Id,Name,ChannelCode,ChannelIcon,ChannelPriority}`
  - pay 侧 `PayChannelResponse` 字段：`id,name,channel_code,channel_priority(int32),channel_icon,status(int32)`；`PayChannelRequest` 字段：`id,name,channel_code,channel_priority(int32),channel_icon`。
- **PayApply**：`PayApplyRequest{OrderId, PayChannelCode}` 不足以直接调 `PayRpc.ApplyPayOrder`（其需 `BizUserId,Amount,PayType,NotifyUrl`）。实现：先用 `l.svcCtx.OrderModel.FindOne(ctx, in.OrderId)` 取订单的 `UserId`/`TotalAmount`（字段名见 ordermodel_gen.go），再 `PayRpc.ApplyPayOrder(ctx, &payclient.ApplyPayOrderRequest{BizUserId: order.UserId, BizOrderNo: in.OrderId, Amount: order.TotalAmount, PayChannelCode: in.PayChannelCode, PayType: 4, NotifyUrl: ""})`，返回 `&pb.PayApplyReply{QrUrl: resp.QrCodeUrl}`。
- **PayResultQuery** / **RefundResultQuery** / **RefundApply**：直接委托 `PayRpc.QueryPayResult` / `PayRpc.QueryRefundResult` / `PayRpc.ApplyRefund`，按字段映射返回（pay 侧 `PayResultResponse{ pay_order_no, biz_order_no, status }`）。字段名以 `apps/pay/rpc/pay.proto` 为准。
- `MQProducer` 可能 nil：发布支付/退款事件时先 `if l.svcCtx.MQProducer != nil { ... }`，避免 nil panic。
- `MQProducer` 可能 nil：发布支付/退款事件时先 `if l.svcCtx.MQProducer != nil { ... }`，避免 nil panic。
- 优惠券（Coupon）无服务接入：`OrderPrePlace` 的 discounts 返回空切片，`total_amount` 直接等于课程价格之和；`OrderPlace` 的 coupon_ids 暂忽略。
- 订单 `OrderStatus`/`OrderCancel` 等状态流转：直接 Model.UpdateStatus，不依赖外部事件。
- 报名（EnrollCourse）：`OrderDetailEnrollCourse` 仅做计数统计（经 `OrderDetailModel.CountPaidByUserIds`），不真正开通课程（无 learning 写接口）。
- `RefundApply`/`RefundResultQuery`（业务侧退款）调 `PayRpc` 同名方法；无 pay 时编译通过即可。

## §6 硬性陷阱

- proto 字段 `no_no`（OrderPageRequest）语义是"订单号精确过滤"，用 `id` 精确匹配。
- `OrderDetailPageRequest` 的 mobile 因无 user 服务接入，图方便**忽略 mobile 过滤**（Model.PageQuery 也未实现 mobile）。
- `RefundApplyVO` 的 `Price` 用 `detail.RealPayAmount`；`StudentName`/`Mobile` 无 user 服务，留空字符串。
- 金额单位：数据库里 `price`/`real_pay_amount` 等已是分（int64）。直接透传。
- 不要在 RPC logic 里 import `tjxt/apps/trade/api/...`（反向依赖禁止）。
- **所有表 `id` 均为 `bigint NOT NULL` 且无 `AUTO_INCREMENT`**：插入时必须 `Id: nextID()`（来自 common.go）。订单(`order`)/明细(`order_detail`)/退款(`refund_apply`) 表还有 `creater`/`updater`（`bigint NOT NULL`），插入时设为当前 `userId`（来自 `auth.UserIdFromCtx`）；`deleted` 默认 0。子表（order_detail）的 `order_id` 用父表插入时分配的 `order.Id`（同一函数内 `nextID()` 生成后持有）。
- `order` / `order_detail` 结构字段名见 `ordermodel_gen.go` / `orderdetailmodel_gen.go`（如 `Creater`/`Updater`/`Status`/`TotalAmount`/`RealAmount`/`CourseId`/`Name`/`CoverUrl`/`Price`/`RealPayAmount` 等）。写 Insert 前先读对应的 `*_gen.go` 确认必填字段。

## §7 文件分配（见各子代理任务描述）

RPC 组：cart(6)+pay(7) / order(8)+order_detail(6) / refund(10)
API 组：cart(6)+pay(7) / order(8)+order_detail(6) / refund(10)

黄金标准（先读这两个文件对齐风格）：
- `apps/trade/rpc/internal/logic/cartaddlogic.go`（RPC 含 CourseRpc 校验 + Model.Insert + xerr）
- `apps/trade/api/internal/logic/cartaddlogic.go`（API 取 user_id + 调 TradeRpc + 返回 NamePlaceVO）
