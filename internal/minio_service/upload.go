package minio_service

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"os"
)

type UploadService interface {
	UploadFile(bucket string, objectName string, file multipart.File, fileSize int64, contentType string) (string, error)
}

type UploadServiceImpl struct {
	client *minio.Client
}

func NewUploadService(client *minio.Client) UploadService {
	return &UploadServiceImpl{client}
}

func (c *UploadServiceImpl) UploadFile(bucket string, objectName string, file multipart.File, fileSize int64, contentType string) (string, error) {
	ctx := context.Background()

	_, err := c.client.PutObject(ctx, bucket, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s/%s/%s", os.Getenv("MINIO_URL"), bucket, objectName), nil
}
