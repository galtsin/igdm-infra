package mq

import (
	"context"
	"fmt"

	"channels-instagram-dm/domain"

	"github.com/nats-io/stan.go"
)

type consumer struct {
	ctx    context.Context
	conn   stan.Conn
	logger domain.Logger
	subs   stan.Subscription
}

func makeConsumer(ctx context.Context, logger domain.Logger, conn stan.Conn) *consumer {
	return &consumer{
		ctx:    ctx,
		conn:   conn,
		logger: logger,
		subs:   nil,
	}
}

func (c *consumer) Subscribe(subject string, handler domain.ConsumerHandler) error {
	// У одного consumer может быть только одна подписка
	if err := c.Close(); err != nil {
		return err
	}

	return c.listen(subject, handler)
}

func (c *consumer) IsActive() bool {
	return c.conn.NatsConn().IsConnected()
}

func (c *consumer) Close() error {
	if c.subs != nil {
		return c.subs.Close()
	}

	return nil
}

func (c *consumer) listen(subject string, handler domain.ConsumerHandler) error {
	subs, err := c.conn.Subscribe(subject, func(msg *stan.Msg) {
		select {
		case <-c.ctx.Done():
			return
		default:
			if err := handler(msg.Data); err != nil {
				c.logger.Error(fmt.Sprintf("Failed to handle message: %s", err), nil)
				return
			}

			if err := msg.Ack(); err != nil {
				c.logger.Error(fmt.Sprintf("Failed to ack message: %s", err), nil)
			}
		}
	},
		stan.DurableName("outbound"),
		stan.DeliverAllAvailable(),
		stan.SetManualAckMode(),
		stan.MaxInflight(1),
	)
	if err != nil {
		return err
	}

	c.subs = subs

	return nil
}
