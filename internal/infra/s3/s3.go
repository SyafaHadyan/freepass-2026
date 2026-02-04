// Package s3 handles s3 file upload
package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/SyafaHadyan/freepass-2026/internal/infra/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type S3Itf interface {
	Upload(ctx context.Context, objectKey string, object []byte) error
}

type S3 struct {
	Client          *s3.Client
	bucketName      string
	accountID       string
	accessKeyID     string
	accessKeySecret string
}

func New(env *env.Env) *S3 {
	config, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(env.S3AccessKeyID, env.S3AccessKeySecret, "")),
		config.WithRegion("auto"))
	if err != nil {
		log.Panic(err)
	}

	client := s3.NewFromConfig(config, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", env.S3AccountID))
	})

	S3 := S3{
		Client:          client,
		bucketName:      env.S3BucketName,
		accountID:       env.S3AccountID,
		accessKeyID:     env.S3AccessKeyID,
		accessKeySecret: env.S3AccessKeySecret,
	}

	return &S3
}

func (s *S3) Upload(ctx context.Context, objectKey string, object []byte) error {
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(object),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
			log.Printf("S3: %v\n", "EntityTooLarge")
		} else {
			log.Printf("S3: %v\n", "can't upload file")
		}
	}

	return err
}
