package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ag-students/support-service/config"
	"github.com/ag-students/support-service/internal/microservices/docsgenerator/services"
	"github.com/ag-students/support-service/pkg/kafka-impl"
	"github.com/ag-students/support-service/pkg/pdf-creator"
	"github.com/dchest/uniuri"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	// "time"
)

func main() {
	// time.Sleep(time.Second * 5)
	fmt.Println("Hello, World! I generate docs")
	config.Init()

	newPatient := &pdf_creator.PatientPersonalData{
		Surname:     uniuri.NewLen(10),
		Name:        "Иван",
		Patronymic:  "Иванович",
		Birthday:    "01.01.2000",
		Gender:      "мужской",
		HomeAddress: "Москва, Красная 213",
		FirstDate:   "01.08.2021",
		SecondDate:  "22.08.2021",
		Vaccine:     "Спутник-V",
		PdfName:     "passport.pdf",
	}

	//create PDF file
	pdf_creator.CreatePDF(newPatient)

	//Uppload PDF file to minIO
	services.UploadNewFile(newPatient.PdfName)

	ctx := context.Background()
	kafkaURL := viper.GetString("mq.kafka.url")
	topic := viper.GetString("mq.kafka.topic")
	groupID := viper.GetString("mq.kafka.group_id")

	writer := kafka_impl.NewKafkaWriter(kafkaURL, topic)
	defer func(writer *kafka.Writer) {
		err := writer.Close()
		if err != nil {
			log.Fatal("Error while closing writer")
		}
	}(writer)
	patientDataJson, err := json.Marshal(newPatient)
	if err != nil {
		log.Fatal(err)
	}
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)
	kafka_impl.Write(ctx, writer, "DataForDocument", patientDataJson)

	reader := kafka_impl.GetKafkaReader(kafkaURL, topic, groupID)
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			log.Fatal("Error while closing reader")
		}
	}(reader)
	go kafka_impl.Listen(ctx, reader)
}
