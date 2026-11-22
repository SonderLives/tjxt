package errorx

import (
	"fmt"
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

// un 解包一层 error，避免直接依赖 errors.Unwrap 的语义。
func un(err error) error {
	type unwrapper interface{ Unwrap() error }
	if u, ok := err.(unwrapper); ok {
		return u.Unwrap()
	}
	return nil
}
