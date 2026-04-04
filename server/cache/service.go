package cache

import (
	"log"
	"sync"
	"time"
)

// CacheKey 缓存 Key 定义
type CacheKey struct {
	Prefix string
	Expire time.Duration
}

const (
	// 帖子缓存
	TopicPrefix    = "topic"
	TopicExpire    = 5 * time.Minute

	// 用户相关缓存
	UserBadgesPrefix = "user:badges"
	UserBadgesExpire = 10 * time.Minute

	UserInfoPrefix = "user:info"
	UserInfoExpire = 10 * time.Minute

	// 板块缓存
	ForumPrefix = "forum"
	ForumExpire = 30 * time.Minute

	// 标签缓存
	TagPrefix = "tag"
	TagExpire = 30 * time.Minute

	// 首页帖子列表缓存
	TopicListPrefix = "topic:list"
	TopicListExpire = 1 * time.Minute
)

// 缓存锁，防止缓存击穿
var locks sync.Map

// BuildKey 构建缓存 key
func BuildKey(prefix string, id interface{}) string {
	return prefix + ":" + toString(id)
}

// toString 转换 id 为字符串
func toString(v interface{}) string {
	switch val := v.(type) {
	case uint:
		return uintToString(val)
	case int:
		return intToString(val)
	case string:
		return val
	default:
		return ""
	}
}

func uintToString(v uint) string {
	return string(rune('0'+v%10)) + uintToString(v/10)
}

func intToString(v int) string {
	if v < 0 {
		return "-" + intToString(-v)
	}
	if v < 10 {
		return string(rune('0' + v))
	}
	return intToString(v/10) + string(rune('0'+v%10))
}

// GetData 通用的缓存获取
// key: 缓存键
// fetchFunc: 缓存不存在时的获取函数
// expire: 过期时间
func GetData(key string, fetchFunc func() (interface{}, error), expire time.Duration) (interface{}, error) {
	// 尝试从缓存获取
	if val, ok := Get(key); ok {
		return val, nil
	}

	// 缓存不存在，执行 fetchFunc
	data, err := fetchFunc()
	if err != nil {
		return nil, err
	}

	// 写入缓存
	Set(key, data, expire)
	return data, nil
}

// GetDataWithLock 带锁的缓存获取，防止缓存击穿
func GetDataWithLock(key string, fetchFunc func() (interface{}, error), expire time.Duration) (interface{}, error) {
	// 尝试从缓存获取
	if val, ok := Get(key); ok {
		return val, nil
	}

	// 加锁
	lockKey := "lock:" + key
	mutex, _ := locks.LoadOrStore(lockKey, &sync.Mutex{})
	mu := mutex.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	// 双重检查
	if val, ok := Get(key); ok {
		return val, nil
	}

	// 缓存不存在，执行 fetchFunc
	data, err := fetchFunc()
	if err != nil {
		return nil, err
	}

	// 写入缓存
	Set(key, data, expire)
	return data, nil
}

// InvalidateCache 失效缓存
func InvalidateCache(prefix string, ids ...interface{}) {
	for _, id := range ids {
		key := BuildKey(prefix, id)
		Delete(key)
	}
}

// InvalidatePrefix 失效指定前缀的所有缓存（通过遍历）
func InvalidatePrefix(prefix string) {
	log.Printf("cache: invalidating prefix %s (note: ristretto doesn't support pattern delete)", prefix)
	// Ristretto 不支持模式删除，日志提示
}

// TopicCache 帖子缓存操作
var TopicCache = &CacheOps{
	Prefix: TopicPrefix,
	Expire: TopicExpire,
}

// UserBadgesCache 用户勋章缓存操作
var UserBadgesCache = &CacheOps{
	Prefix: UserBadgesPrefix,
	Expire: UserBadgesExpire,
}

// ForumCache 板块缓存操作
var ForumCache = &CacheOps{
	Prefix: ForumPrefix,
	Expire: ForumExpire,
}

// TagCache 标签缓存操作
var TagCache = &CacheOps{
	Prefix: TagPrefix,
	Expire: TagExpire,
}

// TopicListCache 帖子列表缓存操作
var TopicListCache = &CacheOps{
	Prefix: TopicListPrefix,
	Expire: TopicListExpire,
}

// CacheOps 通用缓存操作
type CacheOps struct {
	Prefix string
	Expire time.Duration
}

// Get 获取缓存
func (c *CacheOps) Get(id interface{}) (interface{}, bool) {
	return Get(c.BuildKey(id))
}

// Set 设置缓存
func (c *CacheOps) Set(id interface{}, data interface{}) {
	Set(c.BuildKey(id), data, c.Expire)
}

// BuildKey 构建缓存 key
func (c *CacheOps) BuildKey(id interface{}) string {
	return BuildKey(c.Prefix, id)
}

// Invalidate 失效缓存
func (c *CacheOps) Invalidate(id interface{}) {
	Delete(c.BuildKey(id))
}

// GetData 获取缓存数据，不存在时调用 fetchFunc
func (c *CacheOps) GetData(id interface{}, fetchFunc func() (interface{}, error)) (interface{}, error) {
	return GetData(c.BuildKey(id), fetchFunc, c.Expire)
}
