package errorx

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/grpc/status"
)

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

// grpcMsgRe 匹配 *Error.Error() 的输出形态：
//
//	code=<num> msg=<msg>[ cause=<cause>]
//
// RPC 服务端返回的 *Error 经 zrpc/gRPC 传输后在客户端是 status error，
// 其 Message() 即服务端 Error() 的字符串；这里把它还原为 *Error，
// 使业务错误码能够跨服务传递（API 层 response.Fail 优先走这条路径）。
var grpcMsgRe = regexp.MustCompile(`^code=(\d+) msg=(.*)$`)

// FromGRPC 尝试把 gRPC/zrpc 传输来的错误还原为 *Error。
// 非跨服务错误（本地 status error 等）返回 ok=false。
func FromGRPC(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	st, ok := status.FromError(err)
	if !ok {
		return nil, false
	}
	m := grpcMsgRe.FindStringSubmatch(st.Message())
	if m == nil {
		return nil, false
	}
	code, err2 := strconv.ParseInt(m[1], 10, 32)
	if err2 != nil {
		return nil, false
	}
	msg := m[2]
	// cause 仅服务端日志用，不透传给客户端
	if i := strings.Index(msg, " cause="); i >= 0 {
		msg = msg[:i]
	}
	return New(Code(code), msg), true
}

// IsNotFound 判断 err 是否为业务"资源不存在"错误（Code == CodeNotFound）。
//
// 判定顺序：
//  1. 先用 FromGRPC 还原跨 gRPC/zrpc 传输的业务错误——zrpc 客户端拿到的是
//     status error，其 Message 即服务端 *Error.Error() 的 code=... msg=... 文本；
//  2. 再沿错误链（Unwrap）取第一个命中的 *Error（与 CodeOf 语义一致），
//     判断其 Code 是否为 CodeNotFound，可穿透 fmt.Errorf("%w") 等包装。
//
// 注意：model.ErrNotFound（sqlx 裸错误 sql.ErrNoRows 的包装）不是 *Error，
// 本函数对其返回 false；数据层错误请继续用 errors.Is(err, model.ErrNotFound) 判断。
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := FromGRPC(err); ok {
		return e.Code == CodeNotFound
	}
	var biz *Error
	if errors.As(err, &biz) {
		return biz.Code == CodeNotFound
	}
	return false
}
