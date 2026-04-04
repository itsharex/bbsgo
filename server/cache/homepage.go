package cache

import (
	"log"
	"sync"
	"time"
)

// HomePageCache 首页分区缓存模块
// 每个模块独立缓存，独立失效/更新
var HomePageCache = &homePageCache{
	forumsCache:       &sectionCache{prefix: "home:forums", expire: 30 * time.Minute},
	tagsCache:         &sectionCache{prefix: "home:tags", expire: 10 * time.Minute},
	announcementsCache: &sectionCache{prefix: "home:announcements", expire: 5 * time.Minute},
	topicsCache:       &sectionCache{prefix: "home:topics", expire: 1 * time.Minute},
}

type homePageCache struct {
	forumsCache       *sectionCache
	tagsCache         *sectionCache
	announcementsCache *sectionCache
	topicsCache       *sectionCache
}

// sectionCache 单个分区缓存
type sectionCache struct {
	prefix string
	expire time.Duration
	mu     sync.RWMutex
}

// GetKey 获取缓存key
func (s *sectionCache) GetKey() string {
	return s.prefix
}

// GetTTL 获取过期时间
func (s *sectionCache) GetTTL() time.Duration {
	return s.expire
}

// Get 获取分区缓存
func (s *sectionCache) Get() (interface{}, bool) {
	return Get(s.prefix)
}

// Set 设置分区缓存
func (s *sectionCache) Set(data interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	Set(s.prefix, data, s.expire)
}

// Invalidate 使分区缓存失效
func (s *sectionCache) Invalidate() {
	s.mu.Lock()
	defer s.mu.Unlock()
	Delete(s.prefix)
	log.Printf("homepage cache: invalidated %s", s.prefix)
}

// TopicsPageCache 获取指定页的话题缓存
func (s *sectionCache) TopicsPageKey(page int, pageSize int) string {
	return s.prefix + ":" + itoa(page) + ":" + itoa(pageSize)
}

// TopicsPageGet 获取指定页的话题缓存
func (s *sectionCache) TopicsPageGet(page, pageSize int) (interface{}, bool) {
	return Get(s.TopicsPageKey(page, pageSize))
}

// TopicsPageSet 设置指定页的话题缓存
func (s *sectionCache) TopicsPageSet(page, pageSize int, data interface{}) {
	key := s.TopicsPageKey(page, pageSize)
	Set(key, data, s.expire)
}

// TopicsPageInvalidate 使所有话题页缓存失效
func (s *sectionCache) TopicsPageInvalidate() {
	// Ristretto不支持模式删除，遍历清除常用页
	for page := 1; page <= 10; page++ {
		for _, size := range []int{10, 20, 50} {
			Delete(s.TopicsPageKey(page, size))
		}
	}
	log.Printf("homepage cache: invalidated all topics pages")
}

// Forwards 获取所有分区缓存的get方法
func (h *homePageCache) Forums() *sectionCache    { return h.forumsCache }
func (h *homePageCache) Tags() *sectionCache      { return h.tagsCache }
func (h *homePageCache) Announcements() *sectionCache { return h.announcementsCache }
func (h *homePageCache) Topics() *sectionCache    { return h.topicsCache }

// InvalidateAll 使所有分区缓存失效
func (h *homePageCache) InvalidateAll() {
	h.forumsCache.Invalidate()
	h.tagsCache.Invalidate()
	h.announcementsCache.Invalidate()
	h.topicsCache.Invalidate()
	log.Println("homepage cache: all sections invalidated")
}

// InvalidateTopics 只失效话题分区（发帖/评论/删帖/置顶时调用）
func (h *homePageCache) InvalidateTopics() {
	h.topicsCache.TopicsPageInvalidate()
	log.Println("homepage cache: topics section invalidated")
}

// InvalidateAnnouncements 只失效公告分区
func (h *homePageCache) InvalidateAnnouncements() {
	h.announcementsCache.Invalidate()
	log.Println("homepage cache: announcements section invalidated")
}

// itoa 数字转字符串（简单实现）
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	var result []byte
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	return string(result)
}
