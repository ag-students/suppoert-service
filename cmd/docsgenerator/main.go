package main

import (
	"context"
	"fmt"
	"github.com/ag-students/support-service/config"
	"github.com/ag-students/support-service/internal/microservices/communication/repository/postgres"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
	"log"
)

func main() {
	log.Print("Starting docsgenerator service...")
	config.Init()

	log.Print("Connecting to the database...")
	dbParams := postgres.PSQLConfig{
		Host:     viper.GetString("db.postgres.host"),
		Port:     viper.GetString("db.postgres.port"),
		Password: viper.GetString("db.postgres.password"),
		DBName:   viper.GetString("db.postgres.database"),
		Username: viper.GetString("db.postgres.user"),
		SSLMode:  viper.GetString("db.postgres.sslmode"),
	}
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbParams.Username, dbParams.Password, dbParams.Host, dbParams.Port, dbParams.DBName))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		_ = conn.Close(ctx)
	}(conn, context.Background())

	log.Print("Connection established!")
}
