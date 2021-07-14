package services

import (
	"github.com/ag-students/support-service/internal/microservices/communication/models"
)

type SMSNotifier interface {
	NotifyBySMS(msg *models.SMSMessage) error
}

type EmailNotifier interface {
	NotifyByEmail(msg *models.EmailMessage) (int, error)
}

type Service struct {
	SMSNotifier
	EmailNotifier
}
