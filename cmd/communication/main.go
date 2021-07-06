package main

import (
	"github.com/ag-students/support-service/config"
	"github.com/ag-students/support-service/internal/microservices/communication/repository"
	"github.com/ag-students/support-service/internal/microservices/communication/repository/postgres"
	"github.com/spf13/viper"
	"log"
)

func main() {
	log.Print("Starting communication service...")

	config.Init()

	log.Print("Connecting to the database...")
	conn, err := postgres.EstablishPSQLConnection(&postgres.PSQLConfig{
		Host: viper.GetString("db.postgres.host"),
		Port: viper.GetString("db.postgres.port"),
		Password: viper.GetString("db.postgres.password"),
		DBName: viper.GetString("db.postgres.database"),
		Username: viper.GetString("db.postgres.user"),
		SSLMode: viper.GetString("db.postgres.sslmode"),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Connection established!")

	repository.NewRepository(conn)

	log.Print("Service gracefully stopped.")
}
