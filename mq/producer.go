package mq

import (
	"context"

	"channels-instagram-dm/domain"
	"github.com/nats-io/stan.go"
)

type producer struct {
	ctx    context.Context
	conn   stan.Conn
	logger domain.Logger
}

func makeProducer(ctx context.Context, logger domain.Logger, conn stan.Conn) *producer {
	return &producer{
		ctx:    ctx,
		conn:   conn,
		logger: logger,
	}
}

func (p *producer) Publish(subject string, payload []byte) error {
	return p.conn.Publish(subject, payload)
}
