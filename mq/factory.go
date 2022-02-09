package mq

import (
	"context"

	"channels-instagram-dm/domain"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type factory struct {
	ctx      context.Context
	conn     stan.Conn
	logger   domain.Logger
	producer domain.Producer
}

func Factory(ctx context.Context, logger domain.Logger, host, clusterID, clientID string) (domain.MQ, error) {
	mq, err := stan.Connect(
		clusterID,
		clientID,
		stan.NatsURL(host),
		stan.Pings(5, 288),
	)
	if err != nil {
		return nil, err
	}

	mq.NatsConn().SetDisconnectErrHandler(func(conn *nats.Conn, err error) {
		logger.Critical("MQ connection is offline", nil)
	})

	f := &factory{
		ctx:    ctx,
		conn:   mq,
		logger: logger,
	}

	go func() {
		select {
		case <-ctx.Done():
			_ = f.conn.Close()
			return
		}
	}()

	return f, nil
}

func (f *factory) Producer() domain.Producer {
	if f.producer != nil {
		return f.producer
	}

	f.producer = makeProducer(f.ctx, f.logger, f.conn)
	return f.producer
}

func (f *factory) Consumer() domain.Consumer {
	return makeConsumer(f.ctx, f.logger, f.conn)
}
