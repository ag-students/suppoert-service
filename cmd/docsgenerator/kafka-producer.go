package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"log"
)

func Write(ctx context.Context, key, value string) {
	kafkaURL := viper.GetString("KAFKA_URL")
	topic := viper.GetString("TOPIC")
	w := &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	err := w.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(value),
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
