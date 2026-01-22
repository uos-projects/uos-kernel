package mq

import (
	"context"
	"time"
)

// QualityCode 数据质量码
type QualityCode int

const (
	QualityGood QualityCode = iota
	QualityBad
	QualityUncertain
	QualitySubstituted
)

// MQMessage MQ 消息
type MQMessage struct {
	Value     interface{}
	Timestamp time.Time
	Quality   QualityCode
	Source    string
}

// MQConsumer MQ 消费者接口
type MQConsumer interface {
	// Consume 消费一条消息（阻塞）
	Consume(ctx context.Context, topic string) (*MQMessage, error)

	// Subscribe 订阅主题，通过 channel 返回消息
	Subscribe(ctx context.Context, topic string) (<-chan *MQMessage, error)

	// Close 关闭消费者
	Close() error
}
