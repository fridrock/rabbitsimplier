package rabbitsimplier

import amqp "github.com/rabbitmq/amqp091-go"

type Consumer interface {
	CreateChannel(*amqp.Connection) error
	CreateQueue(queueName ...string) (amqp.Queue, error)
	SetBinding(q amqp.Queue, boundingKey, exchangeName string) error
	RegisterHandler(amqp.Queue, Dispatcher)
	Stop()
}

type RConsumer struct {
	Ch *amqp.Channel
}

func (bci *RConsumer) CreateChannel(conn *amqp.Connection) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	bci.Ch = channel
	return nil
}

type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

func DefaultQueueConfig() QueueConfig {
	return QueueConfig{
		Name:       "",
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
}

func (bci *RConsumer) CreateQueue(queueConfig ...QueueConfig) (amqp.Queue, error) {
	var q amqp.Queue
	cfg := extractQueueConfig(queueConfig...)
	q, err := bci.Ch.QueueDeclare(
		cfg.Name,
		cfg.Durable,
		cfg.AutoDelete,
		cfg.Exclusive,
		cfg.NoWait,
		cfg.Args,
	)
	if err != nil {
		return q, err
	}
	return q, nil
}

func extractQueueConfig(qc ...QueueConfig) QueueConfig {
	if len(qc) == 0 {
		return DefaultQueueConfig()
	} else {
		return qc[0]
	}
}

type BindingConfig struct {
	NoWait bool
	Args   amqp.Table
}

func (bci *RConsumer) SetBinding(q amqp.Queue, boundingKey, exchangeName string, bindingConfig ...BindingConfig) error {
	cfg := extractBindingConfig(bindingConfig...)
	return bci.Ch.QueueBind(
		q.Name,
		boundingKey,
		exchangeName,
		cfg.NoWait,
		cfg.Args,
	)

}

func extractBindingConfig(bindingConfig ...BindingConfig) BindingConfig {
	if len(bindingConfig) == 0 {
		return BindingConfig{}
	} else {
		return bindingConfig[0]
	}
}

type ConsumerConfig struct {
	ConsumerName string
	AutoAck      bool
	Exclusive    bool
	NoLocal      bool
	NoWait       bool
	Args         amqp.Table
}

func DefaultConsumerConfig() ConsumerConfig {
	return ConsumerConfig{
		ConsumerName: "",
		AutoAck:      true,
		Exclusive:    false,
		NoLocal:      false,
		NoWait:       false,
		Args:         nil,
	}
}
func (bci *RConsumer) RegisterDispatcher(q amqp.Queue, bh Dispatcher, consumerConfig ...ConsumerConfig) error {
	cfg := extractConsumerConfig(consumerConfig...)
	msgs, err := bci.Ch.Consume(
		q.Name,
		cfg.ConsumerName,
		cfg.AutoAck,
		cfg.Exclusive,
		cfg.NoLocal,
		cfg.NoWait,
		cfg.Args,
	)
	if err != nil {
		return err
	}
	go bh.Dispatch(msgs)

	return nil
}
func extractConsumerConfig(consumerConfig ...ConsumerConfig) ConsumerConfig {
	if len(consumerConfig) == 0 {
		return DefaultConsumerConfig()
	} else {
		return consumerConfig[0]
	}
}
func (bci *RConsumer) Stop() {
	bci.Ch.Close()
}
