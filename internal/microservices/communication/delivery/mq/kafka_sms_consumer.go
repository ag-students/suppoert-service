package mq

import (
	"context"
	"encoding/json"
	"github.com/ag-students/support-service/internal/microservices/communication/models"
	"github.com/ag-students/support-service/internal/microservices/communication/services"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaSMSConsumer struct {
	serv   services.SMSNotifier
	cnf    *KafkaConfig
	reader *kafka.Reader
}

func NewKafkaSMSConsumer(serv *services.Service, cnf *KafkaConfig) KafkaSMSConsumer {
	return KafkaSMSConsumer{
		serv: serv.SMSNotifier,
		cnf:  cnf,
	}
}

func (r *KafkaSMSConsumer) ConsumeSMSRequests() error {
	r.reader = GetKafkaReader(r.cnf)

	for {
		m, err := r.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error while receiving message: %s", err.Error())
			continue
		}

		var smsRequest models.SMSMessage

		if err = json.Unmarshal(m.Value, &smsRequest); err != nil {
			log.Printf("json message is invalid: %s", err.Error())
			continue
		}

		if err = r.serv.NotifyBySMS(&smsRequest); err != nil {
			log.Printf("cant serve sms request: %s", err.Error())
			continue
		}

		log.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Value))
	}
}
