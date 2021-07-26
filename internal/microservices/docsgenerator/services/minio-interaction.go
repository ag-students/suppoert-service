package services

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

// Initialize minio client object.
func NewClient() *minio.Client {

	//Warning! Can change after make app-setup-and-up
	endpoint := "172.19.0.2:9001"

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
			log.Fatalln("Failed to make bucket: ", err)
		}
	} else {
		log.Printf("Successfully maked %s\n", bucketName)
	}

	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"admin:*"
				]
			},
			{
				"Effect": "Allow",
				"Action": [
					"s3:*"
				],
				"Resource": [
					"arn:aws:s3:::*"
				]
			}
		]
	}`

	err = minioClient.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		log.Println("Failed to set bucket policy:", err)
	}
}

// Upload the file
func UploadNewFile(objectName string) (pdfLink string) {
	minioClient := NewClient()
	ctx := context.Background()
	bucketName := "passports"
	filePath := "./" + objectName
	pdfLink = "localhost:9000/passports/" + objectName
	contentType := "application/" + objectName[len(objectName)-3:]
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln("Failed to upload:", err)
	} else {
		log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	}
	return pdfLink
}

// Fetch metadata of an object
func StatObject(objectName string) {
	minioClient := NewClient()
	ctx := context.Background()
	bucketName := "passports"
	objInfo, err := minioClient.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("Failed to fetch metadata of object:", err)
	}
	fmt.Println(objInfo)
}
