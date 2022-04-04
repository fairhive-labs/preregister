package limiter

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

var errNotAllowed = errors.New("not allowed")

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

	if _, ok := limiter.access[ip]; !ok {
		t.Errorf("access map should contain %s", ip)
		t.FailNow()
	}

	if !lmt.Allow() { // consumes 1 token
		t.Errorf("ip %s should be allowed and consume 1 token", ip)
		t.FailNow()
	}

}

func TestBurst(t *testing.T) {
	ip := "10.10.10.10"
	tt := []struct {
		l   rate.Limit
		b   int
		err error
	}{
		{5.0, 3, nil},
		{5.0, 5, nil},
		{10.0, 10, nil},
		{2.0, 3, errNotAllowed},
		{99.0, 100, errNotAllowed},
	}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("l=%.2f b=%d", tc.l, tc.b), func(t *testing.T) {
			limiter := New(tc.l, tc.b)
			n := tc.b + 1
			simulateNRequests(ip, n, limiter)
			time.Sleep(time.Second) //pause, add tokens in the bucket
			if err := simulateNRequests(ip, n, limiter); !errors.Is(err, tc.err) {
				t.Errorf("incorrect error: got %v, want %v", err, tc.err)
				t.FailNow()
			}
		})
	}
}

func simulateNRequests(ip string, n int, limiter *RateLimiter) error {
	for i := 0; i < n; i++ {
		lmt := limiter.GetAccess(ip)
		if !lmt.Allow() && i != limiter.burst { // consumes 1 token
			return errNotAllowed
		}
	}
	return nil
}
