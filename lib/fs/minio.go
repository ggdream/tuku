package fs

import (
	"context"
	"io"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIO 对象存储
// 存储格式：gallery
//			 - raw
//			 - gen
//				- ${hash}.${width}_${height}.${suffix}
type MinIO struct {
	client  *minio.Client
	context context.Context
	bucket  string
}

func NewMinIO(endpoint, bucket, accessKey, secretKey string, useSSL bool) (*MinIO, error) {
	option := &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	}
	client, err := minio.New(endpoint, option)
	if err != nil {
		return nil, err
	}

	isExists, err := client.BucketExists(context.Background(), bucket)
	if err != nil {
		return nil, err
	}
	if !isExists {
		opt := minio.MakeBucketOptions{
			ObjectLocking: true,
		}
		err = client.MakeBucket(context.Background(), bucket, opt)
		if err != nil {
			return nil, err
		}
	}

	return &MinIO{
		client:  client,
		context: context.Background(),
		bucket:  bucket,
	}, nil
}

// WriteFile 保存到MinIO
func (m *MinIO) WriteFile(reader io.Reader, objectName string, fileSize int64, isRaw bool) (*StorageResult, error) {
	ctx, cancel := context.WithTimeout(m.context, 3*time.Minute)
	defer cancel()

	opt := minio.PutObjectOptions{}
	_, err := m.client.PutObject(ctx, m.bucket, m.getPath(objectName, isRaw), reader, fileSize, opt)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ReadFile 获取数据的字节切片
func (m *MinIO) ReadFile(objectName string, isRaw bool) ([]byte, error) {
	ctx, cancel := context.WithTimeout(m.context, 3*time.Minute)
	defer cancel()

	opt := minio.GetObjectOptions{}
	object, err := m.client.GetObject(ctx, m.bucket, m.getPath(objectName, isRaw), opt)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(object)
}

// GetReader 获取Reader，用于响应流
func (m *MinIO) GetReader(objectName string, isRaw bool) (io.Reader, int64, string, error) {
	ctx, cancel := context.WithTimeout(m.context, 3*time.Minute)
	defer cancel()

	opt := minio.GetObjectOptions{}
	object, err := m.client.GetObject(ctx, m.bucket, objectName, opt)
	if err != nil {
		return nil, 0, "", err
	}
	info, err := object.Stat()
	if err != nil {
		return nil, 0, "", err
	}
	return object, info.Size, info.ContentType, nil
}

func (m *MinIO) getPath(objectName string, isRaw bool) string {
	adName := m.adjustName(objectName)
	if isRaw {
		return m.getRawPath(adName)
	} else {
		return m.getGenPath(adName)
	}
}

func (m *MinIO) adjustName(objectName string) string {
	return objectName
}

func (m *MinIO) getRawPath(objectName string) string {
	return filepath.Join(RAW, objectName)
}

func (m *MinIO) getGenPath(objectName string) string {
	return filepath.Join(GEN, objectName)
}
