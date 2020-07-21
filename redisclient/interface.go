// rediscliet
package redisclient

import "time"

// KeyI redis 键操作接口定义
type KeyI interface {
	// Keys 【todo】 仅能获取一个节点上匹配的key，无法获取全部key
	Keys(pattern string) ([]string, error)
	// Del 删除数据库中的键
	Del(keys ...string) (int64, error)
	// Exists 判断数据库中键是否存在
	Exists(keys ...string) (int64, error)
	// Scan 【todo】迭代数据库中的数据库键；仅能获取一个节点上匹配的key，无法获取全部key
	Scan(cursor uint64, match string, count int64) ([]string, uint64, error)
	// Pexpire 设置键的存活时间
	Pexpire(key string, expiration time.Duration) (bool, error)
	// Pttl 键值的剩余存活时间
	Pttl(key string) (time.Duration, error)
}

// StringI
type StringI interface {
	// Set
	Set(key string, value interface{}, expiration time.Duration) error
	// Get
	Get(key string) (string, error)
	// Incrby
	Incrby(key string, i int64) (int64, error)
	// Setnx 如果键不存在则设置
	Setnx(key string, value interface{}, expiration time.Duration) (bool, error)
	// Setex 设置键的到期时间
	Setex(key string, value interface{}, expiration time.Duration) (bool, error)
}

type SetI interface {
	// Sadd 向集合添加多个元素
	Sadd(key string, values ...interface{}) (int64, error)
	// 获取集合中的全部元素
	Smembers(key string) ([]string, error)
}

type ListI interface {
	// Lpush 左边(低位)插入数据
	Lpush(key string, values ...interface{}) (int64, error)
	// Lrange 获取闭区间[first,last]数据
	Lrange(key string, first, last int64) ([]string, error)
	// Ltrim 保留[first,last]闭区间数据，其他数据异常
	Ltrim(key string, first, last int64) error
}

// HashI 哈希数据结构接口
type HashI interface {
	Hmget(key string, fields ...string) ([]interface{}, error)
	Hmset(key string, kvs map[string]interface{}) error
	Hincrby(key string, field string, i int64) (int64, error)
}

// TransactionsI 事务接口
type TransactionsI interface {
}

// PipelineI 管道接口
type PipelineI interface {
}

// ScriptI 脚本接口
type ScriptI interface {
	// Eval  执行指定的lua脚本
	Eval(script string, keys []string, args ...interface{}) (interface{}, error)
	// EvalSha 执行通过sha1指向的lua脚本
	EvalSha(sha1 string, keys []string, args ...interface{}) (interface{}, error)
	// ScriptLoad 缓存脚本，返回脚本对应点sha值
	ScriptLoad(script string) (string, error)
	ScriptExists(hashes ...string) ([]bool, error)
	ScriptFlush() error
	ScriptKill() error
}

type CommandI interface {
	Ping() (string, error)
}

type RedisClientI interface {
	KeyI
	StringI
	SetI
	ListI
	HashI
	TransactionsI
	PipelineI
	ScriptI
	CommandI
}
