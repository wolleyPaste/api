package queue

import (
	"fmt"
	"wolley-api/src/common"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	RabbitConn *amqp.Connection
	RabbitChan *amqp.Channel
)

func InitRabbitMQ(rabbitMQ common.RabbitMQ) error {
	var err error
	RabbitConn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQ.Username, rabbitMQ.Password, rabbitMQ.Host, rabbitMQ.Port))
	if err != nil {
		return err
	}

	RabbitChan, err = RabbitConn.Channel()
	if err != nil {
		return err
	}

	return nil
}

func CleanupRabbitMQ() {
	if RabbitChan != nil {
		RabbitChan.Close()
	}
	if RabbitConn != nil {
		RabbitConn.Close()
	}
}
