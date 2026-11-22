//go:build integration

// 集成测试：依赖本地 MySQL（127.0.0.1:3306）与 Redis（127.0.0.1:6379）。
// 运行：make test-integration 或 go test -tags integration ./internal/model/
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

func newTestPayOrderModel(t *testing.T) PayOrderModel {
	t.Helper()
	dsn := envOrDefault("TJXT_TEST_DSN_PAY",
		"root:0000@tcp(127.0.0.1:3306)/tj_pay?charset=utf8mb4&parseTime=true&loc=Local")
	conn := sqlx.NewMysql(dsn)
	var one int
	if err := conn.QueryRowCtx(context.Background(), &one, "select 1"); err != nil {
		t.Skipf("mysql unavailable, skip integration test: %v", err)
	}
	cc := cache.CacheConf{{RedisConf: redis.RedisConf{Host: envOrDefault("TJXT_TEST_REDIS", "127.0.0.1:6379"), Type: "node"}, Weight: 100}}
	return NewPayOrderModel(conn, cc)
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// TestPayOrderMarkToSuccessIdempotent 锁定支付闭环的幂等语义：
// 并发/重复回调下只有第一次状态流转成功并应发事件，第二次必须返回 false。
func TestPayOrderMarkToSuccessIdempotent(t *testing.T) {
	m := newTestPayOrderModel(t)
	ctx := context.Background()
	id := time.Now().UnixMilli()

	po := &PayOrder{
		Id:             id,
		BizOrderNo:     id,
		PayOrderNo:     id + 1,
		BizUserId:      1,
		PayChannelCode: "mock",
		Amount:         10000,
		PayType:        4,
		Status:         1, // 待支付
		PayOverTime:    time.Now().Add(time.Hour),
	}
	var insertedID int64
	t.Cleanup(func() {
		if insertedID != 0 {
			_ = m.Delete(context.Background(), insertedID)
		}
	})
	result, err := m.Insert(ctx, po)
	if err != nil {
		t.Fatalf("insert pay order: %v", err)
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get inserted pay order id: %v", err)
	}
	insertedID = insertID

	first, err := m.MarkToSuccess(ctx, insertedID, "SUCCESS", "ok")
	if err != nil || !first {
		t.Fatalf("first MarkToSuccess = (%v, %v), want (true, nil)", first, err)
	}
	second, err := m.MarkToSuccess(ctx, insertedID, "SUCCESS", "ok")
	if err != nil || second {
		t.Fatalf("second MarkToSuccess = (%v, %v), want (false, nil)", second, err)
	}

	got, err := m.FindOne(ctx, insertedID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != 3 {
		t.Fatalf("status = %d, want 3 (success)", got.Status)
	}
	if got.PaySuccessTime.Time.IsZero() {
		t.Fatal("pay_success_time should be set")
	}
}

// TestPayOrderMarkToClosedConditional 关单只允许发生在待提交/待支付状态。
func TestPayOrderMarkToClosedConditional(t *testing.T) {
	m := newTestPayOrderModel(t)
	ctx := context.Background()
	id := time.Now().UnixMilli()

	po := &PayOrder{
		Id: id, BizOrderNo: id, PayOrderNo: id + 1, BizUserId: 1,
		PayChannelCode: "mock", Amount: 5000, PayType: 4,
		Status: 1, PayOverTime: time.Now().Add(time.Hour),
	}
	var insertedID int64
	t.Cleanup(func() {
		if insertedID != 0 {
			_ = m.Delete(context.Background(), insertedID)
		}
	})
	result, err := m.Insert(ctx, po)
	if err != nil {
		t.Fatal(err)
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get inserted pay order id: %v", err)
	}
	insertedID = insertID

	// 先关单成功，再标记支付成功必须失败（终态保护）
	closed, err := m.MarkToClosed(ctx, insertedID, "CLOSE", "closed")
	if err != nil || !closed {
		t.Fatalf("MarkToClosed = (%v, %v), want (true, nil)", closed, err)
	}
	transit, err := m.MarkToSuccess(ctx, insertedID, "SUCCESS", "ok")
	if err != nil || transit {
		t.Fatalf("MarkToSuccess after close = (%v, %v), want (false, nil)", transit, err)
	}
	got, _ := m.FindOne(ctx, insertedID)
	if got.Status != 2 {
		t.Fatalf("status = %d, want 2 (closed)", got.Status)
	}
}
