package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func Listen(ctx context.Context) {
	// get kafka reader using environment variables.
	kafkaURL := viper.GetString("KAFKA_URL")
	topic := viper.GetString("TOPIC")
	groupID := viper.GetString("GROUP_ID")

	reader := getKafkaReader(kafkaURL, topic, groupID)

	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			log.Fatal("failed to close reader:", err)
		}
	}(reader)

	fmt.Println("start consuming ... !!")
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n",
			m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
