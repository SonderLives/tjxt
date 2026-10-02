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
// 内部维护一条连接与 channel：发布全程持锁串行化（避免并发发布时
// confirm 通知串台），连接断开后首次发布自动重连并重试一次；
// 采用 confirm + mandatory 模式——消息被确认或被退回（无队列可路由）后才返回。
type Producer struct {
	dsn  string
	conn *amqp091.Connection
	ch   *amqp091.Channel

	mu sync.Mutex
}

// NewProducer 建立 RabbitMQ 连接。
func NewProducer(dsn string) (*Producer, error) {
	p := &Producer{dsn: dsn}
	if err := p.dial(); err != nil {
		return nil, err
	}
	return p, nil
}

// dial 建立连接并开启 confirm 模式的 channel。
func (p *Producer) dial() error {
	conn, err := amqp091.Dial(p.dsn)
	if err != nil {
		return fmt.Errorf("dial rabbitmq: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("open channel: %w", err)
	}
	if err := ch.Confirm(false); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return fmt.Errorf("enable confirm: %w", err)
	}
	p.conn, p.ch = conn, ch
	return nil
}

// ensureChannel 连接不可用时重连一次（供发布失败后重试）。
func (p *Producer) ensureChannel() error {
	if p.conn != nil && !p.conn.IsClosed() && p.ch != nil && !p.ch.IsClosed() {
		return nil
	}
	logx.Info("rabbitmq connection lost, reconnecting producer...")
	_ = p.ch.Close()
	_ = p.conn.Close()
	return p.dial()
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
	if err := p.publishRaw(ctx, exchange, routingKey, data); err != nil {
		// 断线重连后重试一次
		if err2 := p.ensureChannel(); err2 != nil {
			return err
		}
		if err2 := p.publishRaw(ctx, exchange, routingKey, data); err2 != nil {
			return err2
		}
	}
	return nil
}

// PublishJSON 发布一条原始 JSON 字节消息。
func (p *Producer) PublishJSON(ctx context.Context, exchange, routingKey string, data []byte) error {
	return p.Publish(ctx, exchange, routingKey, json.RawMessage(data))
}

// publishRaw 发布并等待 broker 确认。
// 全程持锁：同一 channel 上的 confirm 与 return 通知按发布顺序到达，
// 串行化后才能一一对应，避免并发发布时通知串台。
func (p *Producer) publishRaw(ctx context.Context, exchange, routingKey string, data []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if exchange != "" {
		if err := p.declareExchange(exchange); err != nil {
			return fmt.Errorf("declare exchange %s: %w", exchange, err)
		}
	}

	confirms := p.ch.NotifyPublish(make(chan amqp091.Confirmation, 1))
	returns := p.ch.NotifyReturn(make(chan amqp091.Return, 1))
	err := p.ch.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		true, // mandatory：无队列可路由时消息退回而不是被静默丢弃
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Body:         data,
		},
	)
	if err != nil {
		return fmt.Errorf("publish message: %w", err)
	}

	// 同时等待 confirm 与 return：RabbitMQ 对无法路由的 mandatory 消息
	// 只发 basic.return、不发 confirm，二者必须一起监听
	select {
	case ret := <-returns:
		return fmt.Errorf("message unroutable, exchange=%s routingKey=%s replyText=%s",
			ret.Exchange, ret.RoutingKey, ret.ReplyText)
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
	if p.ch != nil {
		if err := p.ch.Close(); err != nil {
			logx.Errorf("close rabbitmq channel: %v", err)
		}
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			logx.Errorf("close rabbitmq conn: %v", err)
		}
	}
}
