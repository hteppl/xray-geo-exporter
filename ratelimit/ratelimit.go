package ratelimit

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu         sync.Mutex
	rate       int
	tokens     int
	maxTokens  int
	lastRefill time.Time
}

func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	return &RateLimiter{
		rate:       requestsPerMinute,
		tokens:     requestsPerMinute,
		maxTokens:  requestsPerMinute,
		lastRefill: time.Now(),
	}
}

func (r *RateLimiter) Wait() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastRefill)

	// Refill tokens based on elapsed time
	tokensToAdd := int(elapsed.Seconds() * float64(r.rate) / 60.0)
	if tokensToAdd > 0 {
		r.tokens = min(r.tokens+tokensToAdd, r.maxTokens)
		r.lastRefill = now
	}

	// If no tokens available, wait
	if r.tokens <= 0 {
		waitTime := time.Duration(60.0/float64(r.rate)) * time.Second
		time.Sleep(waitTime)
		r.tokens = 1
		r.lastRefill = time.Now()
	}

	r.tokens--
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
