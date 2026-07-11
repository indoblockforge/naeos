package api

import (
	"testing"
	"time"
)

func TestRateLimiterAllow(t *testing.T) {
	limiter := NewRateLimiter(5, time.Minute)

	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Errorf("expected request %d to be allowed", i+1)
		}
	}

	if limiter.Allow() {
		t.Error("expected request to be denied after limit")
	}
}

func TestRateLimiterReset(t *testing.T) {
	limiter := NewRateLimiter(2, time.Minute)

	limiter.Allow()
	limiter.Allow()

	if limiter.Allow() {
		t.Error("expected request to be denied")
	}

	limiter.Reset()

	if !limiter.Allow() {
		t.Error("expected request to be allowed after reset")
	}
}
