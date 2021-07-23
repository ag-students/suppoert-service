package mq

import (
	"context"
	"encoding/json"
	"github.com/ag-students/support-service/internal/microservices/communication/models"
	"github.com/ag-students/support-service/internal/microservices/communication/services"
	"github.com/ag-students/support-service/utils"
	"github.com/segmentio/kafka-go"
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
			logger.Logger.Errorf("error while receiving message: %s", err.Error())
			continue
		}

		logger.Logger.Infof("handling incoming sms request: %s", m.Value)

		var smsRequest models.SMSMessage

		if err = json.Unmarshal(m.Value, &smsRequest); err != nil {
			logger.Logger.Errorf("sms request is invalid: %s", err.Error())
			continue
		}

		if err = r.serv.NotifyBySMS(&smsRequest); err != nil {
			logger.Logger.Errorf("cant serve sms request: %s", err.Error())
			continue
		}

		logger.Logger.Infof("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Value))
	}
}
