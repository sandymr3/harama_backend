package storage

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinioStorage{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (s *MinioStorage) UploadFile(ctx context.Context, objectName string, data []byte, contentType string) (string, error) {
	reader := bytes.NewReader(data)
	_, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	// Generate a temporary URL for access (or return permanent one if public)
	expiry := time.Duration(168) * time.Hour // 1 week
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucketName, objectName, expiry, url.Values{})
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (s *MinioStorage) GetFile(ctx context.Context, objectName string) ([]byte, error) {
	object, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	return io.ReadAll(object)
}
