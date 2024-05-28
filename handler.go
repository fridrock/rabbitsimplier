package rabbitsimplier

import amqp "github.com/rabbitmq/amqp091-go"

type Handler interface {
	Handle(msg amqp.Delivery)
}

type HandlerFunc func(msg amqp.Delivery)

func (hf HandlerFunc) Handle(msg amqp.Delivery) {
	hf(msg)
}

func NewHandlerFunc(f func(msg amqp.Delivery)) HandlerFunc {
	return f
}
