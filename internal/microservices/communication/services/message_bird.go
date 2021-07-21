package services

import (
	"github.com/ag-students/support-service/internal/microservices/communication/models"
	"github.com/ag-students/support-service/internal/microservices/communication/repository"
	"github.com/ag-students/support-service/utils"
	messagebird "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/sms"
)

type MessageBirdConfig struct {
	AccessKey  string
	Originator string
	Params     string
}

type MessageBird struct {
	repo   *repository.Repository
	client *messagebird.Client
	config *MessageBirdConfig
}

func NewMessageBird(repo *repository.Repository, cnf *MessageBirdConfig) *MessageBird {
	return &MessageBird{
		repo:   repo,
		client: messagebird.New(cnf.AccessKey),
		config: cnf,
	}
}

func (r *MessageBird) NotifyBySMS(msg *models.SMSMessage) error {
	params := &sms.Params{Reference: r.config.Params}

	if _, err := sms.Create(
		r.client,
		r.config.Originator,
		[]string{msg.PhoneNumber},
		msg.Body,
		params,
	); err != nil {
		logger.Logger.Errorf("cant serve sms communication request for %s number", msg.PhoneNumber)
		return err
	}

	logger.Logger.Info("Sent SMS to " + msg.PhoneNumber)

	comm := models.Communication{
		CommunicationType: "SMS",
		Delayed:           false,
		Phone:             msg.PhoneNumber,
	}

	if _, err := r.repo.CommunicationHistoryRepository.CreateCommunication(comm); err != nil {
		logger.Logger.Errorf("cant create history record of sms communication: %s", err.Error())
		return err
	}

	return nil
}
