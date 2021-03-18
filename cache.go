package coord_cfg

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	mem "github.com/patrickmn/go-cache"
	"strconv"
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
func connRedis(ctx context.Context, host, port, pwd, db string) error {
	i, _ := strconv.Atoi(db)
	redisCli = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: pwd,
		DB:       i,
	})
	return redisCli.Ping(ctx).Err()
}

func getFromRedis(ctx context.Context, k string) string {
	if redisCli == nil { // 如果redis未进行初始化,则不做操作
		return ""
	}
	res := redisCli.Get(ctx, k)
	if res.Err() != nil {
		return ""
	}
	return res.String()
}

func set2Redis(ctx context.Context, k, v string) error {
	if redisCli == nil { // 如果redis未进行初始化,则不做操作
		return fmt.Errorf("not init redis ")
	}
	return redisCli.Set(ctx, k, v, DefaultCacheTime).Err()
}

func cleanRedis(ctx context.Context, k string) error {
	if redisCli == nil { // 如果redis未进行初始化,则不做操作
		return fmt.Errorf("not init redis ")
	}
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
