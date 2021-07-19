package services

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"context"
	"os"
)

	
    // Initialize minio client object.
	func NewClient() *minio.Client {
		endpoint := "172.19.0.1:9000"
		accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
		secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
		useSSL := false
    		minioClient, err := minio.New(endpoint, &minio.Options{
       		 Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
       			 Secure: useSSL,
   			 })
		if err != nil {
			log.Fatalln(err)
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
            log.Fatalln(err)
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
			log.Fatalln(err)
		}
		log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	}

    