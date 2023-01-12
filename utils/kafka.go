package utils

import (
	"context"
	"time"
	"zestream-server/logger"

	"github.com/segmentio/kafka-go"
)

func PublishMessage(ctx context.Context, kafkaURI, topic string, message string) (string, error) {
	partition := 0
	m := "success"

	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaURI, topic, partition)
	if err != nil {
		logger.Error(ctx, "failed to fial leader", logger.Z{
			"error":         err.Error(),
			"kafka_uri":     kafkaURI,
			"topic":         topic,
			"input_message": message,
		})

		return m, err
	}

	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		m = "failed to set write deadline"
		log.Println(m, err)
	}
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

		m = "failed to write messages"
		return m, err
	}

	if err := conn.Close(); err != nil {
		logger.Error(ctx, "failed to close writer", logger.Z{
			"error":         err.Error(),
			"kafka_uri":     kafkaURI,
			"topic":         topic,
			"input_message": message,
		})

		m = "failed to close writer"
		return m, err
	}

	return m, nil
}
