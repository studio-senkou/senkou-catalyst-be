package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"senkou-catalyst-be/utils/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Config struct {
	Host            string `json:"host"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

func NewS3Config(host, region, bucket, accessKeyID, secretAccessKey string) *S3Config {
	return &S3Config{
		Host:            host,
		Region:          region,
		Bucket:          bucket,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
}

func NewStorage(config *S3Config) *aws.Config {
	return &aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(
			config.AccessKeyID,
			config.SecretAccessKey,
			"",
		),
		Region:       config.Region,
		BaseEndpoint: aws.String(config.Host),
	}
}

func GetStorage() *aws.Config {
	config := NewS3Config(
		config.GetEnv("AWS_S3_HOST", "http://localhost:9989"),
		config.GetEnv("AWS_S3_REGION", "us-east-1"),
		config.GetEnv("AWS_S3_BUCKET", "catadev"),
		config.GetEnv("AWS_S3_ACCESS_KEY_ID", ""),
		config.GetEnv("AWS_S3_SECRET_ACCESS_KEY", ""),
	)

	return NewStorage(config)
}

func GetBucket() string {
	return config.GetEnv("AWS_S3_BUCKET", "catadev")
}

type UploadService struct {
	storage *s3.Client
	bucket  string
}

func NewUploadService() *UploadService {
	config := GetStorage()
	bucket := GetBucket()

	return &UploadService{
		storage: s3.NewFromConfig(*config, func(o *s3.Options) {
			o.UsePathStyle = true
		}),
		bucket: bucket,
	}
}

func (s *UploadService) UploadFile(ctx context.Context, file *multipart.FileHeader, filename string, folder *string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	var key string
	if folder != nil && *folder != "" {
		key = fmt.Sprintf("%s/%s", *folder, filename)
	} else {
		key = filename
	}

	_, err = s.storage.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          src,
		ContentType:   aws.String(file.Header.Get("Content-Type")),
		ContentLength: aws.Int64(file.Size),
		ACL:           types.ObjectCannedACLPublicRead,
		CacheControl:  aws.String("max-age=31536000, public"),
		Metadata: map[string]string{
			"filename":    file.Filename,
			"uploaded-by": "Lentera Cendekia API",
		},
	})
	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *UploadService) GetFileURL(ctx context.Context, path string) (string, error) {
	presignClient := s3.NewPresignClient(s.storage)
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	}

	presignedReq, err := presignClient.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = 30 * 60 // Persist for 30 minutes
	})

	if err != nil {
		return "", fmt.Errorf("failed to presign URL: %w", err)
	}

	return presignedReq.URL, nil
}

func (s *UploadService) RemoveFile(ctx context.Context, path string) error {
	_, err := s.storage.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}
	return nil
}
