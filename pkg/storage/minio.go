package storage

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
	endpoint   string
}

func NewMinioStorage(endpoint, user, password, bucket string) *MinioStorage {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(user, password, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		log.Fatalln(err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln(err)
		}
	}

	return &MinioStorage{
		client:     minioClient,
		bucketName: bucket,
		endpoint:   endpoint,
	}
}

func (s *MinioStorage) Upload(ctx context.Context, fileName string, data []byte) (string, error) {
	reader := bytes.NewReader(data)
	_, err := s.client.PutObject(ctx, s.bucketName, fileName, reader, int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}

	url := fmt.Sprintf("http://%s/%s/%s", s.endpoint, s.bucketName, fileName)
	return url, nil
}
