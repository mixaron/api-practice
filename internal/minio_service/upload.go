package minio_service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type UploadService interface {
	UploadFile(bucket string, objectName string, file multipart.File, fileSize int64, contentType string) (string, error)
}

type UploadServiceImpl struct {
	client *s3.Client
}

func NewUploadService(client *s3.Client) *UploadServiceImpl {
	return &UploadServiceImpl{client}
}

func (s *UploadServiceImpl) UploadFile(bucket string, objectName string, file multipart.File, fileSize int64, contentType string) (string, error) {
	defer file.Close()

	ctx := context.TODO()

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(objectName),
		Body:          file,
		ContentLength: &fileSize,
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки объекта: %w", err)
	}

	endpoint := os.Getenv("YANDEX_STORAGE_ENDPOINT")
	url := fmt.Sprintf("https://%s/%s/%s", endpoint, bucket, objectName)
	return url, nil
}
