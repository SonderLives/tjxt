package svc

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"tjxt/api/exam/internal/config"
	examclient "tjxt/rpc/exam/client/exam"
	"tjxt/pkg/auth"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ServiceContext struct {
	Config  config.Config
	ExamRpc examclient.Exam
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		ExamRpc: examclient.NewExam(zrpc.MustNewClient(c.ExamRpc,
			zrpc.WithUnaryClientInterceptor(userIdClientInterceptor),
			zrpc.WithUnaryClientInterceptor(errClientInterceptor))),
	}
}

// userIdClientInterceptor 将当前登录用户 id 透传到下游 RPC 的 gRPC metadata 中，
// 供 exam-rpc 侧写入题目 creater/updater 字段。
func userIdClientInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if userId, err := auth.UserIdFromCtx(ctx); err == nil {
		ctx = metadata.AppendToOutgoingContext(ctx, "user_id", strconv.FormatInt(userId, 10))
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

// errClientInterceptor 将下游 RPC 返回的错误还原为 xerr 业务错误，
// 避免 *errorx.Error 跨 gRPC 传输后丢失业务错误码（go-zero 默认不还原）。
func errClientInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err == nil {
		return nil
	}
	if _, ok := err.(*errorx.Error); ok {
		return err
	}
	// 下游通过 xerr 返回的错误，其 Error() 文本会原样进入 gRPC status message，
	// 格式为 code=<码> msg=<提示> [cause=<底层原因>]
	if statusErr, ok := status.FromError(err); ok && statusErr.Code() != 0 {
		if code, msg := parseXerrMessage(statusErr.Message()); code != 0 {
			return errorx.New(code, msg)
		}
	}
	return err
}

// xerrMessagePattern 匹配 errorx.Error() 的字符串格式
var xerrMessagePattern = regexp.MustCompile(`^code=(\d+)\s+msg=(.*)$`)

// parseXerrMessage 从 gRPC status message 中还原业务错误码与提示信息
func parseXerrMessage(msg string) (errorx.Code, string) {
	m := xerrMessagePattern.FindStringSubmatch(msg)
	if len(m) != 3 {
		return 0, ""
	}
	code, err := strconv.Atoi(m[1])
	if err != nil || code <= 0 {
		return 0, ""
	}
	text := m[2]
	if idx := strings.Index(text, " cause="); idx >= 0 {
		text = text[:idx]
	}
	return errorx.Code(code), text
}
