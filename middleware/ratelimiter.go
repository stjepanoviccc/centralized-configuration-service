package middleware

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mux      sync.Mutex
	interval time.Duration
	limit    int
	counters map[string]int
}

func NewRateLimiter(interval time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		interval: interval,
		limit:    limit,
		counters: make(map[string]int),
	}
}

func (rl *RateLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIPAddress(r)

		rl.mux.Lock()
		defer rl.mux.Unlock()

		count, ok := rl.counters[ip]
		if !ok {
			count = 1
		} else {
			count++
		}

		if count > rl.limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		fmt.Println("Counters map:", rl.counters)

		rl.counters[ip] = count
		time.AfterFunc(rl.interval, func() {
			rl.mux.Lock()
			defer rl.mux.Unlock()
			delete(rl.counters, ip)
		})

		next.ServeHTTP(w, r)
	})
}

func getIPAddress(r *http.Request) string {
	fullAddr := r.RemoteAddr

	ip, _, err := net.SplitHostPort(fullAddr)
	if err != nil {
		return fullAddr
	}

	return ip
}

func AdaptHandler(handler http.Handler, limiter *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limiter.Limit(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)
		}).ServeHTTP(w, r)
	}
}
