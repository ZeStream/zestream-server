package utils

import (
	"context"
	"time"
	"zestream-server/logger"

	"github.com/segmentio/kafka-go"
)

func PublishMessage(ctx context.Context, kafkaURI, topic string, message string) error {
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaURI, topic, partition)
	if err != nil {
		logger.Error(ctx, "failed to fial leader", logger.Z{
			"error":         err.Error(),
			"kafka_uri":     kafkaURI,
			"topic":         topic,
			"input_message": message,
		})

		return err
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte(message)},
	)
	if err != nil {
		logger.Error(ctx, "failed to write messages", logger.Z{
			"error":         err.Error(),
			"kafka_uri":     kafkaURI,
			"topic":         topic,
			"input_message": message,
		})

		return err
	}

	if err := conn.Close(); err != nil {
		logger.Error(ctx, "failed to close writer", logger.Z{
			"error":         err.Error(),
			"kafka_uri":     kafkaURI,
			"topic":         topic,
			"input_message": message,
		})

		return err
	}

	return nil
}
