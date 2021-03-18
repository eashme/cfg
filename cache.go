package coord_cfg

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	mem "github.com/patrickmn/go-cache"
	"sync"
)

var (
	// 存 redis
	redisCli *redis.Client
	// 存内存
	memCache     *mem.Cache
	memCacheOnce sync.Once
)

// redis
func connRedis(ctx context.Context, host string, port uint16, pwd string, db int) error {
	redisCli = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: pwd,
		DB:       db,
	})
	return redisCli.Ping(ctx).Err()
}

func getFromRedis(ctx context.Context, k string) string {
	return redisCli.Get(ctx, k).String()
}

func set2Redis(ctx context.Context, k, v string) error {
	return redisCli.Set(ctx, k, v, DefaultCacheTime).Err()
}

func cleanRedis(ctx context.Context, k string) error{
	return redisCli.Del(ctx, k).Err()
}

// 内存
func getFromMem(ctx context.Context, k string) string {
	memCacheOnce.Do(func() {
		memCache = mem.New(DefaultCacheTime, ScanInterval)
	})
	v, ok := memCache.Get(k)
	if ok {
		return v.(string)
	}
	return ""
}

func set2Mem(ctx context.Context, k, v string) {
	memCacheOnce.Do(func() {
		memCache = mem.New(DefaultCacheTime, ScanInterval)
	})
	memCache.SetDefault(k, v)
}
