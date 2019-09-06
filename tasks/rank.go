package tasks

import "pilipili/cache"

// RestartDailyRank 重启一天的排名
func ResetDailyRank() error {
	return cache.RedisClient.Del("rank:daily").Err()
}