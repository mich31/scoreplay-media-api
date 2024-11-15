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
	fmt.Println("Initializing storage service..")
	client, err := minio.New(config.Config("STORAGE_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(config.Config("STORAGE_ACCESS_KEY_ID"), config.Config("STORAGE_SECRET_ACCESS_KEY"), ""),
		Secure: false,
	})
	if err != nil {
		return nil, err // TODO
	}

	fmt.Println("Storage service initialized")

	return &StorageService{
		Client: client,
	}, nil
}

func (service *StorageService) CreateBucket(ctx context.Context, bucketName string) {
	err := service.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: config.Config("STORAGE_BUCKET_REGION")})
	service.BucketName = bucketName // TODO
	if err != nil {
		exists, errBucketExists := service.Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Bucket %s created!\n", bucketName)
	}
}

func (service *StorageService) UploadObject(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	fileExtension := filepath.Ext(fileHeader.Filename)
	objectName := fmt.Sprintf("%s%s", uuid.New(), fileExtension)
	file, err := fileHeader.Open()
	if err != nil {
		log.Fatalln("Unable to open file: ", err)
		return "", err
	}
	defer file.Close()

	_, err = service.Client.PutObject(ctx, service.BucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s/%s/%s", config.Config("STORAGE_ENDPOINT"), config.Config("STORAGE_BUCKET_NAME"), objectName), nil
}
