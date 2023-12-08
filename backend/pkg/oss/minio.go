package oss

import (
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	ioconfig "ims-server/pkg/config"
	"mime/multipart"
	"net/url"
	"time"
)

type MinioClient struct {
	client *minio.Client
}

func NewMinioClient() (*MinioClient, error) {
	config := ioconfig.GetMinioConf()
	endpoint := config.Endpoint
	accessKeyID := config.AccessKey
	secretAccessKey := config.SecretKey
	useSSL := config.UseSLL

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL},
	)
	if err != nil {
		return nil, err
	}

	return &MinioClient{
		client: minioClient,
	}, nil
}

func (mc *MinioClient) MakeBucket(name string) error {
	exists, err := mc.client.BucketExists(context.Background(), name)
	if err == nil && exists {
		return errors.New("bucket already exists")
	}
	err = mc.client.MakeBucket(context.Background(), name, minio.MakeBucketOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (mc *MinioClient) DeleteBucket(name string) error {
	exists, err := mc.client.BucketExists(context.Background(), name)
	if err != nil || !exists {
		return errors.New("bucket not exists")
	}
	err = mc.client.RemoveBucket(context.Background(), name)
	if err != nil {
		return err
	}
	return nil
}

// PutObject 上传小文件
func (mc *MinioClient) PutObject(ctx context.Context, bucketName string, name string, fileHeader *multipart.FileHeader) error {
	size := fileHeader.Size
	if size > 8*1024*1024*1024 {
		// 大于 1G，不让传
		return errors.New("file too large")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	_, err = mc.client.PutObject(ctx, bucketName, name, file, fileHeader.Size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return err
}

// GetObject 下载文件
func (mc *MinioClient) GetObject(ctx context.Context, bucketName string, name string) ([]byte, error) {
	object, err := mc.client.GetObject(ctx, bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()
	return nil, err
}

// PresignedGetObject 生成获取文件的预签名 url，不清楚需要什么样的请求头，可将 reqParams 传入 nil（注意：即使桶是私有的仍可获取）
func (mc *MinioClient) PresignedGetObject(ctx context.Context, bucketName, objectName string, expiry time.Duration, reqParams url.Values) (*url.URL, error) {
	// Generates a presigned url which expires in a day.
	url, err := mc.client.PresignedGetObject(ctx, bucketName, objectName, expiry, reqParams)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return url, nil
}

func (mc *MinioClient) ListObjects(ctx context.Context, bucketName string, prefix string) []minio.ObjectInfo {
	var objects []minio.ObjectInfo
	opt := minio.ListObjectsOptions{}
	if prefix != "" {
		// 有传入前缀，按前缀查找
		opt.Prefix = prefix
		opt.Recursive = true
	}
	for object := range mc.client.ListObjects(ctx, bucketName, opt) {
		objects = append(objects, object)
	}
	return objects
}

func (mc *MinioClient) CopyObject(ctx context.Context, originalBucket string, originalName string, newBucket string, newName string) error {
	// destination
	dest := minio.CopyDestOptions{
		Bucket: newBucket,
		Object: newName,
	}
	// source object
	src := minio.CopySrcOptions{
		Bucket: originalBucket,
		Object: originalName,
	}

	_, err := mc.client.CopyObject(ctx, dest, src)
	return err
}

func (mc *MinioClient) DeleteObject(ctx context.Context, bucketName string, name string) error {
	return mc.client.RemoveObject(ctx, bucketName, name, minio.RemoveObjectOptions{})
}
