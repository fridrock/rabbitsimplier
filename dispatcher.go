package rabbitsimplier

import (
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Dispatcher interface {
	Dispatch(<-chan amqp.Delivery)
}

type DispatcherFunc func(<-chan amqp.Delivery)

func (df DispatcherFunc) Dispatch(ch <-chan amqp.Delivery) {
	df(ch)
}
func NewDispactherFunc(f func(ch <-chan amqp.Delivery)) DispatcherFunc {
	return f
}

type RDispatcher struct {
	handlers map[string]Handler
}

func NewRDispacher() RDispatcher {
	return RDispatcher{
		handlers: make(map[string]Handler),
	}
}
func (rd *RDispatcher) RegisterHandler(routingKey string, h Handler) {
	rd.handlers[routingKey] = h
}
func (rd RDispatcher) Dispatch(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		v, ok := rd.handlers[d.RoutingKey]
		if !ok {
			slog.Error(fmt.Sprintf("unknown routing key: %s", d.RoutingKey))
		} else {
			//TODO make some logic for error handling
			go v.Handle(d)
		}
	}
}
