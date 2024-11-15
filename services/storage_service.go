package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/mich31/scoreplay-media-api/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageService struct {
	Client     *minio.Client
	BucketName string
}

func NewStorageService() (*StorageService, error) {
	fmt.Println("initializing storage service..")
	client, err := minio.New(config.Config("STORAGE_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(config.Config("STORAGE_ACCESS_KEY_ID"), config.Config("STORAGE_SECRET_ACCESS_KEY"), ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage client: %w", err)
	}

	fmt.Println("storage service initialized")

	return &StorageService{
		Client: client,
	}, nil
}

func (service *StorageService) CreateBucket(ctx context.Context, bucketName string) error {
	err := service.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: config.Config("STORAGE_BUCKET_REGION")})
	service.BucketName = bucketName
	if err != nil {
		exists, errBucketExists := service.Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("bucket %s already exists\n", bucketName)
			return nil
		} else {
			return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
		}
	}

	log.Printf("bucket %s created!\n", bucketName)
	policy := fmt.Sprintf(`{"Version": "2012-10-17","Statement": [{"Action": ["s3:*"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::%s/*"]}]}`, bucketName)
	err = service.Client.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		log.Printf("unable to set policy for bucket %s: %w", bucketName, err)
	}
	return nil
}

func (service *StorageService) UploadObject(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	fileExtension := filepath.Ext(fileHeader.Filename)
	objectName := fmt.Sprintf("%s%s", uuid.New(), fileExtension)
	file, err := fileHeader.Open()
	if err != nil {
		log.Fatalln("unable to open file: ", err)
		return "", err
	}
	defer file.Close()

	_, err = service.Client.PutObject(ctx, service.BucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s/%s/%s", config.Config("STORAGE_ENDPOINT"), config.Config("STORAGE_BUCKET_NAME"), objectName), nil
}
