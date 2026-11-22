package logic

// 支付单状态（与 rpc/pay pay_order 表注释一致的语义镜像：
// rpc 侧常量在 internal 包内不可跨进程复用，此处仅列出 api 层用到的值）。
const (
	PayOrderStatusClosed  = 2 // 支付超时或取消
	PayOrderStatusSuccess = 3 // 支付成功
)
