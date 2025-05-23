package config

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gopkg.in/ini.v1"
	"strings"
)

// MinioConfig 存储 MinIO 相关配置
type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Secure    bool
}

// 全局 MinIO 客户端和配置
var MinioSettings MinioConfig

var MinioClient *minio.Client

// InitMinio 初始化 MinIO 客户端
func InitMinio() error {
	cfg, err := ini.Load("./config/app.ini")
	if err != nil {
		return fmt.Errorf("加载 MinIO 配置失败: %v", err)
	}

	section := cfg.Section("minio")
	MinioSettings.Endpoint = section.Key("endpoint").String()
	MinioSettings.AccessKey = section.Key("accessKey").String()
	MinioSettings.SecretKey = section.Key("secretKey").String()
	MinioSettings.Bucket = section.Key("bucket").String()
	MinioSettings.Secure = strings.ToLower(section.Key("secure").String()) == "true"

	// 创建客户端
	MinioClient, err = minio.New(MinioSettings.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(MinioSettings.AccessKey, MinioSettings.SecretKey, ""),
		Secure: MinioSettings.Secure,
	})
	if err != nil {
		return fmt.Errorf("MinIO 初始化失败: %v", err)
	}
	fmt.Println("Minio 连接成功!")

	// 检查 Bucket 是否存在，不存在则创建
	ctx := context.Background()
	exists, err := MinioClient.BucketExists(ctx, MinioSettings.Bucket)
	if err != nil {
		return fmt.Errorf("检查桶失败: %v", err)
	}
	if !exists {
		err = MinioClient.MakeBucket(ctx, MinioSettings.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("创建桶失败: %v", err)
		}
		fmt.Println("已创建桶:", MinioSettings.Bucket)
	} else {
		fmt.Println("桶已存在:", MinioSettings.Bucket)
	}

	return nil
}
