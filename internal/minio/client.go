package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

var Client *minio.Client

func Init() {
	var err error
	Client, err = minio.New(os.Getenv("MINIO_URL"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_USER"), os.Getenv("MINIO_PASS"), ""),
		Secure: true,
	})

	if err != nil {
		log.Fatalf("init minio: %v", err)
	}
	bucketName := "articles"

	exists, err := Client.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatalf("ошибка при проверке bucket: %v", err)
	}

	if !exists {
		err = Client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("ошибка при создании bucket: %v", err)
		}
		log.Printf("Bucket %s создан", bucketName)
	} else {
		log.Printf("Bucket %s уже существует", bucketName)
	}

}
