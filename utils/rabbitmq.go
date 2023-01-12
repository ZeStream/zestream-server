package utils

import (
	"context"

	rmq "github.com/rabbitmq/amqp091-go"
)

func PublishEvent(ch *rmq.Channel, q *rmq.Queue, ctx context.Context, body []byte) error {
	err := ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		rmq.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)

	return err
}
