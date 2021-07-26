package services

import (
	"context"
	logger "github.com/ag-students/support-service/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"os"
)

// Initialize miniorepo client object.
func NewClient() *minio.Client {

	//Warning! Can change after make app-setup-and-up
	endpoint := "nginx:9000"

	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Logger.Errorf("Failed to create miniorepo client: %s", err.Error())
	} else {
		logger.Logger.Info("Successfully created miniorepo client")
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
			logger.Logger.Info("we already own: %s", bucketName)
		} else {
			logger.Logger.Errorf("failed to create bucket: %s", err.Error())
		}
	} else {
		logger.Logger.Info("successfully created: %s", bucketName)
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
		logger.Logger.Errorf("Failed to upload:", err.Error())
	} else {
		logger.Logger.Infof("Successfully uploaded %s of size %d\n", objectName, info.Size)
	}
}
