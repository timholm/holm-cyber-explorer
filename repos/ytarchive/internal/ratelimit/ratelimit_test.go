package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestLimiter_Allow(t *testing.T) {
	// Create limiter with 10 requests per second, burst of 5
	limiter := NewLimiter(10, 5)

	// First 5 requests should be allowed immediately (burst)
	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Errorf("Request %d should have been allowed", i+1)
		}
	}

	// 6th request should be denied (burst exhausted)
	if limiter.Allow() {
		t.Error("6th request should have been denied")
	}
}

func TestLimiter_Wait(t *testing.T) {
	// Create limiter with 10 requests per second, burst of 2
	limiter := NewLimiter(10, 2)

	ctx := context.Background()

	// First 2 requests should complete immediately
	start := time.Now()
	for i := 0; i < 2; i++ {
		if err := limiter.Wait(ctx); err != nil {
			t.Errorf("Wait() error = %v", err)
		}
	}
	elapsed := time.Since(start)
	if elapsed > 50*time.Millisecond {
		t.Errorf("First 2 requests took too long: %v", elapsed)
	}

	// 3rd request should wait for ~100ms (1/10 second)
	start = time.Now()
	if err := limiter.Wait(ctx); err != nil {
		t.Errorf("Wait() error = %v", err)
	}
	elapsed = time.Since(start)
	if elapsed < 50*time.Millisecond {
		t.Errorf("3rd request should have waited, but only took %v", elapsed)
	}
}

func TestLimiter_WaitCancelled(t *testing.T) {
	// Create limiter with very slow rate
	limiter := NewLimiter(0.1, 1) // 1 request per 10 seconds

	// Exhaust burst
	limiter.Allow()

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Wait should return immediately with error
	err := limiter.Wait(ctx)
	if err == nil {
		t.Error("Wait() should have returned error for cancelled context")
	}
}

func TestLimiter_Reserve(t *testing.T) {
	// Create limiter with 10 requests per second, burst of 2
	limiter := NewLimiter(10, 2)

	// First 2 reservations should return 0 wait time
	for i := 0; i < 2; i++ {
		wait := limiter.Reserve()
		if wait > 0 {
			t.Errorf("Reservation %d should have 0 wait time, got %v", i+1, wait)
		}
	}

	// 3rd reservation should require waiting
	wait := limiter.Reserve()
	if wait < 50*time.Millisecond {
		t.Errorf("3rd reservation should require waiting, got %v", wait)
	}
}

func TestLimiter_SetRate(t *testing.T) {
	limiter := NewLimiter(1, 1)

	// Exhaust burst
	limiter.Allow()

	// Change to faster rate
	limiter.SetRate(100)

	// Should allow requests quickly now
	time.Sleep(50 * time.Millisecond)
	if !limiter.Allow() {
		t.Error("Should allow after rate increase")
	}
}

func TestLimiter_SetBurst(t *testing.T) {
	limiter := NewLimiter(1, 5)

	// Use some tokens
	for i := 0; i < 3; i++ {
		limiter.Allow()
	}

	// Reduce burst to 2 (less than current tokens)
	limiter.SetBurst(2)

	// Should only allow 2 more at most
	allowed := 0
	for i := 0; i < 5; i++ {
		if limiter.Allow() {
			allowed++
		}
	}

	if allowed > 2 {
		t.Errorf("Should have allowed at most 2 after burst reduction, got %d", allowed)
	}
}

func TestLimiter_TokenRefill(t *testing.T) {
	// Create limiter with 10 requests per second, burst of 2
	limiter := NewLimiter(10, 2)

	// Exhaust burst
	limiter.Allow()
	limiter.Allow()

	// Should not allow immediately
	if limiter.Allow() {
		t.Error("Should not allow immediately after exhausting burst")
	}

	// Wait for tokens to refill (100ms for 1 token at 10/sec)
	time.Sleep(150 * time.Millisecond)

	// Should allow now
	if !limiter.Allow() {
		t.Error("Should allow after waiting for refill")
	}
}
