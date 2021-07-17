package main

import (
	"context"
	"fmt"
	"github.com/ag-students/support-service/config"
	"github.com/ag-students/support-service/pkg/kafka-impl"
	"github.com/ag-students/support-service/pkg/pdf-creator"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"log"
	"time"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	time.Sleep(time.Second * 1)
	fmt.Println("Hello, World! I generate docs")
	config.Init()

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
	kafka_impl.Write(ctx, writer, "Oh", "Test passed")

	reader := kafka_impl.GetKafkaReader(kafkaURL, topic, groupID)
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			log.Fatal("Error while closing reader")
		}
	}(reader)
	go kafka_impl.Listen(ctx, reader)

	surname := "Иванов"
	name := "Иван"
	patronymic := "Иванович"
	pdf_creator.CreatePDF(surname, name, patronymic)


	endpoint := "172.23.0.2:9000"
    accessKeyID := "minio"
    secretAccessKey := "minio123"
    useSSL := false


    // Initialize minio client object.
    minioClient, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
        Secure: useSSL,
    })
    if err != nil {
        log.Fatalln(err)
    }

    // Make a new bucket called mymusic.
    bucketName := "passports"
    location := "my_region"

    err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
    if err != nil {
        // Check to see if we already own this bucket (which happens if you run this twice)
        exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
        if errBucketExists == nil && exists {
            log.Printf("We already own %s\n", bucketName)
        } else {
            log.Fatalln(err)
        }
    } else {
        log.Printf("Successfully created %s\n", bucketName)
    }

    // Upload the zip file
    objectName := "passport.pdf"
    filePath := "./passport.pdf"
    contentType := "application/pdf"

    // Upload the zip file with FPutObject
    info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
    if err != nil {
        log.Fatalln(err)
    }

    log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}
