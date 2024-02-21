package middleware

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/thansetan/pendekin/helper"
)

func TestRateLimitMiddleware(t *testing.T) {
	rl := NewRateLimiter(2, 3*time.Second, helper.UserIPKey)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handlerResponse := "hi"
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, handlerResponse)
	}
	t.Run("FirstRequestShouldNotGetRateLimited", func(t *testing.T) {
		rec := httptest.NewRecorder()
		GetClientIP(rl.RateLimitMiddleware(handler)).ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status code to be %d, got %d\n", http.StatusOK, res.StatusCode)
		}

		data, _ := io.ReadAll(res.Body)
		if string(data) != handlerResponse {
			t.Errorf("expected response body to be %s, got %s\n", handlerResponse, string(data))
		}
	})

	t.Run("SecondRequestShouldNotButThirdRequestShouldGetRateLimited", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			rec := httptest.NewRecorder()
			GetClientIP(rl.RateLimitMiddleware(handler)).ServeHTTP(rec, req)
			res := rec.Result()
			if i == 0 {
				if res.StatusCode != http.StatusOK {
					t.Errorf("expected second request status code to be %d, got %d\n", http.StatusOK, res.StatusCode)
				}
			} else {
				if res.StatusCode != http.StatusTooManyRequests {
					t.Errorf("expected third request status code to be %d, got %d\n", http.StatusTooManyRequests, res.StatusCode)
				}
			}
		}
	})
	t.Run("FourthRequestShouldProceed", func(t *testing.T) {
		time.Sleep(5 * time.Second) // wait for rate limit to refill the token
		rec := httptest.NewRecorder()
		GetClientIP(rl.RateLimitMiddleware(handler)).ServeHTTP(rec, req)
		res := rec.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("expected fourth request status code to be %d, got %d\n", http.StatusOK, res.StatusCode)
		}
	})
}

func TestRateLimiterConcurrentAccess(t *testing.T) {
	var (
		limit, multiplier int64 = 1_000, 5
		interval                = time.Second
		requestDuration         = interval * time.Duration(multiplier)
	)
	rl := NewRateLimiter(limit, interval, helper.UserIPKey)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	ok, totalRequest := new(atomic.Int64), new(atomic.Int64)
	var wg sync.WaitGroup
	timer := time.NewTimer(requestDuration)
loop:
	for {
		select {
		case <-timer.C:
			timer.Stop()
			break loop
		default:
			wg.Add(1)
			go func() {
				defer wg.Done()
				rec := httptest.NewRecorder()
				GetClientIP(rl.RateLimitMiddleware(handler)).ServeHTTP(rec, req)
				res := rec.Result()
				totalRequest.Add(1)
				switch res.StatusCode {
				case http.StatusOK:
					ok.Add(1)
				}
			}()
		}
	}
	wg.Wait()

	expected200 := limit + (limit * multiplier)
	if ok.Load() != expected200 {
		t.Errorf("expected %d 200 OK, got %d\n", expected200, ok.Load())
	}

	t.Logf("request was performed for %.2f seconds with %d (200 OK) of total %d requests\n", requestDuration.Seconds(), ok.Load(), totalRequest.Load())
}
