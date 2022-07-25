package util

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mime/multipart"
	"time"
	"vote/vote/app/config"
)

func UploadFile(img *multipart.FileHeader, userID, fileType string) (*string, error) {

	name := fmt.Sprintf("%v_%s_%s", TimeToFormat(time.Now()), userID, img.Filename)

	//创建ossClient实例
	client, err := oss.New(config.AliOSS.Endpoint, config.AliOSS.AccessKeyId, config.AliOSS.AccessKeySecret)
	if err != nil {
		return nil, errors.New("创建ossClient实例失败：" + err.Error())
	}

	//获取存储空间
	bucket, err := client.Bucket(config.AliOSS.BucketName)
	if err != nil {
		return nil, errors.New("获取存储空间失败：" + err.Error())
	}

	file, err1 := img.Open()
	if err1 != nil {
		return nil, err1
	}
	defer file.Close()

	//上传
	options := []oss.Option{
		oss.ContentType(fileType),
	}
	err3 := bucket.PutObject(name, file, options...)
	if err3 != nil {
		return nil, errors.New("上传文件失败：" + err3.Error())
	}

	url := "http://" + config.AliOSS.BucketName + "." + config.AliOSS.Endpoint + "/" + name

	return &url, nil
}
