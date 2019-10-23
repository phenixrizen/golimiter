package golimiter

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// visitor has a individual rate limiter and the last time seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// LimitHTTP is a http.Handler that limits requests by the Limiters values
func (l *Limiter) LimitHTTP(next http.Handler) http.Handler {
	if !l.cleanup {
		// go routine to remove old entries from the visitors sync map
		go l.cleanupVisitors()
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if l.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LimitHTTPByIP is a http.Handler that limits requests by individual users rates
func (l *Limiter) LimitHTTPByIP(next http.Handler) http.Handler {
	if !l.cleanup {
		// go routine to remove old entries from the visitors sync map
		go l.cleanupVisitors()
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := l.getVisitor(r.RemoteAddr)
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LimitHTTPByHeader is a http.Handler that limits requests by individual users rates
// Examples of keys would be "X-Forwarded-For", "X-Real-Ip", etc.
func (l *Limiter) LimitHTTPByHeader(key string, next http.Handler) http.Handler {
	if !l.cleanup {
		// go routine to remove old entries from the visitors sync map
		go l.cleanupVisitors()
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(key)
		limiter := l.getVisitor(header)
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// addsVisitor adds a new visitor to the sync map
func (l *Limiter) addVisitor(key string) *rate.Limiter {
	limiter := rate.NewLimiter(l.limiter.Limit(), l.limiter.Burst())
	visitor := &visitor{limiter, time.Now()}
	l.visitors.Store(key, visitor)
	return limiter
}

// getVisitor gets a visitor by key from the sync map
func (l *Limiter) getVisitor(key string) *rate.Limiter {
	v, exists := l.visitors.Load(key)
	if !exists {
		return l.addVisitor(key)
	}
	vis := v.(visitor)
	vis.lastSeen = time.Now()
	return vis.limiter
}

// cleanupVisitors cleans the sync map if a visitor hasn't been seen for the cleanup interval
func (l *Limiter) cleanupVisitors() {
	for {
		time.Sleep(l.cleanupInterval)
		l.visitors.Range(func(key interface{}, v interface{}) bool {
			vis := v.(visitor)
			if time.Now().Sub(vis.lastSeen) > l.cleanupInterval {
				l.visitors.Delete(key)
			}
			return true
		})
	}
}
