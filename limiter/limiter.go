package limiter

import (
	"context"
	"sort"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter interface {
	Wait(ctx context.Context) error
	Limit() rate.Limit
}

// Per limiter.Per(20, 1*time.Minute) 代表速率是每 1 分钟补充 20 个令牌。
func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

func NewMultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters, byLimit)
	return &multiLimiter{limiters: limiters}
}

type multiLimiter struct {
	// 聚合多个 RateLimiter，并将速率由小到大排序
	limiters []RateLimiter
}

// Wait  方法会循环遍历多层限速器 multiLimiter 中所有的限速器并索要令牌，只有当所有的限速器规则都满足后，才会正常执行后续的操作。
func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}
