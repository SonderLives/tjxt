//go:build integration

// 集成测试：依赖本地 RabbitMQ（127.0.0.1:5672）。
// 运行：go test -tags integration ./mq/
package mq

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type itEvent struct {
	OrderID int64 `json:"orderId"`
}

const (
	itExchange = "it.tjxt.exchange"
	itQueue    = "it.tjxt.queue"
	itKey      = "it.tjxt"
)

// TestProduceConsumeAndDeadLetter 锁定消息链路三件事：
// 1) confirm 发布 + 消费成功 ack；
// 2) 业务处理失败 → 先重回队列重试，再失败进死信队列；
// 3) 反序列化失败 → 直接进死信队列。
func TestProduceConsumeAndDeadLetter(t *testing.T) {
	dsn := os.Getenv("TJXT_TEST_AMQP")
	if dsn == "" {
		dsn = "amqp://rabbitmq:rabbitmq@127.0.0.1:5672/"
	}

	// 前置：清理可能残留的测试队列
	if conn, err := amqp091.Dial(dsn); err == nil {
		ch, _ := conn.Channel()
		_, _ = ch.QueueDelete(itQueue, false, false, false)
		_, _ = ch.QueueDelete(itQueue+".dead", false, false, false)
		_ = ch.ExchangeDelete(itExchange, false, false)
		_ = conn.Close()
	} else {
		t.Skipf("rabbitmq unavailable, skip integration test: %v", err)
	}

	p, err := NewProducer(dsn)
	if err != nil {
		t.Fatalf("new producer: %v", err)
	}
	defer p.Close()

	client := NewClient(dsn)
	var mu sync.Mutex
	handled := make([]int64, 0, 2)
	delivered := make(chan struct{}, 4)

	Register(client, Binding{Queue: itQueue, Exchange: itExchange, RoutingKey: itKey},
		func(ctx context.Context, evt *itEvent) error {
			mu.Lock()
			handled = append(handled, evt.OrderID)
			mu.Unlock()
			delivered <- struct{}{}
			// OrderID=1：处理成功；OrderID=2：两次处理都失败（先重试后进死信）
			if evt.OrderID == 2 {
				return errAlwaysFail{}
			}
			return nil
		})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_ = client.Start(ctx)
	}()
	t.Cleanup(client.Stop)

	// 等消费者就绪（队列被声明）
	waitForQueue(t, dsn, itQueue)

	// 1) 正常消息：消费成功
	if err := p.Publish(context.Background(), itExchange, itKey, itEvent{OrderID: 1}); err != nil {
		t.Fatalf("publish valid: %v", err)
	}
	select {
	case <-delivered:
	case <-time.After(10 * time.Second):
		t.Fatal("valid message not consumed in time")
	}

	// 2) 业务处理失败：先重回队列重试一次，再失败进死信
	if err := p.Publish(context.Background(), itExchange, itKey, itEvent{OrderID: 2}); err != nil {
		t.Fatalf("publish failing: %v", err)
	}

	// 3) 反序列化失败的消息 → 直接进死信
	if err := p.Publish(context.Background(), itExchange, itKey, map[string]any{"orderId": "not-a-number"}); err != nil {
		t.Fatalf("publish poison: %v", err)
	}
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if n := deadLetterCount(t, dsn, itQueue+".dead"); n >= 2 {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	// 死信队列应含 2 条：业务失败重试超限 1 条 + 反序列化失败 1 条
	if n := deadLetterCount(t, dsn, itQueue+".dead"); n < 2 {
		t.Fatalf("dead letter queue count = %d, want >= 2", n)
	}
	mu.Lock()
	defer mu.Unlock()
	// OrderID=1 恰好被成功消费一次（OrderID=2 的两次投递都失败）
	successCount := 0
	for _, id := range handled {
		if id == 1 {
			successCount++
		}
	}
	if successCount != 1 {
		t.Fatalf("valid message handled %d times, want 1", successCount)
	}
}

type errAlwaysFail struct{}

func (errAlwaysFail) Error() string { return "integration: always fail" }

func waitForQueue(t *testing.T, dsn, queue string) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := amqp091.Dial(dsn)
		if err != nil {
			t.Fatal(err)
		}
		ch, _ := conn.Channel()
		_, err = ch.QueueDeclarePassive(queue, true, false, false, false, nil)
		_ = ch.Close()
		_ = conn.Close()
		if err == nil {
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
	t.Fatalf("queue %s not declared in time", queue)
}

func deadLetterCount(t *testing.T, dsn, queue string) int {
	t.Helper()
	conn, err := amqp091.Dial(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()
	q, err := ch.QueueDeclarePassive(queue, true, false, false, false, nil)
	if err != nil {
		return 0
	}
	return int(q.Messages)
}
