package limiter

import (
	"testing"

	"golang.org/x/time/rate"
)

func TestNew(t *testing.T) {
	var l rate.Limit = 5.0
	b := 3
	limiter := New(l, b)
	if limiter.limit != l {
		t.Errorf("incorrect limit, got %f, want %f", limiter.limit, l)
		t.FailNow()
	}
	if limiter.burst != b {
		t.Errorf("incorrect burst, got %d, want %d", limiter.burst, b)
		t.FailNow()
	}
	if limiter.access == nil {
		t.Errorf("access map cannot be nil")
		t.FailNow()
	}
}

func TestNewUnlimited(t *testing.T) {
	limiter := NewUnlimited()
	if limiter.limit != rate.Inf {
		t.Errorf("incorrect limit, got %f, want %f", limiter.limit, rate.Inf)
		t.FailNow()
	}
	if limiter.burst != 0 {
		t.Errorf("incorrect burst, got %d, want %d", limiter.burst, 0)
		t.FailNow()
	}
	if limiter.access == nil {
		t.Errorf("access map cannot be nil")
		t.FailNow()
	}
}

func TestGetAccess(t *testing.T) {
	var l rate.Limit = 5.0
	b := 3
	limiter := New(l, b)
	ip := "10.10.10.10"

	if _, ok := limiter.access[ip]; ok {
		t.Errorf("access map should not contain %s", ip)
		t.FailNow()
	}

	lmt := limiter.GetAccess(ip)
	if lmt == nil {
		t.Errorf("rate limiter for ip %s cannot be nil", ip)
		t.FailNow()
	}

	if !lmt.Allow() {
		t.Errorf("ip %s should be allowed", ip)
		t.FailNow()
	}

	if _, ok := limiter.access[ip]; !ok {
		t.Errorf("access map should contain %s", ip)
		t.FailNow()
	}

	// for i := 1; i <= 10; i++ {
	// 	lmt = limiter.GetAccess(ip)
	// 	if !lmt.Allow() {
	// 		t.Errorf("ip %s should be allowed on request %d", ip, i)
	// 		t.FailNow()
	// 	}
	// }
}
