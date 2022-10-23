package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fengde/gocommon/storex/redisx"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// 基于IP对访问频率限制
func IPRatelimit(redisClient *redisx.Client, limit int, slidingWindow time.Duration) gin.HandlerFunc {
	if err := redisClient.PingCheck(); err != nil {
		panic(err.Error())
	}

	return func(c *gin.Context) {
		now := time.Now().UnixNano()
		userCntKey := fmt.Sprint(c.ClientIP(), ":ip_ratelimit")

		redisClient.GetClient().ZRemRangeByScore(context.Background(), userCntKey,
			"0",
			fmt.Sprint(now-(slidingWindow.Nanoseconds()))).Result()

		reqs, _ := redisClient.GetClient().ZRange(context.Background(), userCntKey, 0, -1).Result()

		if len(reqs) >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  "failed",
				"message": "too many request",
				"data":    map[string]interface{}{},
			})
			return
		}

		c.Next()
		redisClient.GetClient().ZAddNX(context.Background(), userCntKey, &redis.Z{Score: float64(now), Member: float64(now)})
		redisClient.GetClient().Expire(context.Background(), userCntKey, slidingWindow)
	}
}
