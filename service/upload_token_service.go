package service

import (
	"mime"
	"os"
	"path/filepath"

	"pilipili/serializer"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

// UploadTokenService 获得上传oss token的服务
type UploadTokenService struct {
	Filename string `form:"filename" json:"filename"`
}

// Post 创建 token
func (service *UploadTokenService) Post() serializer.Response {
	client, err := oss.New(
		os.Getenv("OSS_END_POINT"),
		os.Getenv("OSS_ACCESS_KEY_ID"),
		os.Getenv("OSS_ACCESS_KEY_SECRET"))
	if err != nil {
		return serializer.Response{
			Status: 50002,
			Msg:    "OSS配置错误",
			Error:  err.Error(),
		}
	}

	// 获取存储空间。
	bucket, err := client.Bucket(os.Getenv("OSS_BUCKET"))
	if err != nil {
		return serializer.Response{
			Status: 50002,
			Msg:    "OSS配置错误",
			Error:  err.Error(),
		}
	}

	// 获取扩展名
	ext := filepath.Ext(service.Filename)

	// 带可选参数的签名直传，限制文件类型，避免恶意上传
	options := []oss.Option{
		oss.ContentType(mime.TypeByExtension(ext)),
	}

	// 路径限制，方便管理
	key := "upload/avatar/" + uuid.Must(uuid.NewRandom()).String() + ext
	// 签名直传
	// 600： 有效时间
	signedPutURL, err := bucket.SignURL(key, oss.HTTPPut, 600, options...)
	if err != nil {
		return serializer.Response{
			Status: 50002,
			Msg:    "OSS配置错误",
			Error:  err.Error(),
		}
	}

	// 上传完成后，再 GET 请求 查看图片
	signedGetURL, err := bucket.SignURL(key, oss.HTTPGet, 600)
	if err != nil {
		return serializer.Response{
			Status: 50002,
			Msg:    "OSS配置错误",
			Error:  err.Error(),
		}
	}

	// 数据已经打包完成，其他信息发送给前端处理
	return serializer.Response{
		Data: map[string]string{
			"key": key,
			"put": signedPutURL,
			"get": signedGetURL,
		},
	}
}
