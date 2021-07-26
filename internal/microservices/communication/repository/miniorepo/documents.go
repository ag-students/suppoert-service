package miniorepo

import (
	"context"
	"github.com/minio/minio-go/v7"
)

type FileServerConfig struct {
	Client *minio.Client
	Bucket string
}

type DocumentsMinio struct {
	client *minio.Client
	bucketName string
}

func NewDocumentsMinio(cnf *FileServerConfig) *DocumentsMinio {
	return &DocumentsMinio{
		client: cnf.Client,
		bucketName: cnf.Bucket,
	}
}

func (r *DocumentsMinio) GetDocument(docHash string) (string, error) {
	ctx := context.Background()
	if err := r.client.FGetObject(ctx, r.bucketName, docHash, "./" + docHash, minio.GetObjectOptions{}); err != nil {
		return "", err
	}

	return "./" + docHash, nil
}
