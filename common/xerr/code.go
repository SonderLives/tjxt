// Package xerr 定义统一的业务错误码与错误类型。
//
// 所有微服务统一通过该包返回业务错误：逻辑层返回 *Error，
// handler 层通过 result.Write 将错误渲染为统一的 R 响应结构。
package xerr

import "fmt"

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

// Error 携带业务错误码的错误，实现了 error 接口。
type Error struct {
	Code Code   // 业务错误码
	Msg  string // 面向用户的提示信息
	err  error  // 底层原因（仅用于日志/排查，不对外暴露）
}

// New 创建业务错误。
func New(code Code, msg string) *Error {
	return &Error{Code: code, Msg: msg}
}

// Newf 创建业务错误并格式化提示信息。
func Newf(code Code, format string, args ...any) *Error {
	return &Error{Code: code, Msg: fmt.Sprintf(format, args...)}
}

// Wrap 包装底层错误为业务错误，msg 面向用户，err 供日志排查。
func Wrap(err error, code Code, msg string) *Error {
	if err == nil {
		return nil
	}
	return &Error{Code: code, Msg: msg, err: err}
}

// Wrapf 包装底层错误并格式化提示信息。
func Wrapf(err error, code Code, format string, args ...any) *Error {
	if err == nil {
		return nil
	}
	return &Error{Code: code, Msg: fmt.Sprintf(format, args...), err: err}
}

// 常见错误快捷构造

func BadRequestf(format string, args ...any) *Error {
	return Newf(CodeBadRequest, format, args...)
}

func Unauthorized(msg string) *Error {
	if msg == "" {
		msg = MsgUnauthorized
	}
	return New(CodeUnauthorized, msg)
}

func Forbidden(msg string) *Error {
	if msg == "" {
		msg = MsgForbidden
	}
	return New(CodeForbidden, msg)
}

func NotFound(msg string) *Error {
	if msg == "" {
		msg = MsgNotFound
	}
	return New(CodeNotFound, msg)
}

func Conflict(msg string) *Error {
	if msg == "" {
		msg = MsgConflict
	}
	return New(CodeConflict, msg)
}

func Internal(msg string) *Error {
	if msg == "" {
		msg = MsgInternal
	}
	return New(CodeInternal, msg)
}

func ServiceUnavailable(msg string) *Error {
	if msg == "" {
		msg = MsgServiceUnavailable
	}
	return New(CodeServiceUnavailable, msg)
}

// Error 实现 error 接口。
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.err != nil {
		return fmt.Sprintf("code=%d msg=%s cause=%v", e.Code, e.Msg, e.err)
	}
	return fmt.Sprintf("code=%d msg=%s", e.Code, e.Msg)
}

// Unwrap 暴露底层原因，便于 errors.Is / errors.As 链式判断。
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

// Is 让 *Error 与 *Error、以及目标错误被包装的情况可用 errors.Is 判断。
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	if e == nil || t == nil {
		return e == t
	}
	return e.Code == t.Code && e.Msg == t.Msg
}

// CodeOf 沿错误链提取业务错误码；非业务错误返回 CodeInternal。
func CodeOf(err error) Code {
	for err != nil {
		if e, ok := err.(*Error); ok {
			return e.Code
		}
		err = un(err)
	}
	return CodeInternal
}

// MsgOf 沿错误链提取业务提示信息；非业务错误返回内部错误文案。
func MsgOf(err error) string {
	for err != nil {
		if e, ok := err.(*Error); ok {
			return e.Msg
		}
		err = un(err)
	}
	return MsgInternal
}

// HttpStatus 将业务错误码映射为 HTTP 状态码。
func HttpStatus(err error) int {
	code := CodeOf(err)
	switch code {
	case CodeSuccess:
		return 200
	case CodeBadRequest:
		return 400
	case CodeUnauthorized:
		return 401
	case CodeForbidden:
		return 403
	case CodeNotFound:
		return 404
	case CodeConflict:
		return 409
	case CodeTooManyRequests:
		return 429
	case CodeServiceUnavailable:
		return 503
	default:
		return 500
	}
}

// un 解包一层 error，避免直接依赖 errors.Unwrap 的语义。
func un(err error) error {
	type unwrapper interface{ Unwrap() error }
	if u, ok := err.(unwrapper); ok {
		return u.Unwrap()
	}
	return nil
}
