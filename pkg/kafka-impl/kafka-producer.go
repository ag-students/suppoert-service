package kafka_impl

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func NewKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func Write(ctx context.Context, writer *kafka.Writer, key string, value []byte) {
	fmt.Println("start producing ... !!")
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}
	err := writer.WriteMessages(ctx, msg)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("produced %s %s\n", key, value)
	}
}
