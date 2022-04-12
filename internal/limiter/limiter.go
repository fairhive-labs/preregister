package limiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Access struct {
	lat     time.Time //last access tieme
	limiter *rate.Limiter
}

type RateLimiter struct {
	limit      rate.Limit
	burst      int
	access     map[string]*Access
	sync.Mutex //@TODO : RWMutex ?
}

func New(l rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limit:  l,
		burst:  b,
		access: make(map[string]*Access),
	}
}

func NewUnlimited() *RateLimiter {
	return New(rate.Inf, 0)
}

func (rl *RateLimiter) GetAccess(ip string) *rate.Limiter {
	rl.Lock()
	defer rl.Unlock()

	a, ok := rl.access[ip]
	if !ok {
		l := rate.NewLimiter(rl.limit, rl.burst)
		rl.access[ip] = &Access{
			lat:     time.Now(),
			limiter: l,
		}
		return l
	}
	a.lat = time.Now()
	return a.limiter
}

func (rl *RateLimiter) Cleanup(t time.Duration) {
	rl.Lock()
	defer rl.Unlock()

	for ip, a := range rl.access {
		if time.Since(a.lat) > t {
			delete(rl.access, ip)
		}
	}

}
