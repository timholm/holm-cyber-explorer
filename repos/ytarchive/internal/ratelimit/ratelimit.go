package ratelimit

import (
	"context"
	"sync"
	"time"
)

// Limiter implements a token bucket rate limiter
type Limiter struct {
	rate       float64   // tokens per second
	burst      int       // maximum burst size
	tokens     float64   // current tokens
	lastUpdate time.Time // last token update time
	mu         sync.Mutex
}

// NewLimiter creates a new rate limiter
// rate: requests per second
// burst: maximum number of requests that can be made at once
func NewLimiter(rate float64, burst int) *Limiter {
	return &Limiter{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastUpdate: time.Now(),
	}
}

// Wait blocks until a token is available or the context is cancelled
func (l *Limiter) Wait(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(l.lastUpdate).Seconds()
	l.tokens += elapsed * l.rate
	if l.tokens > float64(l.burst) {
		l.tokens = float64(l.burst)
	}
	l.lastUpdate = now

	// If we have a token, use it immediately
	if l.tokens >= 1 {
		l.tokens--
		return nil
	}

	// Calculate wait time for next token
	waitTime := time.Duration((1 - l.tokens) / l.rate * float64(time.Second))

	l.mu.Unlock()
	select {
	case <-ctx.Done():
		l.mu.Lock()
		return ctx.Err()
	case <-time.After(waitTime):
		l.mu.Lock()
		l.tokens = 0
		l.lastUpdate = time.Now()
		return nil
	}
}

// Allow checks if a request can be made without blocking
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Refill tokens
	now := time.Now()
	elapsed := now.Sub(l.lastUpdate).Seconds()
	l.tokens += elapsed * l.rate
	if l.tokens > float64(l.burst) {
		l.tokens = float64(l.burst)
	}
	l.lastUpdate = now

	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}

// Reserve returns a Reservation that indicates how long the caller must wait before n events happen
func (l *Limiter) Reserve() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Refill tokens
	now := time.Now()
	elapsed := now.Sub(l.lastUpdate).Seconds()
	l.tokens += elapsed * l.rate
	if l.tokens > float64(l.burst) {
		l.tokens = float64(l.burst)
	}
	l.lastUpdate = now

	if l.tokens >= 1 {
		l.tokens--
		return 0
	}

	waitTime := time.Duration((1 - l.tokens) / l.rate * float64(time.Second))
	l.tokens = 0
	return waitTime
}

// SetRate updates the rate limit
func (l *Limiter) SetRate(rate float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.rate = rate
}

// SetBurst updates the burst limit
func (l *Limiter) SetBurst(burst int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.burst = burst
	if l.tokens > float64(burst) {
		l.tokens = float64(burst)
	}
}
