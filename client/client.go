package client

import (
	"context"
	"fmt"
	"sync"

	"kawa_blog/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3クライアントのシングルトン
var (
	s3Client *s3.Client
	once     sync.Once
)

// S3クライアントを初期化する関数
func InitS3Client() error {
	var err error
	once.Do(func() { // 一度だけ実行
		s3Client, err = CreateS3Client()
	})
	return err
}

// S3クライアントの取得関数
func GetS3Client() (*s3.Client, error) {
	if s3Client == nil {
		return nil, fmt.Errorf("S3 client is not initialized. Call InitS3Client() first")
	}
	return s3Client, nil
}

func CreateS3Client() (*s3.Client, error) {
	accessKey := utils.GetEnv("R2_ACCESS_KEY")
	secretKey := utils.GetEnv("R2_SECRET_KEY")
	endpoint := utils.GetEnv("R2_ENDPOINT")

	if accessKey == "" || secretKey == "" || endpoint == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"), // Cloudflare R2 の場合は "auto"
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint}, nil
			},
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return s3.NewFromConfig(cfg), nil
}
