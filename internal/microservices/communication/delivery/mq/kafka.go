package mq

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
	"time"
)

type KafkaConsumers struct {
	KafkaSMSConsumer
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

func GetKafkaReader(cnf *KafkaConfig) *kafka.Reader {
	fmt.Println("CONFIG: " + cnf.Brokers[0])

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
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		if err := r.KafkaSMSConsumer.ConsumeSMSRequests(); err != nil {
			log.Panic("error while starting consumers: " + err.Error())
		}
		defer wg.Done()
	}()
	wg.Wait()

	defer func() {
		if err := r.KafkaSMSConsumer.reader.Close(); err != nil {
			log.Printf("error while trying close sms consumer: " + err.Error())
		}
	}()
}
