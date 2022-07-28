package service

import (
	"pilipili/model"
	"pilipili/serializer"
)

// ListVideoService 视频列表服务
type ListVideoService struct {
	Limit int `form:"limit"` // 输出个数
	Start int `form:"start"` // 数据起始位置
}

const (
	_DEFUALT_GET_LIST_LIMIT = 6 // 默认的获取条数
)

// List 视频列表
func (service *ListVideoService) List() serializer.Response {
	var videos []model.Video
	total := 0

	// 设置数据个数默认值
	if service.Limit == 0 {
		service.Limit = _DEFUALT_GET_LIST_LIMIT
	}

	// 符合条件 数据总数
	if err := model.DB.Model(model.Video{}).Count(&total).Error; err != nil {
		return serializer.Response{
			Status: 50000,
			Msg:    "数据库连接错误",
			Error:  err.Error(),
		}
	}

	// Find 查找多个数据
	// 这里处理分页问题
	if err := model.DB.Limit(service.Limit).Offset(service.Start).Find(&videos).Error; err != nil {
		return serializer.Response{
			Status: 50000,
			Msg:    "数据库连接错误",
			Error:  err.Error(),
		}
	}

	// 这里重新封装 JSON 数据，带上 数据内容 and 数据个数
	return serializer.BuildListResponse(
		serializer.BuildVideos(videos), // 数据
		uint(total),
	)
}
