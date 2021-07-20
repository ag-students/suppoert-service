package services

import (
	"github.com/ag-students/support-service/internal/microservices/communication/models"
	"github.com/ag-students/support-service/internal/microservices/communication/repository"
	messagebird "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/sms"
	"log"
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
		log.Panic("Cannot send SMS")
		return err
	}

	log.Println("Sent SMS to " + msg.PhoneNumber)

	return nil
}
