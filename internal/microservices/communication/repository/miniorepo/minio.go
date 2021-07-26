package miniorepo

import (
	logger "github.com/ag-students/support-service/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioConf struct {
	EndPoint 	string
	AccessKey 	string
	Secret 		string
	UseSSL		bool
}

func EstablishMinioConnection(conf *MinioConf) (*minio.Client, error) {
	minioClient, err := minio.New(conf.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.Secret, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		logger.Logger.Errorf("Failed to create miniorepo client:", err)
	} else {
		logger.Logger.Info("Successfully created miniorepo client")
	}

	return minioClient, err
}
