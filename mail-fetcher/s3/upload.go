package s3

import (
	"bytes"
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader struct {
	client *s3.Client
	bucket string
}

func NewS3Uploader(bucket string) *S3Uploader {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	return &S3Uploader{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
	}
}

func (u *S3Uploader) UploadFile(filePath string, key string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, err = u.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &u.bucket,
		Key:    &key,
		Body:   bytes.NewReader(data),
	})
	return err
}
