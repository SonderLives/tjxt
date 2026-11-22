package xerr

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestIsNotFound(t *testing.T) {
	// model.ErrNotFound 的典型形态：sqlx 的裸错误（非 *Error）
	bareNotFound := errors.New("sql: no rows in result set")

	cases := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil 不算 NotFound",
			err:  nil,
			want: false,
		},
		{
			name: "本地业务 NotFound",
			err:  NotFound("课程不存在"),
			want: true,
		},
		{
			name: "fmt.Errorf 包装的业务 NotFound 可穿透",
			err:  fmt.Errorf("查询课程: %w", NotFound("课程不存在")),
			want: true,
		},
		{
			name: "Wrap 携带底层原因的业务 NotFound",
			err:  Wrap(sql.ErrNoRows, CodeNotFound, "订单不存在"),
			want: true,
		},
		{
			name: "跨 gRPC 还原：status error 中 code=404 msg=... 文本",
			err:  toStatusError("code=404 msg=课程不存在"),
			want: true,
		},
		{
			name: "跨 gRPC 还原：带 cause 的 code=404 文本",
			err:  toStatusError("code=404 msg=课程不存在 cause=sql: no rows"),
			want: true,
		},
		{
			name: "其他业务错误码不算 NotFound",
			err:  Conflict("重复报名"),
			want: false,
		},
		{
			name: "业务错误链上第一个 *Error 决定结果（与 CodeOf 一致）",
			err:  fmt.Errorf("wrapper: %w", Conflict("重复报名")),
			want: false,
		},
		{
			name: "普通错误不算",
			err:  errors.New("boom"),
			want: false,
		},
		{
			name: "sql.ErrNoRows 裸错误不算（应继续用 errors.Is 判断）",
			err:  sql.ErrNoRows,
			want: false,
		},
		{
			name: "sqlx model.ErrNotFound 同形态的裸错误不算",
			err:  bareNotFound,
			want: false,
		},
		{
			name: "非 404 的 gRPC status 文本不算",
			err:  toStatusError("code=409 msg=操作冲突"),
			want: false,
		},
		{
			name: "gRPC 内部错误码不算",
			err:  status.Error(codes.Internal, "boom"),
			want: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := IsNotFound(c.err); got != c.want {
				t.Fatalf("IsNotFound(%v) = %v, want %v", c.err, got, c.want)
			}
		})
	}
}

// TestIsNotFoundWrappedStatusError 验证已知限制：FromGRPC 只在顶层判断，
// 被本地 fmt.Errorf 包装过的 status error 无法还原（与既有 FromGRPC 行为一致）。
func TestIsNotFoundWrappedStatusError(t *testing.T) {
	wrapped := fmt.Errorf("调用下游: %w", toStatusError("code=404 msg=课程不存在"))
	if IsNotFound(wrapped) {
		t.Fatal("locally wrapped status error should not be treated as NotFound")
	}
}
