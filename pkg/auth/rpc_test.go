package auth

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"tjxt/common/errorx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// withUserClaim 模拟 rest.WithJwt 鉴权后注入的 ctx（userId claim，见 context.go）。
// go-zero 中间件注入 JWT claim 时用的就是 string key，因此这里必须沿用
// 同样的字符串 key 才能如实模拟，属有意为之。
func withUserClaim(ctx context.Context, userID any) context.Context {
	//nolint:staticcheck // go-zero rest.WithJwt 注入 claim 即用 string key，测试需保持一致
	return context.WithValue(ctx, keyUserID, userID)
}

// captureInvoker 记录拦截器传入的最终 ctx，模拟真实的 invoker 调用
func captureInvoker(captured *context.Context) grpc.UnaryInvoker {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		*captured = ctx
		return nil
	}
}

func TestOutgoingUserIdInterceptorInjectsMetadata(t *testing.T) {
	cases := []struct {
		name    string
		ctx     context.Context
		wantSet bool
		wantID  int64
	}{
		{
			name:    "已登录：注入 user_id",
			ctx:     withUserClaim(context.Background(), int64(1001)),
			wantSet: true,
			wantID:  1001,
		},
		{
			name:    "字符串 claim 同样可注入",
			ctx:     withUserClaim(context.Background(), "2002"),
			wantSet: true,
			wantID:  2002,
		},
		{
			name:    "未登录：静默透传，不注入",
			ctx:     context.Background(),
			wantSet: false,
		},
		{
			name:    "claim 非法：静默透传，不注入",
			ctx:     withUserClaim(context.Background(), "abc"),
			wantSet: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var captured context.Context
			// 拦截器不使用 cc 与调用选项，传 nil 即可
			err := OutgoingUserIdInterceptor(c.ctx, "/pb.Service/Method", nil, nil, nil,
				captureInvoker(&captured))
			if err != nil {
				t.Fatalf("interceptor should not fail: %v", err)
			}
			md, ok := metadata.FromOutgoingContext(captured)
			if !c.wantSet {
				if ok {
					if vals := md.Get(MetadataUserIDKey); len(vals) != 0 {
						t.Fatalf("user_id should not be set, got %v", vals)
					}
				}
				return
			}
			if !ok {
				t.Fatal("outgoing metadata missing")
			}
			vals := md.Get(MetadataUserIDKey)
			if len(vals) != 1 {
				t.Fatalf("user_id values = %v, want 1 value", vals)
			}
			if got, err := strconv.ParseInt(vals[0], 10, 64); err != nil || got != c.wantID {
				t.Fatalf("user_id = %q, want %d", vals[0], c.wantID)
			}
		})
	}
}

func TestOutgoingUserIdInterceptorPassesInvokerError(t *testing.T) {
	// 拦截器不能吞掉下游调用的错误
	wantErr := errors.New("rpc unavailable")
	err := OutgoingUserIdInterceptor(withUserClaim(context.Background(), int64(1)),
		"/pb.Service/Method", nil, nil, nil,
		func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return wantErr
		})
	if !errors.Is(err, wantErr) {
		t.Fatalf("err = %v, want %v", err, wantErr)
	}
}

func TestUserIdFromIncomingMD(t *testing.T) {
	cases := []struct {
		name    string
		ctx     context.Context
		wantID  int64
		wantErr bool
	}{
		{
			name:   "透传成功",
			ctx:    metadata.NewIncomingContext(context.Background(), metadata.Pairs(MetadataUserIDKey, "1001")),
			wantID: 1001,
		},
		{
			name:    "无入站 metadata",
			ctx:     context.Background(),
			wantErr: true,
		},
		{
			name:    "缺少 user_id 键",
			ctx:     metadata.NewIncomingContext(context.Background(), metadata.Pairs("other", "1")),
			wantErr: true,
		},
		{
			name:    "非数字",
			ctx:     metadata.NewIncomingContext(context.Background(), metadata.Pairs(MetadataUserIDKey, "abc")),
			wantErr: true,
		},
		{
			name:    "零值视为无效",
			ctx:     metadata.NewIncomingContext(context.Background(), metadata.Pairs(MetadataUserIDKey, "0")),
			wantErr: true,
		},
		{
			name:    "负数视为无效",
			ctx:     metadata.NewIncomingContext(context.Background(), metadata.Pairs(MetadataUserIDKey, "-5")),
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id, err := UserIdFromIncomingMD(c.ctx)
			if c.wantErr {
				if err == nil {
					t.Fatal("want error, got nil")
				}
				if got := errorx.CodeOf(err); got != errorx.CodeUnauthorized {
					t.Fatalf("code = %d, want %d", got, errorx.CodeUnauthorized)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if id != c.wantID {
				t.Fatalf("id = %d, want %d", id, c.wantID)
			}
		})
	}
}
