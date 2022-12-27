package utils

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func PublishMessage(kafkaURI, topic string, message string) bool {

	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaURI, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte(message)},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
		return false
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
		return false
	}

	return true
}
