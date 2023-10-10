package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/thansetan/pendekin/helpers"
)

type bucket struct {
	maxToken, token int
}

func newBucket(limit int) *bucket {
	return &bucket{
		maxToken: limit,
		token:    limit,
	}
}

func (b *bucket) take() {
	b.token--
}

func (b *bucket) allow() bool {
	return b.token > 0
}

func (b *bucket) refill() {
	b.token = b.maxToken
}

type rateLimiter struct {
	limit    int
	ipMap    map[any]*bucket
	mu       sync.Mutex
	refillAt time.Time
}

func NewRateLimiter(limit int, refillAt time.Time) *rateLimiter {
	rl := &rateLimiter{
		limit:    limit,
		ipMap:    make(map[any]*bucket),
		refillAt: refillAt,
	}

	go rl.refillAll()

	return rl
}

func (rl *rateLimiter) RateLimitMiddleware(f http.HandlerFunc, w http.ResponseWriter, r *http.Request) {

	rl.mu.Lock()
	defer rl.mu.Unlock()
	clientIP := r.Context().Value("user_ip")

	if _, ok := rl.ipMap[clientIP]; !ok {
		rl.ipMap[clientIP] = newBucket(rl.limit)
	}

	bucket := rl.ipMap[clientIP]

	if !bucket.allow() {
		helpers.ResponseBuilder(w, http.StatusTooManyRequests, "rate limit exceeded", nil)
		return
	}

	bucket.take()
	f.ServeHTTP(w, r)

}

func (rl *rateLimiter) refillAll() {
	for {
		now := time.Now().UTC()
		if now.Sub(rl.refillAt).Hours() >= 24 {
			rl.mu.Lock()
			for ip := range rl.ipMap {
				rl.ipMap[ip].refill()
			}
			rl.mu.Unlock()
			rl.refillAt = rl.refillAt.Add(24 * time.Hour)
		}
	}
}

func (rl *rateLimiter) GetUsers() map[any]*bucket {
	return rl.ipMap
}
