package oss

import (
	"encoding/base64"
	"fmt"
	aliyunOss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"rashomon/conf"
	"rashomon/helpers"
	"strings"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

// GetBucket 获取bucket
func (oss *Client) GetBucket() (*aliyunOss.Bucket, error) {
	ossClient, err := aliyunOss.New(conf.AppConfig.Oss.EndPoint, conf.AppConfig.Oss.AccessKeyId, conf.AppConfig.Oss.AccessKeySecret, aliyunOss.SecurityToken(conf.AppConfig.Oss.SecurityToken))
	if err != nil {
		return nil, err
	}
	return ossClient.Bucket(conf.AppConfig.Oss.Bucket)
}

// UploadImage 上传图片接口.  uploadPath: 上传路径
func (oss *Client) UploadImage(uploadPath string, imageBase64 string) (string, error) {
	bucket, err := oss.GetBucket()
	if err != nil {
		return "", err
	}
	imageReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imageBase64))
	if err = bucket.PutObject(uploadPath, imageReader); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", conf.AppConfig.Oss.Address, uploadPath), nil
}

// UploadFile 上传文件接口.  localPath: 文件的本地路径  uploadPath: 上传路径
func (oss *Client) UploadFile(localPath string, uploadPath string) (string, error) {
	bucket, err := oss.GetBucket()
	if err != nil {
		return "", err
	}
	fd, err := os.Open(localPath)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	if err = bucket.PutObject(uploadPath, fd); err != nil {
		return "", err
	}
	return helpers.ResolveUrl(fmt.Sprintf("%s/%s", conf.AppConfig.Oss.Address, uploadPath)), nil
}
