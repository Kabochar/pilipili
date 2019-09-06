package service

import (
	"fmt"
	"strings"

	"pilipili/cache"
	"pilipili/model"
	"pilipili/serializer"
)

// DailyRankService 每日排行的服务
type DailyRankService struct {
}

// Get 获取排行
func (service *DailyRankService) Get() serializer.Response {
	var videos []model.Video

	// 从redis读取点击前十的视频
	vids, _ := cache.RedisClient.ZRevRange(cache.DailyRankKey, 0, 5).Result()

	// 数据大于 1 才进行操作
	if len(vids) > 1 {
		// order，拼接 redis 结果
		order := fmt.Sprintf("FIELD(id, %s)", strings.Join(vids, ","))
		// 使用 ORDER BY 操作，找出所有符合条件的结果
		err := model.DB.Where("id in (?)", vids).Order(order).Find(&videos).Error
		if err != nil {
			return serializer.Response{
				Status: 50000,
				Msg:    "数据库连接错误",
				Error:  err.Error(),
			}
		}
	}

	return serializer.Response{
		Data: serializer.BuildVideos(videos),
	}
}
