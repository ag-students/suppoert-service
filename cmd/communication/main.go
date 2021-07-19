package main

import (
	"fmt"
	"github.com/ag-students/support-service/config"
	"github.com/ag-students/support-service/internal/microservices/communication/delivery/mq"
	"github.com/ag-students/support-service/internal/microservices/communication/repository"
	"github.com/ag-students/support-service/internal/microservices/communication/repository/postgres"
	"github.com/ag-students/support-service/internal/microservices/communication/services"
	"github.com/spf13/viper"
	"log"
	"time"
)

func main() {
	log.Print("Starting communication service...")

	config.Init()

	log.Print("Connecting to the database...")
	conn, err := postgres.EstablishPSQLConnection(&postgres.PSQLConfig{
		Host:     viper.GetString("db.postgres.host"),
		Port:     viper.GetString("db.postgres.port"),
		Password: viper.GetString("db.postgres.password"),
		DBName:   viper.GetString("db.postgres.database"),
		Username: viper.GetString("db.postgres.user"),
		SSLMode:  viper.GetString("db.postgres.sslmode"),
	})
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("error while closing connection: " + err.Error())
		}
	}()
	log.Print("Connection established!")

	repo := repository.NewRepository(conn)

	fmt.Println(repo)

	time.Sleep(time.Second * 5)

	serv := &services.Service{
		SMSNotifier: services.NewMessageBird(repo, &services.MessageBirdConfig{
			AccessKey:  viper.GetString("communication-service.message-bird.access-key"),
			Originator: viper.GetString("communication-service.message-bird.originator"),
			Params:     viper.GetString("communication-service.message-bird.params"),
		}),
		EmailNotifier: nil,
	}

	log.Println(serv)

	time.Sleep(time.Second * 6)
	log.Print("Connecting to kafka...")

	consumers := &mq.KafkaConsumers{
		KafkaSMSConsumer: mq.NewKafkaSMSConsumer(serv, &mq.KafkaConfig{
			Brokers: []string{viper.GetString("communication-service.kafka.broker")},
			Topic:   viper.GetString("communication-service.kafka.sms-topic"),
		}),
	}

	consumers.StartConsumers()

	log.Print("Service gracefully stopped.")
}
