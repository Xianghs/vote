package util

import (
	_ "github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	_ "net/http"
	"sync"
)

// IPRateLimiter .
type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     *sync.RWMutex
	limit rate.Limit
	burst  int
}

// NewIPRateLimiter .
func NewIPRateLimiter(l float64, b int) *IPRateLimiter {
	s := rate.Limit(l)
	i := &IPRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		mu:     &sync.RWMutex{},
		limit: s,
		burst:  b,
	}

	return i
}

// AddIP 创建了一个新的速率限制器，并将其添加到 ips 映射中,
// 使用 IP地址作为密钥
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.limit, i.burst)

	i.ips[ip] = limiter

	return limiter
}

// GetLimiter 返回所提供的IP地址的速率限制器(如果存在的话).
// 否则调用 AddIP 将 IP 地址添加到映射中
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	i.mu.Unlock()

	return limiter
}
