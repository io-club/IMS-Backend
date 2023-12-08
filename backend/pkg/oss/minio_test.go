package oss

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"testing"
	"time"
)

func TestMinio(t *testing.T) {
	ctx := context.Background()

	endpoint := "101.200.63.44:49000"
	accessKeyID := "admin"
	secretAccessKey := "admin123"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL},
	)
	if err != nil || minioClient == nil {
		log.Printf("create minio client failed ,error: %s\n", err)
		return
	}
	client := &MinioClient{
		client: minioClient,
	}
	log.Printf("client create success: %+v", *client)

	err = client.MakeBucket("ims-test")
	if err != nil {
		log.Printf("make bucket failed ,error: %s\n", err)
		return
	}
	log.Printf("make bucket success\n")
	time.Sleep(75 * time.Second)
	// 暂时没法测，应该予以补测
	//err = client.PutObject(ctx, "ims-test", "ims-test", fileHeader)
	//if err != nil {
	//	log.Printf("put object failed ,error: %s\n", err)
	//	return
	//}
	//log.Printf("put object success\n")

	obj := client.ListObjects(ctx, "ims-test", "")
	for _, object := range obj {
		log.Printf("object: %+v\n", object)
		getObject, err := client.PresignedGetObject(ctx, "ims-test", object.Key, 10*time.Minute, nil)
		if err != nil {
			log.Printf("get object failed ,error: %s\n", err)
			return
		}
		log.Printf("getObject: %s\n", getObject)
	}

	time.Sleep(60 * time.Second)

	err = client.DeleteObject(ctx, "ims-test", "ims-test-object.txt")
	if err != nil {
		log.Printf("delete object failed ,error: %s\n", err)
		return
	}
	log.Printf("delete object success\n")
	time.Sleep(30 * time.Second)
	err = client.DeleteBucket("ims-test")
	if err != nil {
		log.Printf("delete bucket failed ,error: %s\n", err)
		return
	}
	log.Printf("delete bucket success\n")
}
