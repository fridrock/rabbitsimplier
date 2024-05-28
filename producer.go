package rabbitsimplier

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer interface {
	CreateChannel(*amqp.Connection) error
	CreateExchange(exchangeName, exchangeType string, exchangeConfig ...ExchangeConfig) error
	PublishMessage(ctx context.Context, exchangeName, routingKey, body string, messageConfig ...MessageConfig) error
	Stop()
}

type RProducer struct {
	Ch *amqp.Channel
}

func (bpi *RProducer) CreateChannel(conn *amqp.Connection) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	bpi.Ch = channel
	return nil
}

type ExchangeConfig struct {
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

func DefaultExchangeConfig() ExchangeConfig {
	return ExchangeConfig{
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
}
func (bpi *RProducer) CreateExchange(exchangeName, exchangeType string, exchangeConfig ...ExchangeConfig) error {
	cfg := extractExchangeConfig(exchangeConfig...)
	return bpi.Ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		cfg.Durable,
		cfg.AutoDelete,
		cfg.Internal,
		cfg.NoWait,
		cfg.Args,
	)
}

func extractExchangeConfig(exchangeConfig ...ExchangeConfig) ExchangeConfig {
	if len(exchangeConfig) == 0 {
		return DefaultExchangeConfig()
	} else {
		return exchangeConfig[0]
	}
}

type MessageConfig struct {
	Mandatory bool
	Immediate bool
}

func DefaultMessageConfig() MessageConfig {
	return MessageConfig{
		Mandatory: false,
		Immediate: false,
	}
}
func (bpi *RProducer) PublishMessage(ctx context.Context, exchangeName, routingKey, body string, messageConfig ...MessageConfig) error {
	cfg := extractMessageConfig(messageConfig...)
	return bpi.Ch.PublishWithContext(
		ctx,
		exchangeName,
		routingKey,
		cfg.Mandatory, //mandatory
		cfg.Immediate, //immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
}
func extractMessageConfig(messageConfig ...MessageConfig) MessageConfig {
	if len(messageConfig) == 0 {
		return DefaultMessageConfig()
	} else {
		return messageConfig[0]
	}
}
func (bpi *RProducer) Stop() {
	bpi.Ch.Close()
}
