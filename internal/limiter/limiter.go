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
	l          rate.Limit // limit
	b          int        // burst
	access     map[string]*Access
	sync.Mutex //@TODO : RWMutex ?
}

func New(l rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		l:      l,
		b:      b,
		access: make(map[string]*Access),
	}
}

func (rl *RateLimiter) GetAccess(ip string) *rate.Limiter {
	rl.Lock()
	defer rl.Unlock()

	a, ok := rl.access[ip]
	if !ok {
		l := rate.NewLimiter(rl.l, rl.b)
		rl.access[ip] = &Access{
			lat:     time.Now(),
			limiter: l,
		}
		return l
	}
	a.lat = time.Now()
	return a.limiter
}

// func (rl *RateLimiter) Cleanup() {
// 	for {
// 		time.Sleep(time.Minute)
// 		rl.Lock()
// 		for ip, a := range rl.access {
// 			if time.Since(a.lat) > 10*time.Minute {
// 				delete(rl.access, ip)
// 			}
// 		}
// 		rl.Unlock()
// 	}
// }
