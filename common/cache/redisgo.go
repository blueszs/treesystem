package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacheDb struct {
	Client  redis.Client
	Context context.Context
	Flag    bool
}

// InitRedisCache 初始化Redis连接，返回Redis客户端对象
func InitRedisCache(addr, pwd string, dbindex int) *RedisCacheDb {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,     // 没有密码，默认值
		DB:       dbindex, // 默认DB 0
	})
	ctx := context.Background()
	return &RedisCacheDb{Client: *rdb, Context: ctx, Flag: true}
}

// Set 设置Redis缓存值
// @key 缓存键
// @value 缓存值
// @expire 过期时间S，值为0表示一直有效。
func (r *RedisCacheDb) Set(key string, value interface{}, expire int) error {
	// 设置过期时间
	var duration time.Duration = time.Duration(expire) * time.Second
	err := r.Client.Set(r.Context, key, value, duration).Err()
	return err
}

// Get 获取Redis中缓存的值
// @key 缓存键
func (r *RedisCacheDb) Get(key string) string {
	val, err := r.Client.Get(r.Context, key).Result()
	if err != nil {
		fmt.Printf("not find key【%s】\n", key)
		return ""
	}
	return val
}

// Expire 设置Redis键的过期时间
// @key 缓存键
// @expire 过期时间S，值为0表示一直有效。
func (r *RedisCacheDb) Expire(key string, expire int) {
	// 设置过期时间
	var duration time.Duration = time.Duration(expire) * time.Second
	err := r.Client.Expire(r.Context, key, duration).Err()
	if err != nil {
		fmt.Printf("expire fail:%v\n", err)
		return
	}
}

// Del 删除Redis缓存中的键
// @key 缓存键
func (r *RedisCacheDb) Del(key string) (bool, error) {
	err := r.Client.Del(r.Context, key).Err()
	if err != nil {
		fmt.Printf("redis delelte failed:%v\n", err)
		return false, err
	}
	return true, nil
}

// Exists 判断键是否存在于Redis中
// @key 缓存键
func (r *RedisCacheDb) Exists(key string) (bool, error) {
	isExit := r.Client.Exists(r.Context, key)
	err := isExit.Err()
	if err != nil {
		fmt.Printf("redis exists failed:%v\n", err)
		return false, err
	}
	return isExit.Val() > 0, nil
}

// Close 关闭Redis客户端对象的连接
func (r *RedisCacheDb) Close() error {
	if r != nil {
		r.Context.Deadline()
		return r.Client.Close()
	} else {
		err := errors.New("redis client object is null or undefined")
		if err != nil {
			return err
		}
		return nil
	}
}
