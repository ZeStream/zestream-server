package utils

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func PublishMessage(kafkaURI, topic string, message string) string {

	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaURI, topic, partition)
	if err != nil {
		log.Println("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte(message)},
	)
	if err != nil {
		log.Println("failed to write messages:", err)
		return err.Error()
	}

	if err := conn.Close(); err != nil {
		log.Println("failed to close writer:", err)
		return err.Error()
	}

	return "success"
}
