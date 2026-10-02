package xerr

import (
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/grpc/status"
)

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
