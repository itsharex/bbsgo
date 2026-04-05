package utils

import (
	"sync"
	"time"
)

// LoginAttempt 登录尝试记录
type LoginAttempt struct {
	Username  string
	Count     int
	LockUntil time.Time
	LastTry   time.Time
}

// LoginLimiter 登录限制器
type LoginLimiter struct {
	attempts map[string]*LoginAttempt
	mu       sync.RWMutex
}

var loginLimiter = &LoginLimiter{
	attempts: make(map[string]*LoginAttempt),
}

// CheckLoginAttempt 检查登录尝试
// 返回是否允许登录和剩余锁定时间
func CheckLoginAttempt(username string) (allowed bool, remainingTime time.Duration) {
	loginLimiter.mu.Lock()
	defer loginLimiter.mu.Unlock()

	attempt, exists := loginLimiter.attempts[username]
	if !exists {
		return true, 0
	}

	// 检查是否被锁定
	if time.Now().Before(attempt.LockUntil) {
		remaining := time.Until(attempt.LockUntil)
		return false, remaining
	}

	// 锁定已过期，重置计数
	if !attempt.LockUntil.IsZero() && time.Now().After(attempt.LockUntil) {
		delete(loginLimiter.attempts, username)
		return true, 0
	}

	return true, 0
}

// RecordLoginFailure 记录登录失败
func RecordLoginFailure(username string) {
	loginLimiter.mu.Lock()
	defer loginLimiter.mu.Unlock()

	attempt, exists := loginLimiter.attempts[username]
	if !exists {
		attempt = &LoginAttempt{
			Username: username,
			Count:    0,
		}
		loginLimiter.attempts[username] = attempt
	}

	attempt.Count++
	attempt.LastTry = time.Now()

	// 根据失败次数设置锁定时间
	// 5次失败：锁定15分钟
	// 10次失败：锁定30分钟
	// 15次失败：锁定1小时
	var lockDuration time.Duration
	switch {
	case attempt.Count >= 15:
		lockDuration = 1 * time.Hour
	case attempt.Count >= 10:
		lockDuration = 30 * time.Minute
	case attempt.Count >= 5:
		lockDuration = 15 * time.Minute
	default:
		// 不锁定
		return
	}

	attempt.LockUntil = time.Now().Add(lockDuration)
}

// RecordLoginSuccess 记录登录成功
func RecordLoginSuccess(username string) {
	loginLimiter.mu.Lock()
	defer loginLimiter.mu.Unlock()

	delete(loginLimiter.attempts, username)
}

// GetLoginAttempts 获取登录失败次数
func GetLoginAttempts(username string) int {
	loginLimiter.mu.RLock()
	defer loginLimiter.mu.RUnlock()

	attempt, exists := loginLimiter.attempts[username]
	if !exists {
		return 0
	}

	return attempt.Count
}

// ClearExpiredAttempts 清理过期的登录尝试记录
func ClearExpiredAttempts() {
	loginLimiter.mu.Lock()
	defer loginLimiter.mu.Unlock()

	now := time.Now()
	for username, attempt := range loginLimiter.attempts {
		// 如果锁定时间已过且超过1小时未尝试，则删除记录
		if !attempt.LockUntil.IsZero() && now.After(attempt.LockUntil) {
			delete(loginLimiter.attempts, username)
		} else if attempt.LockUntil.IsZero() && now.Sub(attempt.LastTry) > 1*time.Hour {
			delete(loginLimiter.attempts, username)
		}
	}
}

// init 初始化定时清理任务
func init() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			ClearExpiredAttempts()
		}
	}()
}
