package aliyun_test

import (
	"testing"

	"github.com/dengyouf/cloudstation/store/provider/aliyun"
	"github.com/stretchr/testify/assert"
)

var (
	Endpoint  = "http://oss-cn-beijing.aliyuncs.com"
	AcessKey  = ""
	SecretKey = ""
)
var (
	bucketName    = "cloudstations"
	objectKey     = "new_store.go"
	localFilePath = "store.go"
)

func TestUploadFile(t *testing.T) {
	should := assert.New(t)

	uploader, err := aliyun.NewUploader(Endpoint, AcessKey, SecretKey)

	if should.NoError(err) {
		err = uploader.UploadFile(bucketName, objectKey, localFilePath)
		should.NoError(err)
	}

}
