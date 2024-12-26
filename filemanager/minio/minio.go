package minio

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/starfork/stargo/filemanager"
)

type Minio struct {
	client     *minio.Client
	bucketName string
}

func NewMinio(conf *filemanager.Config) (filemanager.Filemanager, error) {
	client, err := minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure: false, // 是否启用 HTTPS，设置为 true 使用 HTTPS
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}
	// 检查 bucket 是否存在，不存在则创建
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, conf.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, conf.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &Minio{
		client:     client,
		bucketName: conf.BucketName,
	}, nil
}

func (e *Minio) Upload(filePath string) (*filemanager.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 获取文件信息
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	objectName := stat.Name()
	uploadInfo, err := e.client.PutObject(context.Background(), e.bucketName, objectName, file, stat.Size(), minio.PutObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return &filemanager.File{
		Name:        uploadInfo.Key,
		Size:        uploadInfo.Size,
		ContentType: "unknown", // MinIO 不会直接提供 ContentType，这里可通过 MIME 检测扩展
		ModTime:     time.Now(),
		ETag:        uploadInfo.ETag,
		Owner:       "", //minio.Owner.XMLName.Local,
	}, nil
}

func (e *Minio) Download(fileName string) (*filemanager.File, error) {
	object, err := e.client.GetObject(context.Background(), e.bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer object.Close()

	// 将文件保存到本地
	destFile, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create local file: %w", err)
	}
	defer destFile.Close()
	// 获取文件的详细信息
	stat, err := e.client.StatObject(context.Background(), e.bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	size, err := io.Copy(destFile, object)
	if err != nil {
		return nil, fmt.Errorf("failed to copy object to local file: %w", err)
	}

	return &filemanager.File{
		Name:        stat.Key,
		Size:        size,
		ContentType: stat.ContentType,
		ModTime:     stat.LastModified,
		ETag:        stat.ETag,
		Owner:       "defaultOwner",
	}, nil
}

func (e *Minio) Delete(fileName string) error {
	err := e.client.RemoveObject(context.Background(), e.bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

func (e *Minio) List() ([]filemanager.File, error) {
	ctx := context.Background()
	objectCh := e.client.ListObjects(ctx, e.bucketName, minio.ListObjectsOptions{Recursive: true})

	var files []filemanager.File
	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to list files: %w", object.Err)
		}
		files = append(files, filemanager.File{
			Name: object.Key,
			Size: object.Size,
		})
	}
	return files, nil
}

func (e *Minio) Get(fileName string) (*filemanager.File, error) {
	ctx := context.Background()
	stat, err := e.client.StatObject(ctx, e.bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &filemanager.File{
		Name:         stat.Key,
		Size:         stat.Size,
		ContentType:  stat.ContentType,
		ModTime:      stat.LastModified,
		ETag:         stat.ETag,
		StorageClass: stat.StorageClass,
	}, nil
}

func (e *Minio) Rename(oldName, newName string) error {
	// MinIO 不支持直接重命名，通过复制实现
	err := e.Copy(oldName, newName)
	if err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}
	return e.Delete(oldName)
}

func (e *Minio) Copy(sourceFileName, destFileName string) error {
	ctx := context.Background()
	src := minio.CopySrcOptions{
		Bucket: e.bucketName,
		Object: sourceFileName,
	}
	dest := minio.CopyDestOptions{
		Bucket: e.bucketName,
		Object: destFileName,
	}
	_, err := e.client.CopyObject(ctx, dest, src)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	return nil
}

func (e *Minio) Move(sourceFileName, destFileName string) error {
	// 通过复制和删除实现移动
	err := e.Copy(sourceFileName, destFileName)
	if err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}
	return e.Delete(sourceFileName)
}

// func (e *Minio) stat2File(stat minio.ObjectInfo) *filemanager.File {
// 	return &filemanager.File{
// 		Name:         stat.Key,
// 		Size:         stat.Size,
// 		ContentType:  stat.ContentType,
// 		ModTime:      stat.LastModified,
// 		ETag:         stat.ETag,
// 		StorageClass: stat.StorageClass,
// 	}
// }
