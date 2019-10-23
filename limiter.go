package golimiter

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Limiter controls how frequently events are allowed to happen.
// It uses the golang x time rate library which implements a
// token bucket" of size b, initially full and refilled at rate
// r tokens per second. Informally, in any large enough time
// interval, the Limiter limits the rate to r tokens per second,
// with a maximum burst size of b events. As a special case,
// if r == Inf (the infinite rate), b is ignored.
// See https://en.wikipedia.org/wiki/Token_bucket for more about token buckets.
type Limiter struct {
	limiter         *rate.Limiter
	visitors        sync.Map
	cleanup         bool
	cleanupInterval time.Duration
}

type Reservation struct {
	*rate.Reservation
}

// New returns a new limiter that will limit http requests
func New(r float64, b int) *Limiter {
	limiter := &Limiter{
		limiter: rate.NewLimiter(rate.Limit(r), b),
	}
	return limiter
}

// Allow is shorthand for AllowN(time.Now(), 1).
func (l *Limiter) Allow() bool {
	return l.limiter.Allow()
}

// AllowN reports whether n events may happen at time now.
// Use this method if you intend to drop / skip events that exceed
// the rate limit. Otherwise use Reserve or Wait.
func (l *Limiter) AllowN(t time.Time, n int) bool {
	return l.limiter.AllowN(t, n)
}

// Burst returns the maximum burst size. Burst is the maximum number
// of tokens that can be consumed in a single call to Allow, Reserve,
// or Wait, so higher Burst values allow more events to happen at once.
// A zero Burst allows no events, unless limit == Inf.
func (l *Limiter) Burst() int {
	return l.limiter.Burst()
}

// Limit returns the maximum overall event rate.
func (l *Limiter) Limit() float64 {
	return float64(l.limiter.Limit())
}

// Reserve is shorthand for ReserveN(time.Now(), 1).
func (l *Limiter) Reserve() *Reservation {
	r := l.limiter.Reserve()
	return &Reservation{Reservation: r}
}

// ReserveN  returns a Reservation that indicates how long the caller
// must wait before n events happen. The Limiter takes this Reservation
// into account when allowing future events. ReserveN returns false if
// n exceeds the Limiter's burst size.
func (l *Limiter) ReserveN(ctx context.Context, t time.Time, n int) *Reservation {
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	r := l.limiter.ReserveN(t, n)
	return &Reservation{Reservation: r}
}

// SetBurst is shorthand for SetBurstAt(time.Now(), newBurst).
func (l *Limiter) SetBurst(nb int) {
	l.limiter.SetBurst(nb)
}

// SetBurstAt sets a new burst size for the limiter.
func (l *Limiter) SetBurstAt(t time.Time, nb int) {
	l.limiter.SetBurstAt(t, nb)
}

// SetLimit is shorthand for SetLimitAt(time.Now(), newLimit).
func (l *Limiter) SetLimit(nl float64) {
	l.limiter.SetLimit(rate.Limit(nl))
}

// SetLimitAt sets a new Limit for the limiter. The new Limit,
// and Burst, may be violated or underutilized by those which
// reserved (using Reserve or Wait) but did not yet act before
// SetLimitAt was called.
func (l *Limiter) SetLimitAt(t time.Time, nl float64) {
	l.limiter.SetLimitAt(t, rate.Limit(nl))
}

// Wait is shorthand for WaitN(ctx, 1).
func (l *Limiter) Wait(ctx context.Context) error {
	c, cancel := context.WithCancel(ctx)
	defer cancel()
	return l.limiter.Wait(c)
}

// WaitN WaitN blocks until lim permits n events to happen. It
// returns an error if n exceeds the Limiter's burst size, the
// Context is canceled, or the expected wait time exceeds the
// Context's Deadline. The burst limit is ignored if the rate
// limit is Inf.
func (l *Limiter) WaitN(ctx context.Context, n int) error {
	c, cancel := context.WithCancel(ctx)
	defer cancel()
	return l.limiter.WaitN(c, n)
}
