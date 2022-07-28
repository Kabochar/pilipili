package service

import (
	"pilipili/model"
	"pilipili/serializer"

	"github.com/gin-gonic/gin"
)

// DeleteVideoService 删除投稿的服务
type DeleteVideoService struct{}

// Delete 删除视频
func (service *DeleteVideoService) Delete(c *gin.Context) serializer.Response {
	var video model.Video
	err := model.DB.First(&video, c.Param("id")).Error
	if err != nil {
		return serializer.Response{
			Status: 404,
			Msg:    "视频不存在",
			Error:  err.Error(),
		}
	}

	err = model.DB.Delete(&video).Error
	if err != nil {
		return serializer.Response{
			Status: 50000,
			Msg:    "视频删除失败",
			Error:  err.Error(),
		}
	}

	return serializer.Response{}
}
