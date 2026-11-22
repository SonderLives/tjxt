package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// DeadLetterExchange 死信交换机：每个业务队列配套一个 <queue>.dead 死信队列，
// 消费失败重试超限 / 消息反序列化失败的消息路由到这里，不再无限重回或静默丢弃。
const DeadLetterExchange = "tjxt.dlx"

const (
	initialReconnectDelay = 2 * time.Second
	maxReconnectDelay     = 30 * time.Second
)

// Binding 队列绑定配置
type Binding struct {
	Queue      string
	Exchange   string
	RoutingKey string
	Kind       string // direct, topic, fanout, headers
}

// Handler 泛型处理器类型：func(ctx context.Context, msg *T) error
type Handler[T any] func(ctx context.Context, msg *T) error

// registration 内部注册项
type registration struct {
	binding Binding
	handler any
	msgType reflect.Type
}

// Client MQ 消费者客户端。
//
// Start 内部带自动重连：连接断开后按指数退避重连并重新声明队列；
// 每个队列消费失败先重回队列重试一次，再次失败进入死信队列。
type Client struct {
	dsn      string
	regs     []registration
	mu       sync.Mutex
	stopCh   chan struct{}
	stopOnce sync.Once
	autoAck  bool
	prefetch int
}

// NewClient 创建消费者客户端。
func NewClient(dsn string) *Client {
	return &Client{
		dsn:      dsn,
		regs:     make([]registration, 0),
		stopCh:   make(chan struct{}),
		autoAck:  false,
		prefetch: 10,
	}
}

// SetPrefetch 设置预取数量（连接建立时生效）。
func (c *Client) SetPrefetch(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prefetch = n
}

// SetAutoAck 设置是否自动确认。
func (c *Client) SetAutoAck(autoAck bool) {
	c.autoAck = autoAck
}

// Register 注册消费者。handler 签名必须为 func(ctx context.Context, msg *T) error。
func Register[T any](c *Client, binding Binding, handler Handler[T]) {
	msgType := reflect.TypeFor[T]()
	if msgType.Kind() != reflect.Pointer {
		msgType = reflect.PointerTo(msgType)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.regs = append(c.regs, registration{
		binding: binding,
		handler: handler,
		msgType: msgType,
	})
}

// Stop 停止消费，幂等。
func (c *Client) Stop() {
	c.stopOnce.Do(func() { close(c.stopCh) })
}

// Start 启动消费（阻塞直到 Stop / ctx 取消）。
// 连接断开后自动重连：指数退避，上限 30s。
func (c *Client) Start(ctx context.Context) error {
	backoff := initialReconnectDelay
	for {
		err := c.run(ctx)
		if err == nil {
			return nil // Stop 或 ctx 取消
		}
		logx.Errorf("mq consumer disconnected: %v, reconnecting in %v", err, backoff)
		select {
		case <-c.stopCh:
			return nil
		case <-ctx.Done():
			return nil
		case <-time.After(backoff):
		}
		if backoff < maxReconnectDelay {
			backoff *= 2
			if backoff > maxReconnectDelay {
				backoff = maxReconnectDelay
			}
		}
	}
}

// run 单轮消费：建连、声明拓扑、启动各队列消费循环；
// 任一队列的投递通道断开即返回错误（由 Start 触发整体重连）。
func (c *Client) run(ctx context.Context) error {
	conn, err := amqp091.Dial(c.dsn)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel: %w", err)
	}
	if err := ch.Qos(c.prefetch, 0, false); err != nil {
		_ = ch.Close()
		return fmt.Errorf("set qos: %w", err)
	}

	if err := c.setupBindings(ch); err != nil {
		_ = ch.Close()
		return err
	}

	errCh := make(chan error, len(c.regs))
	var wg sync.WaitGroup
	for i := range c.regs {
		wg.Add(1)
		go c.consumeLoop(ctx, ch, &c.regs[i], errCh, &wg)
	}
	logx.Infof("MQ client started, consuming %d queues", len(c.regs))

	select {
	case err := <-errCh:
		// 关闭 channel 触发所有 deliveries 通道关闭，消费循环退出后整体重连
		_ = ch.Close()
		wg.Wait()
		return err
	case <-ctx.Done():
		return nil
	case <-c.stopCh:
		return nil
	}
}

// setupBindings 声明交换机、业务队列（带死信配置）、死信队列与绑定。
func (c *Client) setupBindings(ch *amqp091.Channel) error {
	for i := range c.regs {
		reg := &c.regs[i]
		b := reg.binding

		if b.Kind == "" {
			b.Kind = "direct"
		}
		if b.Exchange != "" {
			if err := ch.ExchangeDeclare(b.Exchange, b.Kind, true, false, false, false, nil); err != nil {
				return fmt.Errorf("declare exchange %s: %w", b.Exchange, err)
			}
		}

		// 死信拓扑：DLX + <queue>.dead 队列（路由键为原队列名）
		if err := ch.ExchangeDeclare(DeadLetterExchange, "direct", true, false, false, false, nil); err != nil {
			return fmt.Errorf("declare dead letter exchange: %w", err)
		}
		deadQueue := b.Queue + ".dead"
		if _, err := ch.QueueDeclare(deadQueue, true, false, false, false, nil); err != nil {
			return fmt.Errorf("declare dead letter queue %s: %w", deadQueue, err)
		}
		if err := ch.QueueBind(deadQueue, b.Queue, DeadLetterExchange, false, nil); err != nil {
			return fmt.Errorf("bind dead letter queue %s: %w", deadQueue, err)
		}

		// 业务队列：消费失败重试超限后经 DLX 进入死信队列
		q, err := ch.QueueDeclare(b.Queue, true, false, false, false, amqp091.Table{
			"x-dead-letter-exchange":    DeadLetterExchange,
			"x-dead-letter-routing-key": b.Queue,
		})
		if err != nil {
			return fmt.Errorf("declare queue %s: %w", b.Queue, err)
		}
		if b.Exchange != "" && b.RoutingKey != "" {
			if err := ch.QueueBind(b.Queue, b.RoutingKey, b.Exchange, false, nil); err != nil {
				return fmt.Errorf("bind queue %s: %w", b.Queue, err)
			}
		}
		logx.Infof("mq queue ready, queue=%s deadLetterQueue=%s", q.Name, deadQueue)
	}
	return nil
}

// consumeLoop 单队列消费循环；投递通道关闭（连接断开）时上报错误。
func (c *Client) consumeLoop(ctx context.Context, ch *amqp091.Channel, reg *registration, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	b := reg.binding
	consumerTag := fmt.Sprintf("consumer-%s", b.Queue)

	deliveries, err := ch.Consume(b.Queue, consumerTag, c.autoAck, false, false, false, nil)
	if err != nil {
		select {
		case errCh <- fmt.Errorf("consume queue %s: %w", b.Queue, err):
		default:
		}
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.stopCh:
			return
		case d, ok := <-deliveries:
			if !ok {
				select {
				case errCh <- fmt.Errorf("delivery channel closed for queue %s", b.Queue):
				default:
				}
				return
			}
			c.handleDelivery(ctx, reg, d)
		}
	}
}

// handleDelivery 处理单条消息：
// 反序列化失败 → 直接进死信队列；业务失败 → 首次重回队列重试，再次失败进死信队列。
func (c *Client) handleDelivery(ctx context.Context, reg *registration, d amqp091.Delivery) {
	handleCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	msgPtr := reflect.New(reg.msgType.Elem())
	if err := json.Unmarshal(d.Body, msgPtr.Interface()); err != nil {
		logx.Errorf("unmarshal message failed, dead-letter it, queue=%s, err=%v, body=%s",
			reg.binding.Queue, err, string(d.Body))
		if !c.autoAck {
			_ = d.Nack(false, false)
		}
		return
	}

	results := reflect.ValueOf(reg.handler).Call([]reflect.Value{
		reflect.ValueOf(handleCtx),
		msgPtr,
	})
	if len(results) == 1 && !results[0].IsNil() {
		err := results[0].Interface().(error)
		if !d.Redelivered {
			logx.Errorf("handle message failed, retry once via requeue, queue=%s, err=%v", reg.binding.Queue, err)
			if !c.autoAck {
				_ = d.Nack(false, true)
			}
			return
		}
		logx.Errorf("handle message failed again, dead-letter it, queue=%s, err=%v", reg.binding.Queue, err)
		if !c.autoAck {
			_ = d.Nack(false, false)
		}
		return
	}
	if !c.autoAck {
		_ = d.Ack(false)
	}
}
