package file_repository

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Repository struct {
	s3Client   *s3.Client
	bucketName string
	ctx        context.Context
}

func NewS3Repository(
	ctx context.Context,
	key string,
	secret string,
	endpoint string,
	bucketName string,
	region string,
) *S3Repository {
	creds := credentials.NewStaticCredentialsProvider(key, secret, "")
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
		config.WithBaseEndpoint(endpoint))

	if err != nil {
		fmt.Println("could not create s3 cfg")
		return nil
	}

	client := s3.NewFromConfig(cfg)

	if client == nil {
		fmt.Println("newly created client is nil")
		return nil
	}

	return &S3Repository{
		s3Client:   client,
		bucketName: bucketName,
		ctx:        ctx,
	}
}

func (repo *S3Repository) AddBinary(
	fileName string,
	fileData []byte,
) error {
	_, err := repo.s3Client.PutObject(repo.ctx, &s3.PutObjectInput{
		Bucket: aws.String(repo.bucketName),
		Key:    aws.String("bin/" + fileName),
		Body:   bytes.NewReader(fileData),
	})

	return err
}

func (repo *S3Repository) GetBinary(
	fileName string,
) ([]byte, error) {
	result, err := repo.s3Client.GetObject(repo.ctx, &s3.GetObjectInput{
		Bucket: &repo.bucketName,
		Key:    aws.String("bin/" + fileName),
	})

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = result.Body.Close()
	}()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Println("could not read object body")
		return nil, err
	}

	return body, nil
}
