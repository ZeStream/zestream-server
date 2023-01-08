package configs

import (
	"context"
	"log"
	"zestream-server/constants"

	rmq "github.com/rabbitmq/amqp091-go"
)

var rabbitConn *rmq.Connection
var rabbitCh *rmq.Channel
var rabbitQ rmq.Queue
var rabbitCtx context.Context
var rabbitCtxCancel context.CancelFunc

func InitRabbitMQ() (*rmq.Connection, *rmq.Channel, *rmq.Queue, *context.Context, context.CancelFunc) {
	var err error
	rabbitConn, err = rmq.Dial(EnvVar[RABBITMQ_URI])
	failOnError(err)

	rabbitCh, err = rabbitConn.Channel()
	failOnError(err)

	rabbitQ, err = rabbitCh.QueueDeclare(
		constants.RABBIT_MQ_CHANNEL,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err)

	rabbitCtx, rabbitCtxCancel = context.WithTimeout(context.Background(), constants.RABBIT_MQ_TIMEOUT)

	return rabbitConn, rabbitCh, &rabbitQ, &rabbitCtx, rabbitCtxCancel
}

func GetRabbitMQ() (*rmq.Connection, *rmq.Channel, *rmq.Queue, *context.Context, context.CancelFunc) {
	return rabbitConn, rabbitCh, &rabbitQ, &rabbitCtx, rabbitCtxCancel
}

func failOnError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
