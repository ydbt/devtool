package redisclient

import (
	"context"
	"errors"
	"time"

	"github.com/ydbt/devtool/v3/logger"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(cfg *RedisCfg, lg logger.LogI) *RedisClient {
	rc := new(RedisClient)
	rc.client = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          cfg.Addrs,
		RouteByLatency: false,
		RouteRandomly:  false,
		Username:       cfg.User,
		Password:       cfg.Password,
	})
	rc.ctx = context.Background()
	rc.log = lg
	err := rc.client.ForEachShard(rc.ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		panic(err)
	}
	st, err := rc.Ping()
	rc.log.Debugf("connect test %s", st)
	return rc
}

// UpdateCfg
// 热加载配置
func (rc *RedisClient) UpdateCfg(cfg interface{}) {
}

/* --------------- redis key字符串操作 ----------------- */
// Keys
func (rc *RedisClient) Keys(pattern string) ([]string, error) {
	keys, err := rc.client.Keys(rc.ctx, pattern).Result()
	return keys, err
}

// Del
func (rc *RedisClient) Del(keys ...string) (int64, error) {
	affects, err := rc.client.Del(rc.ctx, keys...).Result()
	return affects, err
}

// Exists
func (rc *RedisClient) Exists(keys ...string) (int64, error) {
	affects, err := rc.client.Exists(rc.ctx, keys...).Result()
	return affects, err
}

// Scan
func (rc *RedisClient) Scan(cursor uint64, match string, count int64) ([]string, uint64, error) {
	keys, nextCursor, err := rc.client.Scan(rc.ctx, cursor, match, count).Result()
	return keys, nextCursor, err
}

// Pexpire
func (rc *RedisClient) Pexpire(key string, expiration time.Duration) (bool, error) {
	status, err := rc.client.PExpire(rc.ctx, key, expiration).Result()
	return status, err
}

// Pttl
func (rc *RedisClient) Pttl(key string) (time.Duration, error) {
	expire, err := rc.client.PTTL(rc.ctx, key).Result()
	return expire, err
}

/* -------------- redis string操作 ------------------- */
// Set
func (rc *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	s, err := rc.client.Set(rc.ctx, key, value, expiration).Result()
	if err == nil && s != "OK" {
		return errors.New(s)
	}
	return err
}

// Get
func (rc *RedisClient) Get(key string) (string, error) {
	val, err := rc.client.Get(rc.ctx, key).Result()
	return val, err
}

// Incrby
func (rc *RedisClient) Incrby(key string, i int64) (int64, error) {
	val, err := rc.client.IncrBy(rc.ctx, key, i).Result()
	return val, err
}

// Setnx
func (rc *RedisClient) Setnx(key string, value interface{}, expiration time.Duration) (bool, error) {
	status, err := rc.client.SetNX(rc.ctx, key, value, expiration).Result()
	return status, err
}

// Setex
func (rc *RedisClient) Setex(key string, value interface{}, expiration time.Duration) (bool, error) {
	status, err := rc.client.SetXX(rc.ctx, key, value, expiration).Result()
	return status, err
}

/* ----------------- redis set集合操作 ------------------ */
// Sadd
func (rc *RedisClient) Sadd(key string, values ...interface{}) (int64, error) {
	newMemberCnt, err := rc.client.SAdd(rc.ctx, key, values...).Result()
	return newMemberCnt, err
}

// Smembers
func (rc *RedisClient) Smembers(key string) ([]string, error) {
	members, err := rc.client.SMembers(rc.ctx, key).Result()
	return members, err
}

/* ----------------- redis list列表操作 ---------------- */
// Lpush
func (rc *RedisClient) Lpush(key string, values ...interface{}) (int64, error) {
	totalSize, err := rc.client.LPush(rc.ctx, key, values...).Result()
	return totalSize, err
}

// Rpush
func (rc *RedisClient) Rpush(key string, values ...interface{}) (int64, error) {
	totalSize, err := rc.client.RPush(rc.ctx, key, values...).Result()
	return totalSize, err
}

// Lrange
func (rc *RedisClient) Lrange(key string, first, last int64) ([]string, error) {
	sectionData, err := rc.client.LRange(rc.ctx, key, first, last).Result()
	return sectionData, err
}

// Ltrim
func (rc *RedisClient) Ltrim(key string, first, last int64) error {
	s, err := rc.client.LTrim(rc.ctx, key, first, last).Result()
	if err == nil && s != "OK" {
		return errors.New(s)
	}
	return err
}

// Brpop
func (rc *RedisClient) Brpop(timeout time.Duration, keys ...string) ([]string, error) {
	vals, err := rc.client.BRPop(rc.ctx, timeout, keys...).Result()
	return vals, err
}

/* ------------------ reids hash哈希结构操作 ------------------ */
// Hset
func (rc *RedisClient) Hset(key string, values ...interface{}) error {
	_, err := rc.client.HSet(rc.ctx, key, values...).Result()
	return err
}

// Hmget
func (rc *RedisClient) Hmget(key string, fields ...string) ([]interface{}, error) {
	vals, err := rc.client.HMGet(rc.ctx, key, fields...).Result()
	return vals, err
}

// Hmset
func (rc *RedisClient) Hmset(key string, values ...interface{}) error {
	_, err := rc.client.HMSet(rc.ctx, key, values...).Result()
	return err
}
func (rc *RedisClient) Hincrby(key string, field string, i int64) (int64, error) {
	curVal, err := rc.client.HIncrBy(rc.ctx, key, field, i).Result()
	return curVal, err
}

/* ---------------- redis transactions事务操作接口 -------------------- */
/* ---------------- redis pipline管道操作接口 -------------------- */

/* ---------------- redis script脚本操作接口 -------------------- */
// Eval
func (rc *RedisClient) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	val, err := rc.client.Eval(rc.ctx, script, keys, args...).Result()
	return val, err
}

// EvalSha
func (rc *RedisClient) EvalSha(sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	val, err := rc.client.EvalSha(rc.ctx, sha1, keys, args...).Result()
	return val, err
}

// ScriptLoad
func (rc *RedisClient) ScriptLoad(script string) (string, error) {
	sha, err := rc.client.ScriptLoad(rc.ctx, script).Result()
	return sha, err
}

// ScriptExists
func (rc *RedisClient) ScriptExists(hashes ...string) ([]bool, error) {
	status, err := rc.client.ScriptExists(rc.ctx, hashes...).Result()
	return status, err
}

// ScriptFlush
func (rc *RedisClient) ScriptFlush() error {
	s, err := rc.client.ScriptFlush(rc.ctx).Result()
	if err == nil && s != "OK" {
		return errors.New(s)
	}
	return nil
}

// ScriptKill
func (rc *RedisClient) ScriptKill() error {
	s, err := rc.client.ScriptKill(rc.ctx).Result()
	if err == nil && s != "OK" {
		return errors.New(s)
	}
	return err
}

/* ------------------ redis 命令接口 --------------- */
func (rc *RedisClient) Ping() (string, error) {
	s, err := rc.client.Ping(rc.ctx).Result()
	return s, err
}
func (rc *RedisClient) Do(args ...interface{}) (interface{}, error) {
	tv, err := rc.client.Do(rc.ctx, args...).Result()
	return tv, err
}

// RedisClient
type RedisClient struct {
	client *redis.ClusterClient
	log    logger.LogI
	ctx    context.Context
}

const (
	Nil = redis.Nil // redis when key not exists
)
