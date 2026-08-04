package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// Binding 队列绑定配置
type Binding struct {
	Queue      string
	Exchange   string
	RoutingKey string
	Kind       string // direct, topic, fanout, headers
}

// Handler 泛型处理器类型
type Handler[T any] func(ctx context.Context, msg *T) error

// registration 内部注册项
type registration struct {
	binding     Binding
	handler     any
	msgType     reflect.Type
	queue       *amqp.Queue
	consumerTag string
}

// Client MQ 客户端
type Client struct {
	conn       *amqp.Connection
	ch         *amqp.Channel
	dsn        string
	regs       []registration
	mu         sync.Mutex
	reconnect  chan struct{}
	stopCh     chan struct{}
	wg         sync.WaitGroup
	autoAck    bool
	prefetch   int
	retryDelay time.Duration
}

// NewClient 创建 MQ 客户端
func NewClient(dsn string) *Client {
	return &Client{
		dsn:        dsn,
		regs:       make([]registration, 0),
		reconnect:  make(chan struct{}, 1),
		stopCh:     make(chan struct{}),
		autoAck:    false,
		prefetch:   10,
		retryDelay: 5 * time.Second,
	}
}

// SetPrefetch 设置预取数量
func (c *Client) SetPrefetch(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prefetch = n
	if c.ch != nil {
		_ = c.ch.Qos(n, 0, false)
	}
}

// SetAutoAck 设置自动确认
func (c *Client) SetAutoAck(autoAck bool) {
	c.autoAck = autoAck
}

// Register 注册消费者
// binding: 队列绑定配置
// handler: 业务处理函数，签名为 func(ctx context.Context, msg *T) error
func Register[T any](c *Client, binding Binding, handler Handler[T]) {
	var zero T
	msgType := reflect.TypeOf(zero)
	if msgType.Kind() != reflect.Ptr {
		msgType = reflect.PtrTo(msgType)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.regs = append(c.regs, registration{
		binding: binding,
		handler: handler,
		msgType: msgType,
	})
}

// connect 建立连接
func (c *Client) connect() error {
	conn, err := amqp.Dial(c.dsn)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	if err := ch.Qos(c.prefetch, 0, false); err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	c.conn = conn
	c.ch = ch
	return nil
}

// setupBindings 声明交换机、队列、绑定
func (c *Client) setupBindings() error {
	for i := range c.regs {
		reg := &c.regs[i]
		b := reg.binding

		if b.Kind == "" {
			b.Kind = "direct"
		}

		// 声明交换机
		if b.Exchange != "" {
			if err := c.ch.ExchangeDeclare(
				b.Exchange,
				b.Kind,
				true,  // durable
				false, // autoDelete
				false, // internal
				false, // noWait
				nil,
			); err != nil {
				return fmt.Errorf("declare exchange %s: %w", b.Exchange, err)
			}
		}

		// 声明队列
		q, err := c.ch.QueueDeclare(
			b.Queue,
			true,  // durable
			false, // autoDelete
			false, // exclusive
			false, // noWait
			nil,
		)
		if err != nil {
			return fmt.Errorf("declare queue %s: %w", b.Queue, err)
		}
		reg.queue = &q

		// 绑定队列到交换机
		if b.Exchange != "" && b.RoutingKey != "" {
			if err := c.ch.QueueBind(
				b.Queue,
				b.RoutingKey,
				b.Exchange,
				false, // noWait
				nil,
			); err != nil {
				return fmt.Errorf("bind queue %s to exchange %s with key %s: %w", b.Queue, b.Exchange, b.RoutingKey, err)
			}
		}
	}
	return nil
}

// Start 启动消费
func (c *Client) Start(ctx context.Context) error {
	if err := c.connect(); err != nil {
		logx.Errorf("connect to RabbitMQ failed: %v", err)
		return err
	}
	defer c.conn.Close()
	defer c.ch.Close()

	if err := c.setupBindings(); err != nil {
		logx.Errorf("setup bindings failed: %v", err)
		return err
	}

	// 为每个注册项启动消费协程
	for i := range c.regs {
		reg := &c.regs[i]
		c.wg.Add(1)
		go c.consumeLoop(ctx, reg)
	}

	logx.Infof("MQ client started, consuming %d queues", len(c.regs))

	// 等待停止信号
	<-c.stopCh
	c.wg.Wait()
	return nil
}

// Stop 停止消费
func (c *Client) Stop() {
	close(c.stopCh)
}

// consumeLoop 消费循环
func (c *Client) consumeLoop(ctx context.Context, reg *registration) {
	defer c.wg.Done()

	logx.Infof("consume loop started for queue %s", reg.binding.Queue)

	b := reg.binding
	consumerTag := fmt.Sprintf("consumer-%s", b.Queue)
	reg.consumerTag = consumerTag

	deliveries, err := c.ch.Consume(
		b.Queue,
		consumerTag,
		c.autoAck, // autoAck
		false,     // exclusive
		false,     // noLocal
		false,     // noWait
		nil,
	)
	if err != nil {
		logx.Errorf("consume queue %s failed: %v", b.Queue, err)
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
				logx.Infof("delivery channel closed for queue %s", b.Queue)
				return
			}
			c.handleDelivery(ctx, reg, d)
		}
	}
}

// handleDelivery 处理单条消息
func (c *Client) handleDelivery(ctx context.Context, reg *registration, d amqp.Delivery) {
	handleCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 创建消息实例
	msgPtr := reflect.New(reg.msgType.Elem())
	msgInterface := msgPtr.Interface()

	// 反序列化
	if err := json.Unmarshal(d.Body, msgInterface); err != nil {
		logx.Errorf("unmarshal message failed, queue=%s, err=%v, body=%s", reg.binding.Queue, err, string(d.Body))
		if !c.autoAck {
			_ = d.Nack(false, false) // 拒绝不重新入队，避免死循环
		}
		return
	}

	// 调用 handler
	handler := reg.handler
	handlerValue := reflect.ValueOf(handler)

	// 验证 handler 签名: func(context.Context, *T) error
	if handlerValue.Kind() != reflect.Func {
		logx.Errorf("handler is not a function for queue %s", reg.binding.Queue)
		if !c.autoAck {
			_ = d.Nack(false, false)
		}
		return
	}

	// 调用 handler
	results := handlerValue.Call([]reflect.Value{
		reflect.ValueOf(handleCtx),
		msgPtr,
	})

	// 检查错误
	if len(results) == 1 && !results[0].IsNil() {
		err := results[0].Interface().(error)
		logx.Errorf("handle message failed, queue=%s, err=%v", reg.binding.Queue, err)
		if !c.autoAck {
			_ = d.Nack(false, true) // 重新入队重试
		}
		return
	}

	if !c.autoAck {
		_ = d.Ack(false)
	}
}
