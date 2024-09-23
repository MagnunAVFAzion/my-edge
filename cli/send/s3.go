package send

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func uploadFile(ctx context.Context, client *s3.Client, bucketName string, filePath string, key string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filePath, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info %q, %v", filePath, err)
	}
	contentLength := fileInfo.Size()

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: &contentLength,
		ContentType:   aws.String("application/octet-stream"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file %q, %v", filePath, err)
	}

	fmt.Printf("Successfully uploaded %s to %s\n", filePath, key)
	return nil
}

func uploadDir(ctx context.Context, client *s3.Client, bucketName string, dirPath string, baseKey string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relativePath := strings.TrimPrefix(path, dirPath)
		s3Key := filepath.Join(baseKey, relativePath)

		err = uploadFile(ctx, client, bucketName, path, s3Key)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to upload directory %q, %v", dirPath, err)
	}

	return nil
}

func UploadDirFiles(dirPath, customEndpoint, bucketName, projectId string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-1"),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				if service == s3.ServiceID {
					return aws.Endpoint{
						URL:               customEndpoint,
						HostnameImmutable: true,
					}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested for service: %s", service)
			}),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to load aws config, %w", err)
	}

	client := s3.NewFromConfig(cfg)

	err = uploadDir(context.TODO(), client, bucketName, dirPath, projectId)
	if err != nil {
		return err
	}

	return nil
}
