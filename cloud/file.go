package cloud

import (
	"context"
	"fmt"
	"io"
	"kawa_blog/client"
	"kawa_blog/utils"
	"log"
	"mime"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func UploadFile(objectKey string, r io.Reader) (string, error) {
	var bucketName string = utils.GetEnv("R2_BUCKET_NAME")
	var publicURL string = utils.GetEnv("R2_PUBLIC_URL")
	var fileName string = uuid.NewString()

	var objectKeyParts []string = strings.Split(objectKey, ".")
	var ext string = "." + objectKeyParts[len(objectKeyParts)-1]
	var contentType string = mime.TypeByExtension(ext)

	s3Client, err := client.GetS3Client()
	if err != nil {
		log.Fatalf("S3 client is not initialized: %v", err)
		return "", err
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileName),
		Body:        r,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Cloudflare R2 の公開 URL を生成
	imageURL := fmt.Sprintf("%s/%s", publicURL, fileName)
	return imageURL, nil
}
