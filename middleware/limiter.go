package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/thansetan/pendekin/helper"
)

type bucket struct {
	maxToken, token int64
	lastRefillTime  time.Time
	interval        time.Duration
}

func newBucket(limit int64, interval time.Duration) *bucket {
	bucket := &bucket{
		maxToken:       limit,
		token:          limit,
		interval:       interval,
		lastRefillTime: time.Now(),
	}

	return bucket
}

func (b *bucket) take() {
	b.token--
}

func (b *bucket) allow() bool {
	if b.canRefill() {
		b.refill()
	}
	return b.token > 0
}

func (b *bucket) refill() {
	b.token = b.maxToken
	b.lastRefillTime = time.Now()
}

func (b *bucket) canRefill() bool {
	return time.Now().After(b.lastRefillTime.Add(b.interval))
}

type rateLimiter struct {
	limit      int64
	ipMap      map[any]*bucket
	mu         sync.Mutex
	interval   time.Duration
	contextKey any
}

func NewRateLimiter(limit int64, interval time.Duration, contextKey any) *rateLimiter {
	rl := &rateLimiter{
		limit:      limit,
		ipMap:      make(map[any]*bucket),
		interval:   interval,
		contextKey: contextKey,
	}

	return rl
}

func (rl *rateLimiter) RateLimitMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Context().Value(rl.contextKey)
		rl.mu.Lock()
		if _, ok := rl.ipMap[key]; !ok {
			rl.ipMap[key] = newBucket(rl.limit, rl.interval)
		}

		bucket := rl.ipMap[key]
		if !bucket.allow() {
			rl.mu.Unlock()
			helper.ResponseBuilder[any](w, http.StatusTooManyRequests, "rate limit exceeded", nil)
			return
		}
		bucket.take()
		rl.mu.Unlock()
		f.ServeHTTP(w, r)
	}
}
