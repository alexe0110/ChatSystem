package storage

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client     *s3.Client
	bucketName string
	region     string
}

func NewS3Storage(bucketName, region string) *S3Storage {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
	)

	if err != nil {
		log.Fatalln(err)
	}

	return &S3Storage{
		client:     s3.NewFromConfig(cfg),
		bucketName: bucketName,
		region:     region,
	}
}

func (s *S3Storage) Upload(ctx context.Context, fileName string, data []byte) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return "", fmt.Errorf("s3 upload failed: %w", err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		s.bucketName, s.region, fileName)
	return url, nil
}
