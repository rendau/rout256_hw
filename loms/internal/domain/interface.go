package domain

//go:generate mockery --name IEventProducer --output ./mocks --filename event_producer.go
type IEventProducer interface {
	SendMessage(key string, value []byte) error
}
