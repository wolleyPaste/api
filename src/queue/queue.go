package queue

import (
	"fmt"
	"wolley-api/src/common"

	"github.com/charmbracelet/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	RabbitConn *amqp.Connection
	RabbitChan *amqp.Channel
)

func InitRabbitMQ(rabbitMQ common.RabbitMQ) {
	var err error
	RabbitConn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQ.Username, rabbitMQ.Password, rabbitMQ.Host, rabbitMQ.Port))
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	RabbitChan, err = RabbitConn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
}

func CleanupRabbitMQ() {
	if RabbitChan != nil {
		RabbitChan.Close()
	}
	if RabbitConn != nil {
		RabbitConn.Close()
	}
}
