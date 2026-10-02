package xerr

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// toStatusError 构造 zrpc 客户端实际收到的错误形态：*status.Error(Unknown, 远端Error())
func toStatusError(msg string) error {
	return status.Error(codes.Unknown, msg)
}

func TestFromGRPCRestoresBizError(t *testing.T) {
	// 模拟 zrpc 传输后的错误：远端 *Error 以 status error 的 Message 传回
	remote := New(CodeNotFound, "课程不存在")
	wire := errors.New(remote.Error()) // 客户端收到的是普通 error

	// status.FromError 对非 status 错误返回 ok=false，说明真实场景里
	// zrpc 客户端拿到的是 *status.Error；这里直接用 status.Error 构造。
	got, ok := FromGRPC(toStatusError(remote.Error()))
	if !ok {
		t.Fatal("should restore from status error")
	}
	if got.Code != CodeNotFound || got.Msg != "课程不存在" {
		t.Fatalf("restored = %+v", got)
	}
	_ = wire
}

func TestFromGRPCStripsCause(t *testing.T) {
	remote := Wrap(errors.New("db down"), CodeInternal, "查询订单失败")
	got, ok := FromGRPC(toStatusError(remote.Error()))
	if !ok {
		t.Fatal("should restore")
	}
	if got.Msg != "查询订单失败" {
		t.Fatalf("cause leaked into msg: %q", got.Msg)
	}
}

func TestFromGRPCRejectsPlainErrors(t *testing.T) {
	if _, ok := FromGRPC(nil); ok {
		t.Fatal("nil should not restore")
	}
	if _, ok := FromGRPC(errors.New("some random error")); ok {
		t.Fatal("random error without code= prefix should not restore")
	}
}
