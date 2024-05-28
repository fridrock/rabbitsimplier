package rabbitsimplier

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Username string
	Password string
	Host     string
}
type Configurer interface {
	Configure(Config) error
	Stop()
}

type RConfigurer struct {
	Connection *amqp.Connection
}

func (bci *RConfigurer) Configure(config Config) error {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", config.Username, config.Password, config.Host))
	if err != nil {
		return err
	}
	bci.Connection = conn
	return nil
}

func (bci *RConfigurer) Stop() {
	bci.Connection.Close()
}
