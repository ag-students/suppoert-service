package mq

import (
	logger "github.com/ag-students/support-service/utils"
	"github.com/segmentio/kafka-go"
	"time"
)

type KafkaConsumers struct {
	KafkaSMSConsumer
	KafkaEmailConsumer
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

func GetKafkaReader(cnf *KafkaConfig) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:         cnf.Brokers,
		Topic:           cnf.Topic,
		MinBytes:        10e3,
		MaxBytes:        10e6,
		MaxWait:         1 * time.Second,
		ReadLagInterval: -1,
	})
}

func (r *KafkaConsumers) StartConsumers() {
	go func() {
		if err := r.KafkaSMSConsumer.ConsumeSMSRequests(); err != nil {
			logger.Logger.Fatalf("error while starting consumers: %s", err.Error())
		}

		defer func() {
			if err := r.KafkaSMSConsumer.reader.Close(); err != nil {
				logger.Logger.Errorf("error while trying close sms consumer: %s", err.Error())
			}
		}()
	}()

	go func() {
		if err := r.KafkaEmailConsumer.ConsumeEmailRequests(); err != nil {
			logger.Logger.Fatalf("error while starting email consumer: %s", err.Error())
		}
		defer func() {
			if err := r.KafkaEmailConsumer.reader.Close(); err != nil {
				logger.Logger.Errorf("error while trying close email consumer: %s", err.Error())
			}
		}()
	}()
}
