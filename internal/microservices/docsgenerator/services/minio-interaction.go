package services

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"fmt"
)

// Initialize minio client object.
func NewClient() *minio.Client {

	//Warning! Can change after make app-setup-and-up
	endpoint := "172.18.0.2:9000"

	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("Failed to create minio client:", err)
	} else {
		fmt.Println("Successfully created minio client")
	}
	return minioClient
}

// Make a new bucket
func NewBucket(bucketName string) {
	minioClient := NewClient()
	location := "my_region"
	ctx := context.Background()
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln("Failed to create bucket: ", err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
}

// Upload the file
func UploadNewFile(objectName string) {
	minioClient := NewClient()
	ctx := context.Background()
	bucketName := "passports"
	filePath := "./" + objectName
	contentType := "application/" + objectName[len(objectName)-3:]
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln("Failed to upload:", err)
	} else {
		log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	}
}
