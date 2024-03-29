package usecase

import (
	"fmt"
	"sync"

	"golang.org/x/time/rate"
)

type Limiter interface {
	Allow() bool
}

type Limit float64

func NewLimiter(r Limit, b int) Limiter { //TODO: implement the interface in order to avoid usage of rate library
	return rate.NewLimiter(rate.Limit(r), b) //FIXME: check rate limits
}

type IPRateLimiter struct {
	ips map[string]Limiter //FIXME: sync.Map
	mu  *sync.RWMutex
	r   Limit
	b   int
}

func NewIPRateLimiter(r Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
	return i
}

func (i *IPRateLimiter) AddIP(ip string) Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()
	limiter := NewLimiter(i.r, i.b)
	i.ips[ip] = limiter
	return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) Limiter {
	i.mu.Lock()
	fmt.Println(ip)
	limiter, exists := i.ips[ip]
	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}
	i.mu.Unlock()
	return limiter
}
