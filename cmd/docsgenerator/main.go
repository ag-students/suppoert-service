package main

import (
	"context"
	"fmt"
	"github.com/ag-students/support-service/config"
	"github.com/ag-students/support-service/pkg/kafka-impl"
	"github.com/ag-students/support-service/pkg/pdf-creator"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"log"
	"time"
)

func main() {
	time.Sleep(time.Second * 10)
	fmt.Println("Hello, World! I generate docs")
	config.Init()

	ctx := context.Background()
	kafkaURL := viper.GetString("mq.kafka.url")
	topic := viper.GetString("mq.kafka.topic")
	groupID := viper.GetString("mq.kafka.group_id")

	writer := kafka_impl.NewKafkaWriter(kafkaURL, topic)
	defer func(writer *kafka.Writer) {
		err := writer.Close()
		if err != nil {
			log.Fatal("Error while closing writer")
		}
	}(writer)
	kafka_impl.Write(ctx, writer, "Oh", "Test passed")

	reader := kafka_impl.GetKafkaReader(kafkaURL, topic, groupID)
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			log.Fatal("Error while closing reader")
		}
	}(reader)
	go kafka_impl.Listen(ctx, reader)

	surname := "Иванов"
	name := "Иван"
	patronymic := "Иванович"
	pdf_creator.CreatePDF(surname, name, patronymic)
}
