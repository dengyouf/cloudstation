package aliyun

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/dengyouf/cloudstation/store"
	"github.com/go-playground/validator/v10"
)

// 迁移参数校验逻辑，第三方库
var (
	validate = validator.New()
)

func (a *aliyun) validate() error {
	return validate.Struct(a)
}

type aliyun struct {
	Endpoint  string `validate:"required"`
	AcessKey  string `validate:"required"`
	SecretKey string `validate:"required"`

	listener oss.ProgressListener
}

func NewUploader(Endpoint, AcessKey, SecretKey string) (store.Uploader, error) {
	uploader := &aliyun{
		Endpoint:  Endpoint,
		AcessKey:  AcessKey,
		SecretKey: SecretKey,
		listener:  NewListener(),
	}
	if err := uploader.validate(); err != nil {
		return nil, err
	}
	return uploader, nil
}

func (a *aliyun) UploadFile(bucketName, objectKey, localFilePath string) error {
	if bucketName == "" || objectKey == "" || localFilePath == "" {
		return fmt.Errorf("bucketName or objectKey or localFilePath is missed")
	}

	client, err := oss.New(a.Endpoint, a.AcessKey, a.SecretKey)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 第一个参数是云端路径，第二个是本地路径
	err = bucket.PutObjectFromFile(objectKey, localFilePath, oss.Progress(a.listener))
	if err != nil {
		return err
	}

	// 打印下载URL
	signedURL, err := bucket.SignURL(localFilePath, oss.HTTPGet, 60*60*24)
	if err != nil {
		return fmt.Errorf("SignURL error, %s", err)
	}
	fmt.Printf("下载链接: %s\n", signedURL)
	fmt.Println("\n注意: 文件下载有效期为1天, 中转站保存时间为3天, 请及时下载")

	// return fmt.Errorf("no impl")
	return nil
}
