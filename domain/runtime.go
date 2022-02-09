package domain

import "context"

type RuntimeContext interface {
	Syncer() Syncer                                 // Синхронизировать глобальные горутины
	WithContext(ctx context.Context) RuntimeContext // Изменять контекст
	WithLogger(logger Logger) RuntimeContext
	Context() context.Context
	EventBus() EventBus
	Repository() Repository
	Service() Service
	Logger() Logger
	MQ() MQ
}

type Syncer interface {
	Add()
	Remove()
	Done() <-chan struct{}
}
