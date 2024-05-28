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
	GetConnection() *amqp.Connection
}

type RConfigurer struct {
	connection *amqp.Connection
}

func (bci *RConfigurer) Configure(config Config) error {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", config.Username, config.Password, config.Host))
	if err != nil {
		return err
	}
	bci.connection = conn
	return nil
}
func (bci RConfigurer) GetConnection() *amqp.Connection {
	return bci.connection
}
func (bci *RConfigurer) Stop() {
	bci.connection.Close()
}
