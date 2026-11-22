// Package errorx 定义统一的业务错误码与错误类型。
//
// 所有微服务统一通过该包返回业务错误：逻辑层返回 *Error，
// handler 层通过 response.Write 将错误渲染为统一的 R 响应结构。
package errorx

// Code 业务错误码。
// 采用与 HTTP 语义对齐的整数码段，避免客户端再做一层映射；
// 服务内部需要更细粒度错误时，可在 Code 上扩展专用码段。
type Code int32

// 通用业务错误码
const (
	CodeSuccess            Code = 200 // 成功
	CodeBadRequest         Code = 400 // 请求参数错误
	CodeUnauthorized       Code = 401 // 未认证或凭证失效
	CodeForbidden          Code = 403 // 无权限
	CodeNotFound           Code = 404 // 资源不存在
	CodeConflict           Code = 409 // 状态冲突（重复操作、状态不允许等）
	CodeTooManyRequests    Code = 429 // 请求过于频繁
	CodeInternal           Code = 500 // 内部错误
	CodeServiceUnavailable Code = 503 // 服务不可用
)

// 常用业务提示文案
const (
	MsgSuccess            = "OK"
	MsgBadRequest         = "请求参数错误"
	MsgUnauthorized       = "未登录或登录已过期"
	MsgForbidden          = "无权访问"
	MsgNotFound           = "资源不存在"
	MsgConflict           = "操作冲突，请刷新后重试"
	MsgTooManyRequests    = "请求过于频繁，请稍后重试"
	MsgInternal           = "系统繁忙，请稍后重试"
	MsgServiceUnavailable = "服务暂不可用"
)
