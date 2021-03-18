package coord_cfg

import (
	"context"
	"log"
)

// 读取环境变量逻辑
func Get(ctx context.Context, k string, defaultV ...string) (v string) {
	// 读内存缓存
	v = getFromMem(ctx, k)
	if v != "" {
		return v
	}
	// 读redis
	v = getFromRedis(ctx, k)
	if v != "" {
		// 缓存结果到内存
		set2Mem(ctx, k, v)
		return v
	}
	// 读db
	v = getFromDB(ctx, k)
	if v != "" {
		// 缓存结果到redis
		_ = set2Redis(ctx, k, v)
		return v
	}
	// 读yaml
	v = getFromYaml(k)
	if v != "" {
		return v
	}
	// 读环境变量
	v = getFromEnv(k)
	if v != "" {
		return v
	}

	dft := ""
	if len(defaultV) > 0 {
		dft = defaultV[0]
	}
	return dft
}

// 设置配置逻辑
func Set(ctx context.Context, k, v string) error {
	err := set2DB(ctx, k, v)
	if err != nil {
		log.Print("failed set to db ", err)
		return err
	}
	go func() {
		err := cleanRedis(ctx, k)
		if err != nil {
			log.Print("failed clean redis ", k, err)
		}
	}()
	set2Mem(ctx, k, "")
	return nil
}
