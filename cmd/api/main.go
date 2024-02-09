package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bvaledev/go-filemanager/configs"
	"github.com/bvaledev/go-filemanager/internal/infra/storage"
	"github.com/bvaledev/go-filemanager/internal/infra/web"
	"github.com/bvaledev/go-filemanager/internal/infra/web/server"
	"github.com/bvaledev/go-filemanager/internal/service"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	fileService := service.NewFileService(buildS3StorageAdapter(configs))
	fileHandler := web.NewFileHandler(fileService)

	httpServer := server.NewWebServer(configs.Port)

	httpServer.AddHandler("/file/upload", "POST", fileHandler.Upload)

	httpServer.Start()
}

func buildS3StorageAdapter(configs *configs.Conf) *storage.S3AdapterImpl {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{URL: configs.S3Endpoint}, nil
	})
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(configs.S3Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(configs.S3Key, configs.S3Secret, "")),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		log.Fatal(err)
	}
	s3Client := s3.NewFromConfig(cfg)
	if _, err = s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{Bucket: &configs.S3Bucket}); err != nil {
		log.Fatal(err)
	}

	return storage.NewS3Adapter(s3Client, configs.S3Bucket)
}
