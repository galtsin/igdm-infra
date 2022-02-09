package domain

type ConsumerHandler func([]byte) error

type MQ interface {
	Producer() Producer
	Consumer() Consumer
}

type Producer interface {
	Publish(subject string, data []byte) error
}

type Consumer interface {
	Subscribe(subject string, consumer ConsumerHandler) error
	IsActive() bool
	Close() error
}
