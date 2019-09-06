package model

import (
	"pilipili/cache"
	"strconv"

	"github.com/jinzhu/gorm"
)

// Video 视频模型
type Video struct {
	gorm.Model
	Title  string
	Info   string
	URL    string
	Avatar string
}

// View 点击数
// JSON 需要的数值是 uint 类型，redis 的结果数据是字符串类型
func (video *Video) View() uint64 {
	countStr, _ := cache.RedisClient.Get(cache.VideoViewKey(video.ID)).Result()
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

// AddView 视频游览
func (video *Video) AddView() {
	// 增加视频点击数
	cache.RedisClient.Incr(cache.VideoViewKey(video.ID))
	// 增加排行点击数
	cache.RedisClient.ZIncrBy(cache.DailyRankKey, 1, strconv.Itoa(int(video.ID)))
}
