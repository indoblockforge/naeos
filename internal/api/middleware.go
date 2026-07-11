package api

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.window)

	// Clean old requests
	if reqs, ok := r.requests["global"]; ok {
		valid := make([]time.Time, 0)
		for _, t := range reqs {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}
		r.requests["global"] = valid
	}

	// Check limit
	if len(r.requests["global"]) >= r.limit {
		return false
	}

	// Add request
	r.requests["global"] = append(r.requests["global"], now)
	return true
}

func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests = make(map[string][]time.Time)
}
