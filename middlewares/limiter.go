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
	limit      int
	ipMap      map[any]*bucket
	mu         sync.Mutex
	lastRefill time.Time
	interval   time.Duration
}

func NewRateLimiter(limit int, interval time.Duration) *rateLimiter {
	rl := &rateLimiter{
		limit: limit,
		ipMap: make(map[any]*bucket),

		// not the best practice, but at least it works!
		lastRefill: time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(), time.Now().UTC().Day(), 0, 0, 0, 0, time.UTC),
		interval:   interval,
	}

	go rl.refillAll()

	return rl
}

func (rl *rateLimiter) RateLimitMiddleware(f http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
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

}

func (rl *rateLimiter) refillAll() {
	for {
		now := time.Now().UTC()
		if now.Sub(rl.lastRefill).Nanoseconds() >= rl.interval.Nanoseconds() {
			rl.mu.Lock()
			for ip := range rl.ipMap {
				rl.ipMap[ip].refill()
			}
			rl.mu.Unlock()
			rl.lastRefill = rl.lastRefill.Add(rl.interval)
		}
	}
}
