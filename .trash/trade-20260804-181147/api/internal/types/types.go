package types

// ============ 请求 ============

// PageRequest 通用分页请求
type PageRequest struct {
	PageNo   int64  `form:"pageNo,omitempty"`
	PageSize int64  `form:"pageSize,omitempty"`
	IsAsc    bool   `form:"isAsc,omitempty"`
	SortBy   string `form:"sortBy,omitempty"`
}

// CartsAddReq 添加购物车请求
type CartsAddReq struct {
	CourseId int64 `json:"courseId"`
}

// CartsDeleteReq 删除购物车请求
type CartsDeleteReq struct {
	Ids string `form:"ids"`
}

// PrePlaceOrderReq 预下单请求
type PrePlaceOrderReq struct {
	CourseIds string `form:"courseIds"`
}

// PlaceOrderReq 提交订单请求
type PlaceOrderReq struct {
	OrderId   int64   `json:"orderId"`
	CourseIds []int64 `json:"courseIds"`
	CouponIds []int64 `json:"couponIds,omitempty"`
}

// FreeCourseReq 免费课报名请求
type FreeCourseReq struct {
	CourseId int64 `path:"courseId"`
}

// OrderPageReq 订单分页请求
type OrderPageReq struct {
	PageRequest
	Status int64 `form:"status,omitempty"`
}

// OrderIdReq 路径订单 id 请求
type OrderIdReq struct {
	Id int64 `path:"id"`
}

// OrderDetailCourseReq 订单明细课程 id 请求
type OrderDetailCourseReq struct {
	Id int64 `path:"id"`
}

// EnrollCourseReq 学员报名数请求
type EnrollCourseReq struct {
	StudentIds string `form:"studentIds"`
}

// EnrollNumReq 课程报名数请求
type EnrollNumReq struct {
	CourseIdList string `form:"courseIdList"`
}

// OrderDetailPageReq 订单明细分页请求
type OrderDetailPageReq struct {
	PageRequest
	Id             int64  `form:"id,omitempty"`
	Mobile         string `form:"mobile,omitempty"`
	Status         int64  `form:"status,omitempty"`
	RefundStatus   int64  `form:"refundStatus,omitempty"`
	PayChannel     string `form:"payChannel,omitempty"`
	OrderStartTime string `form:"orderStartTime,omitempty"`
	OrderEndTime   string `form:"orderEndTime,omitempty"`
}

// PurchaseInfoReq 课程购买信息请求
type PurchaseInfoReq struct {
	CourseId int64 `form:"courseId"`
}

// PayApplyReq 支付申请请求
type PayApplyReq struct {
	OrderId        int64  `json:"orderId"`
	PayChannelCode string `json:"payChannelCode"`
}

// RefundApplyReq 退款申请请求
type RefundApplyReq struct {
	OrderDetailId int64  `json:"orderDetailId"`
	RefundReason  string `json:"refundReason"`
	QuestionDesc  string `json:"questionDesc,omitempty"`
}

// ApproveReq 退款审批请求
type ApproveReq struct {
	Id             int64  `json:"id"`
	ApproveType    int64  `json:"approveType"`
	ApproveOpinion string `json:"approveOpinion,omitempty"`
	Remark         string `json:"remark,omitempty"`
}

// RefundCancelReq 取消退款请求
type RefundCancelReq struct {
	Id            int64 `json:"id"`
	OrderDetailId int64 `json:"orderDetailId"`
}

// RefundApplyPageReq 退款申请分页请求
type RefundApplyPageReq struct {
	PageRequest
	Id             int64  `form:"id,omitempty"`
	Mobile         string `form:"mobile,omitempty"`
	OrderId        int64  `form:"orderId,omitempty"`
	OrderDetailId  int64  `form:"orderDetailId,omitempty"`
	RefundStatus   int64  `form:"refundStatus,omitempty"`
	ApplyStartTime string `form:"applyStartTime,omitempty"`
	ApplyEndTime   string `form:"applyEndTime,omitempty"`
}

// RefundIdReq 退款 id 请求
type RefundIdReq struct {
	Id int64 `path:"id"`
}

// ============ 响应 ============

// CartVO 购物车条目
type CartVO struct {
	Id         int64  `json:"id"`
	CourseId   int64  `json:"courseId"`
	CourseName string `json:"courseName"`
	CoverUrl   string `json:"coverUrl"`
	Price      int64  `json:"price"`
	NowPrice   int64  `json:"nowPrice"`
	Expired    bool   `json:"expired"`
}

// OrderCourseVO 订单确认页课程信息
type OrderCourseVO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	CoverUrl string `json:"coverUrl"`
	Price    int64  `json:"price"`
}

// CouponDiscountDTO 优惠券抵扣信息
type CouponDiscountDTO struct {
	DiscountAmount int64   `json:"discountAmount"`
	DiscountDetail any     `json:"discountDetail"`
	Ids            []int64 `json:"ids"`
	Rules          []any   `json:"rules"`
}

// OrderConfirmVO 预下单返回
type OrderConfirmVO struct {
	OrderId     int64               `json:"orderId"`
	Courses     []OrderCourseVO     `json:"courses"`
	Discounts   []CouponDiscountDTO `json:"discounts"`
	TotalAmount int64               `json:"totalAmount"`
}

// PlaceOrderResultVO 下单/查询支付状态返回
type PlaceOrderResultVO struct {
	OrderId    int64  `json:"orderId"`
	PayAmount  int64  `json:"payAmount"`
	PayOutTime string `json:"payOutTime"`
	Status     int64  `json:"status"`
}

// OrderProgressNodeVO 订单进度节点
type OrderProgressNodeVO struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

// OrderDetailVO 订单明细（用户端）
type OrderDetailVO struct {
	Id            int64  `json:"id"`
	OrderId       int64  `json:"orderId"`
	CourseId      int64  `json:"courseId"`
	Name          string `json:"name"`
	CoverUrl      string `json:"coverUrl"`
	Price         int64  `json:"price"`
	RealPayAmount int64  `json:"realPayAmount"`
	RefundStatus  int64  `json:"refundStatus"`
	CouponDesc    string `json:"couponDesc"`
	CanRefund     bool   `json:"canRefund"`
}

// OrderPageVO 订单分页条目
type OrderPageVO struct {
	Id          int64           `json:"id"`
	Status      int64           `json:"status"`
	StatusDesc  string          `json:"statusDesc"`
	TotalAmount int64           `json:"totalAmount"`
	RealAmount  int64           `json:"realAmount"`
	CreateTime  string          `json:"createTime"`
	Details     []OrderDetailVO `json:"details"`
}

// OrderVO 订单详情
type OrderVO struct {
	Id             int64                 `json:"id"`
	Status         int64                 `json:"status"`
	StatusDesc     string                `json:"statusDesc"`
	Message        string                `json:"message"`
	TotalAmount    int64                 `json:"totalAmount"`
	RealAmount     int64                 `json:"realAmount"`
	DiscountAmount int64                 `json:"discountAmount"`
	CouponDesc     string                `json:"couponDesc"`
	CreateTime     string                `json:"createTime"`
	Details        []OrderDetailVO       `json:"details"`
	ProgressNodes  []OrderProgressNodeVO `json:"progressNodes"`
}

// OrderDetailPageVO 订单明细分页条目（管理端）
type OrderDetailPageVO struct {
	Id               int64  `json:"id"`
	OrderId          int64  `json:"orderId"`
	Name             string `json:"name"`
	Mobile           string `json:"mobile"`
	Price            int64  `json:"price"`
	RealPayAmount    int64  `json:"realPayAmount"`
	PayChannel       string `json:"payChannel"`
	Status           int64  `json:"status"`
	StatusDesc       string `json:"statusDesc"`
	RefundStatus     int64  `json:"refundStatus"`
	RefundStatusDesc string `json:"refundStatusDesc"`
	CreateTime       string `json:"createTime"`
}

// OrderDetailAdminVO 订单明细详情（管理端）
type OrderDetailAdminVO struct {
	Id                 int64                 `json:"id"`
	OrderId            int64                 `json:"orderId"`
	Name               string                `json:"name"`
	StudentName        string                `json:"studentName"`
	Mobile             string                `json:"mobile"`
	Price              int64                 `json:"price"`
	RealPayAmount      int64                 `json:"realPayAmount"`
	DiscountAmount     int64                 `json:"discountAmount"`
	PayChannel         string                `json:"payChannel"`
	PayOrderNo         int64                 `json:"payOrderNo"`
	Status             int64                 `json:"status"`
	RefundStatus       int64                 `json:"refundStatus"`
	RefundReason       string                `json:"refundReason"`
	RefundMessage      string                `json:"refundMessage"`
	RefundChannel      string                `json:"refundChannel"`
	RefundOrderNo      int64                 `json:"refundOrderNo"`
	RefundApplyId      int64                 `json:"refundApplyId"`
	RefundProposerName string                `json:"refundProposerName"`
	FailedReason       string                `json:"failedReason"`
	Message            string                `json:"message"`
	Remark             string                `json:"remark"`
	CouponDesc         string                `json:"couponDesc"`
	CanRefund          bool                  `json:"canRefund"`
	StudyValidTime     string                `json:"studyValidTime"`
	Nodes              []OrderProgressNodeVO `json:"nodes"`
}

// PayChannelVO 支付渠道
type PayChannelVO struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	ChannelCode     string `json:"channelCode"`
	ChannelPriority int64  `json:"channelPriority"`
	ChannelIcon     string `json:"channelIcon"`
}

// RefundApplyVO 退款申请详情
type RefundApplyVO struct {
	Id                 int64  `json:"id"`
	OrderId            int64  `json:"orderId"`
	OrderDetailId      int64  `json:"orderDetailId"`
	StudentName        string `json:"studentName"`
	Mobile             string `json:"mobile"`
	Name               string `json:"name"`
	Price              int64  `json:"price"`
	RealPayAmount      int64  `json:"realPayAmount"`
	DiscountAmount     int64  `json:"discountAmount"`
	PayChannel         string `json:"payChannel"`
	PayOrderNo         int64  `json:"payOrderNo"`
	PaySuccessTime     string `json:"paySuccessTime"`
	OrderTime          string `json:"orderTime"`
	RefundReason       string `json:"refundReason"`
	QuestionDesc       string `json:"questionDesc"`
	RefundProposerName string `json:"refundProposerName"`
	RefundChannel      string `json:"refundChannel"`
	RefundOrderNo      int64  `json:"refundOrderNo"`
	Status             int64  `json:"status"`
	ApproveOpinion     string `json:"approveOpinion"`
	ApproveTime        string `json:"approveTime"`
	FailedReason       string `json:"failedReason"`
	Message            string `json:"message"`
	Remark             string `json:"remark"`
	CouponDesc         string `json:"couponDesc"`
	CreateTime         string `json:"createTime"`
}

// RefundApplyPageVO 退款申请分页条目
type RefundApplyPageVO struct {
	Id                int64  `json:"id"`
	OrderId           int64  `json:"orderId"`
	OrderDetailId     int64  `json:"orderDetailId"`
	ProposerName      string `json:"proposerName"`
	ProposerMobile    string `json:"proposerMobile"`
	ApproverName      string `json:"approverName"`
	RefundAmount      int64  `json:"refundAmount"`
	Status            int64  `json:"status"`
	RefundStatusDesc  string `json:"refundStatusDesc"`
	RefundSuccessTime string `json:"refundSuccessTime"`
	CreateTime        string `json:"createTime"`
	ApproveTime       string `json:"approveTime"`
}

// PurchaseInfoVO 课程购买信息
type PurchaseInfoVO struct {
	EnrollNum     int64 `json:"enrollNum"`
	RealPayAmount int64 `json:"realPayAmount"`
	RefundNum     int64 `json:"refundNum"`
}

// Page 通用分页响应
type Page struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
	Pages int64 `json:"pages"`
}
