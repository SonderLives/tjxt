//go:build integration

// 集成测试：依赖本地 MySQL 与 Redis。运行：go test -tags integration ./internal/model/
package model

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func newTestOrderModel(t *testing.T) OrderModel {
	t.Helper()
	dsn := envOrDefaultT("TJXT_TEST_DSN_TRADE",
		"root:0000@tcp(127.0.0.1:3306)/tj_trade?charset=utf8mb4&parseTime=true&loc=Local")
	conn := sqlx.NewMysql(dsn)
	var one int
	if err := conn.QueryRowCtx(context.Background(), &one, "select 1"); err != nil {
		t.Skipf("mysql unavailable, skip integration test: %v", err)
	}
	cc := cache.CacheConf{{RedisConf: redis.RedisConf{Host: envOrDefaultT("TJXT_TEST_REDIS", "127.0.0.1:6379"), Type: "node"}, Weight: 100}}
	return NewOrderModel(conn, cc)
}

func envOrDefaultT(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// TestOrderMarkPaidIdempotent 锁定订单回写的幂等语义：
// pay.success 事件重复投递时，只有第一次把待支付订单置为已支付。
func TestOrderMarkPaidIdempotent(t *testing.T) {
	m := newTestOrderModel(t)
	ctx := context.Background()
	const id = 991101

	o := &Order{
		Id: id, UserId: 1, Status: 1, Message: "it",
		TotalAmount: 10000, RealAmount: 10000, Creater: 1, Updater: 1,
	}
	t.Cleanup(func() { _ = m.Delete(context.Background(), id) })
	if _, err := m.Insert(ctx, o); err != nil {
		t.Fatalf("insert order: %v", err)
	}

	first, err := m.MarkPaid(ctx, id, 991102, "mock", time.Now(), 10000)
	if err != nil || !first {
		t.Fatalf("first MarkPaid = (%v, %v), want (true, nil)", first, err)
	}
	second, err := m.MarkPaid(ctx, id, 991102, "mock", time.Now(), 10000)
	if err != nil || second {
		t.Fatalf("second MarkPaid = (%v, %v), want (false, nil)", second, err)
	}

	got, err := m.FindOne(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != 2 {
		t.Fatalf("status = %d, want 2 (paid)", got.Status)
	}
	if !got.PayOrderNo.Valid || got.PayOrderNo.Int64 != 991102 {
		t.Fatalf("pay_order_no = %+v, want 991102", got.PayOrderNo)
	}
	if !got.PayTime.Valid {
		t.Fatal("pay_time should be set")
	}
}

// TestOrderMarkPaidSkipsNonPending 非待支付订单不允许被标记已支付（终态保护）。
func TestOrderMarkPaidSkipsNonPending(t *testing.T) {
	m := newTestOrderModel(t)
	ctx := context.Background()
	const id = 991103

	o := &Order{
		Id: id, UserId: 1, Status: 3, Message: "it-closed",
		TotalAmount: 10000, RealAmount: 10000, Creater: 1, Updater: 1,
	}
	t.Cleanup(func() { _ = m.Delete(context.Background(), id) })
	if _, err := m.Insert(ctx, o); err != nil {
		t.Fatal(err)
	}

	transit, err := m.MarkPaid(ctx, id, 991104, "mock", time.Now(), 10000)
	if err != nil || transit {
		t.Fatalf("MarkPaid on closed order = (%v, %v), want (false, nil)", transit, err)
	}
	got, _ := m.FindOne(ctx, id)
	if got.Status != 3 {
		t.Fatalf("status = %d, want 3 (closed, unchanged)", got.Status)
	}
}
