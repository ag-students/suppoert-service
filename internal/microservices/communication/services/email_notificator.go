package services

import (
	"crypto/tls"
	"github.com/ag-students/support-service/internal/microservices/communication/models"
	"github.com/ag-students/support-service/internal/microservices/communication/repository"
	"github.com/ag-students/support-service/internal/microservices/communication/services/email_templates"
	logger "github.com/ag-students/support-service/utils"
	"github.com/xhit/go-simple-mail/v2"
	"os"
	"time"
)

type SMTPConfig struct {
	Host 		string
	Port 		int
	Username 	string
	Password 	string
}

type EmailNotificator struct {
	repo 		*repository.Repository
	smtpClient 	*mail.SMTPClient
}

func NewEmailNotificator(repo *repository.Repository, smtpConfig *SMTPConfig) *EmailNotificator {
	server := mail.NewSMTPClient()

	server.Host = smtpConfig.Host
	server.Port = smtpConfig.Port
	server.Username = smtpConfig.Username
	server.Password = smtpConfig.Password
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	client, err := server.Connect()
	if err != nil {
		logger.Logger.Errorf("error while creating email notificator: %s", err.Error())
	}

	return &EmailNotificator{
		repo: repo,
		smtpClient: client,
	}
}

func (r *EmailNotificator) NotifyByEmail(msg *models.EmailMessage) error {
	email := mail.NewMSG()
	email.SetFrom("Covid-19 Protection <noreply@shirnis-koronoy.ga>").
		AddTo(msg.EmailAddress).
		SetSubject("Ваш код подтверждения")
	email.SetBody(mail.TextHTML, email_templates.ConfirmEmailTemplate(msg.ReplyLink))

	if email.Error != nil {
		logger.Logger.Error(email.Error)
	}

	comm := models.Communication{
		CommunicationType: "EMAIL",
		Delayed:           false,
		Email: 			   msg.EmailAddress,
	}

	if _, err := r.repo.CommunicationHistoryRepository.CreateCommunication(comm); err != nil {
		logger.Logger.Errorf("cant create history record of sms communication: %s", err.Error())
		return err
	}

	if err := email.Send(r.smtpClient); err != nil {
		logger.Logger.Error(err.Error())
		return err
	} else {
		logger.Logger.Info("Email Sent")
		return nil
	}
}

func (r *EmailNotificator) NotifyByEmailPass(msg *models.EmailPassportMessage) error {
	email := mail.NewMSG()
	email.SetFrom("Covid-19 Protection <noreply@shirnis-koronoy.ga>").
		AddTo(msg.EmailAddress).
		SetSubject("Ваш код подтверждения")
	email.SetBody(mail.TextHTML, email_templates.ConfirmEmailTemplate(msg.ReplyLink))

	root, err := os.Getwd()
	if err != nil {
		logger.Logger.Errorf("cant stat root dir: %s", err.Error())
	}

	email.Attach(&mail.File{FilePath: root + "tmp/" + msg.FileName, Name:"passport.pdf", Inline: true})

	if email.Error != nil {
		logger.Logger.Error(email.Error)
	}

	comm := models.Communication{
		CommunicationType: "EMAIL",
		Delayed:           false,
		Email: 			   msg.EmailAddress,
	}

	if _, err := r.repo.CommunicationHistoryRepository.CreateCommunication(comm); err != nil {
		logger.Logger.Errorf("cant create history record of sms communication: %s", err.Error())
		return err
	}

	if err := email.Send(r.smtpClient); err != nil {
		logger.Logger.Error(err.Error())
		return err
	} else {
		logger.Logger.Info("Email Sent")
		return nil
	}
}
