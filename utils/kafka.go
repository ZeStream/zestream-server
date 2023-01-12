package utils

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func PublishMessage(kafkaURI, topic string, message string) (string, error) {

	partition := 0
	m := "success"

	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaURI, topic, partition)
	if err != nil {
		m = "failed to dial leader"
		log.Println(m, err)
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
		m = "failed to write messages"
		log.Println(m, err)
		return m, err
	}

	if err := conn.Close(); err != nil {
		m = "failed to close writer"
		log.Println(m, err)
		return m, err
	}

	return m, nil
}
