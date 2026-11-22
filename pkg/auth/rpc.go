package auth

import (
	"context"
	"strconv"

	"tjxt/pkg/xerr"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// 用户身份的 gRPC 透传链路：
//
//	JWT 只在 API 进程由 rest.WithJwt 注入 ctx，RPC 进程的 ctx 里没有 userId，
//	需要经 gRPC metadata 跨进程传递。完整链路为：
//
//	API 侧 rest.WithJwt 注入 ctx（userId claim，见 context.go）
//	  → API 调 RPC 时由 OutgoingUserIdInterceptor 把 userId 写入出站 metadata
//	  → RPC 侧用 UserIdFromIncomingMD 从入站 metadata 读取。
//
// RPC 服务之间互调时，若也需要透传身份，同样在下一段 zrpc 客户端上挂
// OutgoingUserIdInterceptor 即可（其 ctx 通常先经 UserIdFromIncomingMD 恢复出 userId）。
//
// 现状说明：exam 与 message 服务此前各自在 logic/common.go 里定义了私有的
// metadataKeyUserID = "user_id"，exam 的 api 侧还在 servicecontext.go 手写了
// 同功能的 userIdClientInterceptor；本文件将该约定收敛到 pkg/auth 统一维护，
// 键名保持 "user_id" 不变，与既有透传链路完全兼容。
const MetadataUserIDKey = "user_id"

// OutgoingUserIdInterceptor gRPC 一元调用客户端拦截器：将当前登录用户 id
// 写入出站 metadata（键为 MetadataUserIDKey），供下游 RPC 服务读取。
//
// 用法：zrpc.MustNewClient(c.XxxRpc, zrpc.WithUnaryClientInterceptor(auth.OutgoingUserIdInterceptor))。
//
// UserIdFromCtx 失败说明当前请求未登录（如游客可访问的接口），
// 属正常业务场景：静默跳过注入、原样透传请求，不吞后续调用的错误。
func OutgoingUserIdInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if userID, err := UserIdFromCtx(ctx); err == nil {
		ctx = metadata.AppendToOutgoingContext(ctx, MetadataUserIDKey, strconv.FormatInt(userID, 10))
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

// UserIdFromIncomingMD 从 gRPC 入站 metadata 读取上游透传的登录用户 id。
//
// 缺少 metadata、缺少键、无法解析为 int64 或值 <=0 时，
// 统一返回 xerr.Unauthorized("")（与 UserIdFromCtx 的失败语义一致）。
func UserIdFromIncomingMD(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, xerr.Unauthorized("")
	}
	vals := md.Get(MetadataUserIDKey)
	if len(vals) == 0 {
		return 0, xerr.Unauthorized("")
	}
	userID, err := strconv.ParseInt(vals[0], 10, 64)
	if err != nil || userID <= 0 {
		return 0, xerr.Unauthorized("")
	}
	return userID, nil
}
