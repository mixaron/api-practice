package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"os"
)

func UploadFile(bucket string, objectName string, file multipart.File, fileSize int64, contentType string) (string, error) {
	ctx := context.Background()

	_, err := Client.PutObject(ctx, bucket, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s/%s/%s", os.Getenv("MINIO_URL"), bucket, objectName), nil
}
