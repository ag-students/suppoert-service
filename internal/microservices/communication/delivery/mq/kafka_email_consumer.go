package mq

import (
	"context"
	"encoding/json"
	"github.com/ag-students/support-service/internal/microservices/communication/models"
	"github.com/ag-students/support-service/internal/microservices/communication/services"
	"github.com/ag-students/support-service/utils"
	"github.com/segmentio/kafka-go"
)

type KafkaEmailConsumer struct {
	serv 	services.EmailNotifier
	cnf  	*KafkaConfig
	reader 	*kafka.Reader
}

func NewKafkaEmailConsumer(serv *services.Service, cnf *KafkaConfig) KafkaEmailConsumer {
	return KafkaEmailConsumer{
		serv: 	serv.EmailNotifier,
		cnf: 	cnf,
	}
}

func (r *KafkaEmailConsumer) ConsumeEmailRequests() error {
	r.reader = GetKafkaReader(r.cnf)

	for {
		m, err := r.reader.ReadMessage(context.Background())
		if err != nil {
			logger.Logger.Errorf("error while receiving message: %s", err.Error())
			continue
		}

		logger.Logger.Infof("handling incoming email request: %s", m.Value)

		var emailRequest models.EmailMessage

		if err = json.Unmarshal(m.Value, &emailRequest); err != nil {
			logger.Logger.Errorf("email request is invalid: %s", err.Error())
			continue
		}

		if err = r.serv.NotifyByEmail(&emailRequest); err != nil {
			logger.Logger.Errorf("cant serve email request: %s", err.Error())
			continue
		}

		logger.Logger.Infof("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Value))
	}
}
