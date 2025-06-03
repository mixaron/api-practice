package minio_service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
)

type ClientInitService interface {
	InitYandexClient() (*s3.Client, error)
}

type ClientServiceImpl struct{}

func NewClientInitService() ClientInitService {
	return &ClientServiceImpl{}
}

func (c *ClientServiceImpl) InitYandexClient() (*s3.Client, error) {
	endpoint := os.Getenv("YANDEX_STORAGE_ENDPOINT")
	region := "ru-central1"
	accessKeyID := os.Getenv("YANDEX_ACCESS_KEY")
	secretAccessKey := os.Getenv("YANDEX_SECRET_KEY")

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           fmt.Sprintf("https://%s", endpoint),
			SigningRegion: region,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки AWS конфигурации: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	return s3Client, nil
}
