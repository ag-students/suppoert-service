package main

import (
	"github.com/ag-students/support-service/config"
	"github.com/ag-students/support-service/internal/microservices/communication/delivery/mq"
	"github.com/ag-students/support-service/internal/microservices/communication/repository"
	"github.com/ag-students/support-service/internal/microservices/communication/repository/postgres"
	"github.com/ag-students/support-service/internal/microservices/communication/services"
	"github.com/ag-students/support-service/utils"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger.InitLogger()

	logger.Logger.Info("Starting communication service...")

	config.Init()

	logger.Logger.Info("Connecting to the database...")
	conn, err := postgres.EstablishPSQLConnection(&postgres.PSQLConfig{
		Host:     viper.GetString("db.postgres.host"),
		Port:     viper.GetString("db.postgres.port"),
		Password: viper.GetString("db.postgres.password"),
		DBName:   viper.GetString("db.postgres.database"),
		Username: viper.GetString("db.postgres.user"),
		SSLMode:  viper.GetString("db.postgres.sslmode"),
	})
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Logger.Error(err.Error())
		}
	}()
	logger.Logger.Info("Connection established!")

	repo := repository.NewRepository(conn)

	time.Sleep(time.Second * 5)

	serv := &services.Service{
		SMSNotifier: services.NewMessageBird(repo, &services.MessageBirdConfig{
			AccessKey:  viper.GetString("communication-service.message-bird.access-key"),
			Originator: viper.GetString("communication-service.message-bird.originator"),
			Params:     viper.GetString("communication-service.message-bird.params"),
		}),
		EmailNotifier: services.NewEmailNotificator(repo, &services.SMTPConfig{
			Host:     viper.GetString("communication-service.smtp-config.host"),
			Port:     viper.GetInt("communication-service.smtp-config.port"),
			Username: viper.GetString("communication-service.smtp-config.username"),
			Password: viper.GetString("communication-service.smtp-config.password"),
		}),
	}

	consumers := &mq.KafkaConsumers{
		KafkaSMSConsumer: mq.NewKafkaSMSConsumer(serv, &mq.KafkaConfig{
			Brokers: []string{viper.GetString("communication-service.kafka.broker")},
			Topic:   viper.GetString("communication-service.kafka.sms-topic"),
		}),
		KafkaEmailConsumer: mq.NewKafkaEmailConsumer(serv, &mq.KafkaConfig{
			Brokers: []string{viper.GetString("communication-service.kafka.broker")},
			Topic:   viper.GetString("communication-service.kafka.email-topic"),
		}),
	}

	consumers.StartConsumers()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Logger.Info("Service gracefully stopped.")
}
