# trade（交易订单）Specification

## Purpose

trade（交易订单）微服务覆盖购物车、订单（预下单 / 下单 / 免费课报名 / 取消 / 删除 / 查询）、订单明细与退款申请审批流，并承载迁移自 pay 服务的支付渠道 / 支付 / 退款业务能力（经 `pay.rpc` 代理，本服务不落地支付流水）。HTTP 服务 `trade-api` 监听 8809，RPC 服务 `trade.rpc` 监听 8089（gRPC），经 etcd（`127.0.0.1:2379`，key `trade.rpc`）注册与发现；数据库为 MySQL `tj_trade` 库（cart / order / order_detail / refund_apply / undo_log 共 5 张表）。下游依赖 `pay.rpc`（已装配）以及 course / promotion（课程快照取价与优惠券试算，设计意图依赖，RPC 客户端配置当前未接线）；作为生产者向 RabbitMQ `order.exchange` 发布 `order.pay` / `order.refund` 交易域事件供 learning 消费开通 / 撤销课程（生产端配置就绪，当前代码无发布调用）。73 个 logic（RPC 37 + API 36）均为真实实现，两层 `go build` 编译通过，全仓无 goctl 占位 stub。

## Requirements

### Requirement: HTTP 购物车接口

trade-api SHALL 提供以下购物车维护接口（全部经 `TradeRpc` 转发）：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /carts | 获取购物车中的课程，响应 `R{data: R«List«CartVO»»}` |
| POST | /carts | 添加课程到购物车，请求体 `CartsAddDTO`，响应 `R{data: R}` |
| DELETE | /carts/{id} | 删除指定的购物车条目 |
| DELETE | /carts | 批量删除购物车条目 |

补充约束：

- 响应统一为 `pkg/response.R{Code,Msg,RequestId,Data any}`；分页约定 `PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`；错误码遵循共享错误码规范（trade 域 103000-103999）
- 按 configs.md 记载，trade-api 全部路由声明 `@server (jwt: Auth)`，`Auth.AccessSecret` 须与 auth.rpc 签发密钥一致（聚合文档 `docs/tjxt.openapi.json` 将上述端点认证列标注为「否」，以配置记载为准）
- 聚合文档仅收录 12 个 cart/order 端点；trade-api 实际实现 36 个 API logic，支付发起 / 支付渠道管理 / 退款申请审批等 HTTP 转发逻辑的路由未收录于聚合文档

#### Scenario: 添加课程到购物车

- **WHEN** 学员携带 `CartsAddDTO` 请求 POST /carts
- **THEN** API 层转发 `TradeRpc.CartAdd{CourseId}`，返回 `NamePlaceVO{Existed:true, Message:"ok"}`
- **AND** 同一用户对同一课程重复加购不重复插入，仍返回成功（幂等）

#### Scenario: 查询购物车列表

- **WHEN** 学员请求 GET /carts
- **THEN** 返回当前用户的购物车课程列表（CartVO 含加购快照价 price 与当前售价 now_price，单位分）

### Requirement: HTTP 下单接口

trade-api SHALL 提供预下单、下单与免费课报名接口：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /orders/prePlaceOrder | 预下单接口，确认订单可用优惠券信息，响应 `R{data: R«OrderConfirmVO»}` |
| POST | /orders/placeOrder | 下单接口，请求体 `PlaceOrderDTO`，响应 `R{data: R«PlaceOrderResultVO»}` |
| POST | /orders/freeCourse/{courseId} | 免费课立刻报名接口，响应 `R{data: R«PlaceOrderResultVO»}` |

> 已知缺口：预下单实际不生成 / 落库订单 id（`OrderConfirmVO.OrderId` 恒为 0），优惠券试算列表恒为空（promotion 未接线）；聚合文档路径摘要中「生成订单id」为设计意图描述。

#### Scenario: 正式下单返回支付截止时间

- **WHEN** 学员请求 POST /orders/placeOrder
- **THEN** 返回 `PlaceOrderResultVO`（OrderId / PayAmount / Status=1 / PayOutTime=now+15min）
- **AND** 订单与逐课明细落库后订单为「待支付」状态

#### Scenario: 免费课立刻报名

- **WHEN** 学员请求 POST /orders/freeCourse/{courseId}
- **THEN** 直接生成已支付订单（status=2）与明细，金额 0，跳过支付环节

### Requirement: HTTP 订单查询与生命周期接口

trade-api SHALL 提供订单查询与生命周期管理接口：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /orders/page | 分页查询我的订单，响应 `R{data: R«PageDTO«OrderPageVO»»}` |
| GET | /orders/{id} | 根据 id 查询订单详细信息，响应 `R{data: R«OrderVO»}` |
| GET | /orders/{id}/status | 查询订单支付状态，响应 `R{data: R«PlaceOrderResultVO»}` |
| PUT | /orders/{id}/cancel | 取消订单接口 |
| DELETE | /orders/{id} | 删除订单接口 |

#### Scenario: 轮询支付状态

- **WHEN** 客户端在支付后轮询 GET /orders/{id}/status
- **THEN** 返回 `PlaceOrderResultVO`（status：1待支付 2已支付 3已关闭 4已完成 5已报名 6申请退款）

#### Scenario: 删除订单为逻辑删除

- **WHEN** 用户请求 DELETE /orders/{id}
- **THEN** 订单置 `deleted=1` 逻辑删除，不做物理删除

### Requirement: RPC 服务契约

`Trade` 服务 SHALL 经 etcd 注册（key `trade.rpc`，监听 `0.0.0.0:8089`）对外提供 37 个 RPC 方法，按 7 个业务域分组：

| 业务域 | 方法 | 数量 |
|--------|------|------|
| 购物车 | CartAdd、CartList、CartGet、CartUpdate、CartDelete、CartBatchDelete | 6 |
| 订单 | OrderPrePlace、OrderPlace、OrderFreeCourse、OrderPageQuery、OrderGet、OrderStatus、OrderCancel、OrderDelete | 8 |
| 订单明细 | OrderDetailGet、OrderDetailCourseCheck、OrderDetailEnrollCourse、OrderDetailEnrollNum、OrderDetailPageQuery、OrderDetailPurchaseInfo | 6 |
| 支付渠道（管理端） | PayChannelAdd、PayChannelList、PayChannelGet、PayChannelDelete | 4 |
| 支付（学员侧） | PayApply、PayResultQuery、PayChannels | 3 |
| 退款（业务侧） | RefundApply、RefundResultQuery | 2 |
| 退款申请（审批流） | RefundApplyCreate、RefundApplyApprove、RefundApplyCancel、RefundApplyDetail、RefundApplyNext、RefundApplyPageQuery、RefundApplyGet、RefundApplyDelete | 8 |

关键字段约定（金额字段单位均为分）：

- `PlaceOrderResultVO`：order_id、pay_amount、pay_out_time、status（1待支付 2已支付 3已关闭 4已完成 5已报名 6申请退款）
- `OrderVO`：total_amount / real_amount / discount_amount、status / status_desc / message、coupon_desc、details（订单明细列表）、progress_nodes（订单进度节点 name / desc / status / time）
- `CartVO`：course_id、course_name、cover_url、price（加购时单价）、now_price（当前售价）、expired（课程是否已下架/失效）
- `PayResultDTO.status`：1 支付中 2 失败 3 成功；`RefundResultDTO.status`：1 退款中 2 失败 3 成功
- `ApproveRequest.approve_type`：1 同意、2 拒绝；`PayChannelDTO.channel_code`：wxPay / aliPay
- `OrderDetailPageRequest` 支持分页排序与 no_no（订单号）、id、mobile、status、refund_status、pay_channel、下单时间区间筛选；`RefundApplyPageRequest` 支持 order_detail_id / order_id / refund_status / mobile / 申请时间区间筛选

消费关系：当前仅 trade-api 自身消费该 RPC（HTTP Handler → `tradeclient.Trade`）；对 learning 的课程开通不走 RPC 直连，而是经 RabbitMQ 发布事件（见「交易域 MQ 事件发布」）。

#### Scenario: 加购到下单链路

- **WHEN** 用户从加购走到正式下单
- **THEN** 依次调用 CartAdd 加入购物车 → CartList 查看 → OrderPrePlace 预下单（返回试算信息与可用优惠券）→ OrderPlace 正式下单

#### Scenario: 扫码支付链路

- **WHEN** 用户对订单发起支付
- **THEN** PayChannels 拉取可选渠道 → PayApply 返回 qr_url 扫码 → 轮询 PayResultQuery / OrderStatus 直到 status=3（成功）

#### Scenario: 退款审批链路

- **WHEN** 学员提交退款申请并由管理员处理
- **THEN** RefundApplyCreate → RefundApplyNext 取件 → RefundApplyApprove（approve_type=1 同意）→ 渠道退款由独立 RefundApply RPC 触发 → 轮询 RefundResultQuery

#### Scenario: 管理端运营看板

- **WHEN** 管理端查询运营数据
- **THEN** OrderDetailPageQuery 分页查明细 → OrderDetailPurchaseInfo 查课程销售汇总 → OrderDetailEnrollNum 查各课程报名人数

### Requirement: 购物车管理规则

购物车 SHALL 按 `user_id` 隔离，并在加购时写入课程快照以避免每次列表回查课程服务：

- 按用户隔离：列表 / 删除均带 `user_id` 条件，禁止跨用户
- 课程快照冗余：加购时写入课程 name / cover / price（`cart.cover_url` / `course_name` / `price`），列表不再回查
- 幂等去重：同 userId+courseId 已存在仍返回成功，不重复插入
- 批量删除：支持一次多 id（`CartBatchDeleteRequest.ids`）

> 已知缺口：`toCartVO` 的 `NowPrice` 恒等于加购快照价、`Expired` 恒为 false，未与课程服务当前价比对 / 标记下架（设计意图要求列表展示时比对现价并标记 expired，缺口 #8）。
> 注（源冲突）：business-rules.md 记载加购经 `CourseRpc.CourseSimpleInfoList` 写入快照；configs.md「已知配置缺口」则记载 trade.yaml 未声明 course 服务连接配置，两处表述存在出入。

#### Scenario: 幂等加购

- **WHEN** 同一用户对已在其购物车中的课程再次加购
- **THEN** 不重复插入购物车条目，接口仍返回成功

#### Scenario: 加购写入课程快照

- **WHEN** 加购成功
- **THEN** 课程名称 / 封面 / 价格以快照形式写入 cart 表
- **AND** 后续购物车列表查询不再回查课程服务

#### Scenario: 删除归属校验

- **WHEN** 用户删除购物车条目（CartDelete）
- **THEN** 删除操作带 user 归属校验，禁止跨用户删除

### Requirement: 订单下单规则

下单 SHALL 采用两段式入口（PrePlace 试算 + Place 落库），一单多课拆分明细：

- 两段式入口：`OrderPrePlace` 仅试算不落库，返回 `OrderConfirmVO{OrderId:0, TotalAmount, Courses, Discounts:空}`；`OrderPlace` 真正生成订单与明细
- 金额口径：当前 `real_amount = total_amount`、`discount_amount = 0`（优惠券未接入），金额单位分
- 一单多课：order(1) → order_detail(N)，每课一条明细，独立 `real_pay_amount` 与 `refund_status`
- 免费课直通：金额 0，直接置 status=2 已支付（非设计意图的「已报名 5」）
- 支付超时：下单返回 `PayOutTime = now+15min`（关单动作依赖支付回调，当前未落地）

> 已知缺口：预下单不生成 / 落库 order_id，`OrderPlace` 自行生成雪花 ID；设计意图的「预下单 order_id 回传做幂等」未落地，下单本身无 order_id 幂等保护，幂等仅依赖购物车去重（缺口 #2）。

#### Scenario: 预下单仅试算

- **WHEN** 调用 OrderPrePlace 携带 courseIds
- **THEN** 不落库，返回试算总价与课程列表
- **AND** OrderId 恒为 0、Discounts 恒为空

#### Scenario: OrderPlace 下单流程（实际实现）

- **WHEN** 用户提交下单（courseIds 非空，userId 取自 JWT）
- **THEN** 批量取课程快照，累加 total_amount（= real_amount，discount=0）
- **AND** 以雪花 nextID() 生成 order.id 并 Insert order(status=1)，逐课 Insert order_detail(status=1, real_pay_amount=price)
- **AND** 返回 OrderId / PayAmount(total) / Status=1 / PayOutTime=now+15min

#### Scenario: 支付超时时间下发

- **WHEN** 下单成功
- **THEN** 返回支付截止时间 PayOutTime = 当前时间 + 15 分钟
- **AND** 超时关单动作依赖支付回调链路，trade 内当前未实现

### Requirement: 订单状态机

`order.status` 与 `order_detail.status` SHALL 使用同一状态枚举并同步流转：

| 值 | 含义 | 当前写入路径 |
|----|------|-------------|
| 1 | 待支付 | 下单写入（OrderPlace） |
| 2 | 已支付 | 仅 OrderFreeCourse 写入；支付回调未落地 |
| 3 | 已关闭 | OrderCancel（用户取消） |
| 4 | 已完成 | 支付后 30 天 — 无写入路径 |
| 5 | 已报名 | 课程开通 — 无写入路径 |
| 6 | 已申请退款 | 仅 order 表定义（order_detail 注释止于 5）— 无写入路径 |

约束规则：

- 取消仅限待支付：`OrderCancel` 仅 status=1 时允许，否则拒绝
- 删除为逻辑删除：`OrderDelete` 置 `deleted=1`，不物理删除
- 明细状态无 6：退款流转走 `order_detail.refund_status` 独立字段（1~6）
- 进度节点：`buildOrderProgressNodes` 按 order.status 组装「提交订单→支付成功/关闭→已完成/退款中」

> 已知缺口：status 2/4/5/6 在 trade 内无写入路径（支付回调、关单定时任务、退款联动均未落地），当前仅 1（下单）、3（取消）由 trade 自身驱动（缺口 #9）。

#### Scenario: 取消仅限待支付

- **WHEN** OrderCancel 作用于 status=1（待支付）订单
- **THEN** 更新订单状态为 status=3（已关闭）并写入 message
- **AND** status≠1 时取消请求被拒绝

#### Scenario: 状态写入路径核查

- **WHEN** 审查 trade 代码的订单状态写入
- **THEN** 仅 status=1（OrderPlace）与 status=3（OrderCancel）由常规下单流程驱动，status=2 仅由 OrderFreeCourse 写入，4/5/6 无写入路径

#### Scenario: 订单进度节点组装

- **WHEN** 查询订单详情（OrderGet）
- **THEN** `buildOrderProgressNodes` 按 order.status 组装进度节点「提交订单→支付成功/关闭→已完成/退款中」

### Requirement: 支付委托规则

trade SHALL 不落地支付流水，支付发起 / 结果查询 / 渠道管理全部经 `PayRpc payclient.Pay` 代理 pay 服务：

- 支付发起：`PayApply` → `PayRpc.ApplyPayOrder{BizUserId:order.UserId, BizOrderNo:order.Id, Amount:order.TotalAmount, PayChannelCode, PayType:4（native 扫码）}` → 返回 QrUrl
- 支付结果：`PayResultQuery` → `PayRpc.QueryPayResult` → `PayResultDTO`
- 渠道管理：`PayChannelAdd` / `PayChannelList` / `PayChannelGet` / `PayChannelDelete` 全部委托 PayRpc（无本地 `pay_channel` 表）
- 渠道展示：`PayChannels` / `PayChannelList` → `PayRpc.ListPayChannels` 映射 VO，按 `channel_priority` 排序由客户端处理
- 渠道删除：`PayChannelDelete` 委托 `UpdatePayChannelStatus{Status:2}` 软删

> 已知缺口：设计意图要求 PayApply 成功后回写 `order.pay_channel` / `order.pay_order_no` 并以 `real_amount` 作为支付金额；实际 PayApply 不回写 order，且金额取自 TotalAmount（缺口 #3）。

#### Scenario: 发起扫码支付

- **WHEN** PayApply 对已存在的订单发起支付
- **THEN** 经 `PayRpc.ApplyPayOrder` 下单（PayType=4 native 扫码），返回支付二维码 URL

#### Scenario: 查询支付结果

- **WHEN** PayResultQuery 按 biz_order_id 查询支付结果
- **THEN** 返回 `PayResultDTO`（status：1 支付中 2 失败 3 成功）

#### Scenario: 支付渠道管理委托

- **WHEN** 管理端新增 / 查询 / 删除支付渠道
- **THEN** 全部委托 PayRpc 处理，trade 本地无 pay_channel 表
- **AND** 删除为软删（`UpdatePayChannelStatus{Status:2}`）

### Requirement: 退款申请审批流规则

退款 SHALL 分「业务审批」（refund_apply 表）与「渠道退款」（pay 服务）两阶段，`refund_apply.status` 与 `order_detail.refund_status` 双写同步（枚举一致，1~6）：

| 阶段 | 方法 | refund_apply.status | 实际行为 |
|------|------|---------------------|---------|
| 学员提交 | RefundApplyCreate | 1 待审批 | Insert(status=1) + `UpdateRefundStatus(detail,1)` 双写 |
| 学员撤回 | RefundApplyCancel | 2 取消退款 | status→2 + 同步 detail refund_status=2 |
| 管理员同意 | RefundApplyApprove（approve_type=1） | 3 同意退款 | `UpdateApprove`（approver 硬编码 0）+ `UpdateRefundStatus(detail,3)` |
| 管理员拒绝 | RefundApplyApprove（approve_type=2） | 4 拒绝退款 | 同上，status→4 |
| 渠道退款 | RefundApply（独立 RPC） | — | 委托 `PayRpc.ApplyRefund`，返回 RefundResultDTO |
| 渠道结果 | RefundResultQuery | — | 委托 `PayRpc.QueryRefundResult` |

退款状态枚举（`refund_apply.status` / `order_detail.refund_status`）：1 待审批、2 取消退款、3 同意退款、4 拒绝退款、5 退款成功、6 退款失败。

约束规则：

- 申请粒度为明细：`refund_apply.order_detail_id`，一单多课可单独退某门
- 可退判定：`toOrderDetailItemVO.can_refund` 仅当明细 status∈{2,4,5} 且 refund_status 为空/0/1
- 待办取件：`RefundApplyNext` 经 `FindNextPending` 取 status=1 的下一条申请
- 审批留痕：approver / approve_opinion / remark / approve_time 落库（approver 当前恒为 0）

> 已知缺口（重大）：`RefundApplyApprove` 是业务审批与渠道退款的衔接点，但实际仅置 refund_apply.status=3/4 与 order_detail.refund_status，不调用 `PayRpc.ApplyRefund`、不置 `order.status=6`、不发射 MQ；审批→渠道退款链路在 trade 内断开，渠道退款须由独立 `RefundApply` RPC 另行触发（缺口 #4）。审批人 approver 硬编码为 0（缺口 #5）。

#### Scenario: 提交退款申请双写

- **WHEN** 学员对订单明细提交退款申请
- **THEN** Insert refund_apply(status=1)，并同步更新 order_detail.refund_status=1

#### Scenario: 审批同意流程（实际实现）

- **WHEN** RefundApplyApprove 以 approve_type=1 审批
- **THEN** FindOne refund_apply 校验存在后置 status=3
- **AND** UpdateApprove(id, status, approver=0, opinion, remark, approveTime=now)，并 UpdateRefundStatus(detail, 3)
- **AND** 不触发 PayRpc 调用、不动 order 状态、不发射 MQ

#### Scenario: 可退判定

- **WHEN** 组装订单明细 VO
- **THEN** can_refund 仅当明细 status∈{2,4,5} 且 refund_status 为空/0/1 时为真

#### Scenario: 待办取件

- **WHEN** 管理员调用 RefundApplyNext 取下一条待审批申请
- **THEN** 经 `FindNextPending` 返回 status=1 的下一条退款申请

### Requirement: 交易域 MQ 事件发布

trade SHALL 作为事件生产者向 RabbitMQ `order.exchange` 交换机发布支付成功事件（路由键 `order.pay`）与退款事件（路由键 `order.refund`），由 learning 消费开通 / 撤销课程（消费端队列 `learning.lesson.pay.queue` / `learning.lesson.refund.queue`；生产端不声明队列）。`MQProducer` 初始化失败时仅记录日志 `init rabbitmq producer failed, will skip event publish` 并保持 nil，服务照常启动，业务代码发布事件前必须判空。

> 已知缺口：trade 全部 logic 均未调用 `MQProducer.Publish`，`order.pay` / `order.refund` 事件无发射点，learning 的 MQ 消费端无生产者数据来源（缺口 #6）。

#### Scenario: Producer 初始化容错

- **WHEN** RabbitMQ 连接失败导致 `mq.NewProducer` 出错
- **THEN** 仅打印 `init rabbitmq producer failed, will skip event publish` 日志，`MQProducer` 保持 nil，服务照常启动
- **AND** 后续业务代码发布事件前必须判空跳过

#### Scenario: 事件拓扑一致性

- **WHEN** 配置 trade 生产端与 learning 消费端
- **THEN** 两端交换机（`order.exchange`）与路由键（`order.pay` / `order.refund`）必须一致
- **AND** 生产端不声明队列，队列由 learning 消费端声明

### Requirement: 数据模型与不变量

trade 数据库 `tj_trade` SHALL 包含 5 张表：

| 表 | 用途 | 关键约束 / 不变量 |
|----|------|------------------|
| cart | 购物车条目 | cover_url / course_name / price 为加购时点的课程快照冗余字段 |
| order | 订单 | 表名为 MySQL 保留字，SQL 中必须用反引号包裹；status 默认 1；total_amount / real_amount / discount_amount 单位分（discount 默认 0）；coupon_ids 为 json；pay_order_no 可空 |
| order_detail | 订单明细 | 索引 idx_order(order_id)、idx_user_course(user_id, course_id)（校验是否已购）、idx_course_expire_time（扫描到期课程）、idx_pay_channel（按支付渠道统计）；valid_duration 单位月，空则永久有效；course_expire_time 支付成功开始计时 |
| refund_apply | 退款申请 | status 默认 1；粒度为 order_detail_id；pay_order_no / refund_order_no 可空 |
| undo_log | Seata AT 回滚日志 | 唯一键 ux_undo_log(xid)；字符集 utf8（区别于其余表 utf8mb4），遗留自 Java 版本，Go 侧未生成 model |

关系与跨库引用（均无外键）：

- cart (N) —course_id→ course 域（跨库，无外键）
- order (1) — (N) order_detail；order_detail — (1) refund_apply（经 order_detail_id / order_id）
- order.pay_order_no → tj_pay.pay_order；refund_apply.refund_order_no → tj_pay.refund_order（跨库，无外键）
- `*_gen.go` 由 goctl 生成禁止修改，扩展方法统一放同名自定义 .go 文件（如 ordermodel.go 扩展 OrderModel）

> 已知缺口：tj_trade 库不含 pay_channel / pay_order / refund_order 表（定义在 `sql/ddl/tj_pay.sql`，trade 经 PayRpc 代理调用 pay 服务，本库不落地支付流水）；undo_log 无对应 model，未接入分布式事务框架。
> 注（源冲突）：data-model.md 记载 4 个自定义 Model（cartmodel.go / ordermodel.go / orderdetailmodel.go / refundapplymodel.go）均为 goctl 空壳、仅 Insert / FindOne / Update / Delete 四个基础方法；business-rules.md（v1.1）则记载上述 Model 已补齐分页、聚合、按 user_id 查询、状态更新、双写等扩展方法且编译通过。本规格按 data-model.md 表述记录并如实加注。

#### Scenario: order 表保留字

- **WHEN** 编写涉及订单表的 SQL
- **THEN** 表名必须用反引号包裹（`order` 为 MySQL 保留字）

#### Scenario: 金额字段单位

- **WHEN** 读写订单 / 明细 / 退款金额字段（total_amount、real_amount、discount_amount、price、real_pay_amount、refund_amount）
- **THEN** 单位均为分

#### Scenario: 跨库引用支付流水

- **WHEN** 需要关联支付 / 退款流水
- **THEN** 经 order.pay_order_no、refund_apply.refund_order_no 跨库引用 tj_pay.pay_order / tj_pay.refund_order（无外键约束）
- **AND** tj_trade 本库不落地支付流水

### Requirement: 服务配置

trade SHALL 按以下配置部署运行：

| 配置项 | 默认值 |
|--------|--------|
| trade-api Name / Port | trade-api / 8809（HTTP，监听 0.0.0.0） |
| trade.rpc Name / ListenOn | trade.rpc / 0.0.0.0:8089（gRPC） |
| etcd | 127.0.0.1:2379，服务注册 / 发现 key `trade.rpc` |
| DataSource | root:0000@tcp(127.0.0.1:3306)/tj_trade?charset=utf8mb4&parseTime=true&loc=Local |
| Redis Cache | 127.0.0.1:6379，node 单机模式，密码为空 |
| RabbitMQ | 127.0.0.1:5672（rabbitmq/rabbitmq），Exchange `order.exchange`，PayRoutingKey `order.pay`，RefundRoutingKey `order.refund` |
| PayRpc | etcd 127.0.0.1:2379，key `pay.rpc` |
| Auth | AccessSecret 默认 `change-me-in-production`，AccessExpire 7200 秒 |

约束规则：

- `Auth.AccessSecret` 必须与签发方 auth.rpc 的 `Jwt.AccessSecret` 一致（trade-api 只校验 token 不签发）；生产环境必须修改默认值，否则 JWT 可被伪造
- Telemetry / Prometheus / Log 三段可观测性注入配置已注入，勿删（统一约定见 `specs/infra/observability/spec.md`）

> 已知配置缺口：trade-api.yaml 仅声明 TradeRpc、trade.yaml 仅声明 PayRpc——预下单所需课程信息与优惠券试算（OrderConfirmVO.courses / discounts）需补 CourseRpc / PromotionRpc 客户端配置或在 RPC 层补相应依赖；当前无 course 服务连接配置。trade 无 MQ 消费端配置（仅 Producer 侧 Exchange + RoutingKey），符合「trade 只生产不消费」的定位。

#### Scenario: JWT 密钥一致性

- **WHEN** trade-api 校验用户 token
- **THEN** `Auth.AccessSecret` 必须与 auth.rpc 签发密钥一致，否则 token 校验失败
- **AND** 生产环境必须修改默认值 `change-me-in-production`，否则 JWT 可被伪造

#### Scenario: 端口与服务发现

- **WHEN** 部署 trade 服务
- **THEN** trade-api 监听 8809（HTTP）、trade.rpc 监听 8089（gRPC）
- **AND** 经 etcd key `trade.rpc` 注册，调用方通过该 key 服务发现

