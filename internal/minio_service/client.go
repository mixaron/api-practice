package minio_service

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

type ClientInitService interface {
	Init() (*minio.Client, error)
}

type ClientServiceImpl struct {
}

func NewClientInitService() ClientInitService {
	return &ClientServiceImpl{}
}

func (c *ClientServiceImpl) Init() (*minio.Client, error) {
	client, err := minio.New(os.Getenv("MINIO_URL"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_USER"), os.Getenv("MINIO_PASS"), ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("init minio_service: %w", err)
	}

	bucketName := "articles"

	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке bucket: %w", err)
	}

	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("ошибка при создании bucket: %w", err)
		}
		log.Printf("Bucket %s создан", bucketName)
	} else {
		log.Printf("Bucket %s уже существует", bucketName)
	}

	return client, nil
}
