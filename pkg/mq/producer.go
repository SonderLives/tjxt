package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// Producer 轻量 RabbitMQ 发布者。
//
// 内部维护一条连接与 channel，并发安全；
// 采用 confirm 模式，Publish 返回前确保消息已送达交换机。
type Producer struct {
	dsn  string
	conn *amqp091.Connection
	ch   *amqp091.Channel

	mu   sync.Mutex
	done bool
}

// NewProducer 建立 RabbitMQ 连接。
func NewProducer(dsn string) (*Producer, error) {
	conn, err := amqp091.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open channel: %w", err)
	}
	if err := ch.Confirm(false); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("enable confirm: %w", err)
	}
	return &Producer{dsn: dsn, conn: conn, ch: ch}, nil
}

// declareExchange 确保交换机存在（幂等声明）。
func (p *Producer) declareExchange(exchange string) error {
	return p.ch.ExchangeDeclare(
		exchange,
		"direct", // 与下游消费者声明的交换机类型保持一致
		true,     // durable
		false,    // autoDelete
		false,    // internal
		false,    // noWait
		nil,
	)
}

// Publish 发布一条 JSON 消息。
func (p *Producer) Publish(ctx context.Context, exchange, routingKey string, body any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}
	return p.publishRaw(ctx, exchange, routingKey, data)
}

// PublishJSON 发布一条原始 JSON 字节消息。
func (p *Producer) PublishJSON(ctx context.Context, exchange, routingKey string, data []byte) error {
	return p.publishRaw(ctx, exchange, routingKey, data)
}

func (p *Producer) publishRaw(ctx context.Context, exchange, routingKey string, data []byte) error {
	if exchange != "" {
		if err := p.declareExchange(exchange); err != nil {
			return fmt.Errorf("declare exchange %s: %w", exchange, err)
		}
	}

	confirms := p.ch.NotifyPublish(make(chan amqp091.Confirmation, 1))
	err := p.ch.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Body:         data,
		},
	)
	if err != nil {
		return fmt.Errorf("publish message: %w", err)
	}
	select {
	case c := <-confirms:
		if !c.Ack {
			return fmt.Errorf("broker rejected message")
		}
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("confirm message timeout")
	}
}

// Close 关闭连接，幂等。
func (p *Producer) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.done {
		return
	}
	p.done = true
	if err := p.ch.Close(); err != nil {
		logx.Errorf("close rabbitmq channel: %v", err)
	}
	if err := p.conn.Close(); err != nil {
		logx.Errorf("close rabbitmq conn: %v", err)
	}
}
