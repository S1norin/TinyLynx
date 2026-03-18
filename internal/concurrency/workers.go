package concurrency

import (
	"context"
	"sync"
	"time"
)

type AnalyticsEvent struct {
	LinkID    int
	IPAddress string
	UserAgent string
	Referrer  string
	Country   string
	Device    string
	Browser   string
	Platform  string
}

type WorkerPool struct {
	workerCount int
	taskQueue   chan AnalyticsEvent
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewWorkerPool(workerCount int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workerCount: workerCount,
		taskQueue:   make(chan AnalyticsEvent, 1000),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case event := <-wp.taskQueue:
			// Process analytics event
			processAnalyticsEvent(event)
		}
	}
}

func (wp *WorkerPool) Submit(event AnalyticsEvent) {
	select {
	case wp.taskQueue <- event:
		// Successfully queued
	default:
		// Queue is full, drop the event
		// In production, you might want to log this
	}
}

func (wp *WorkerPool) Shutdown() {
	wp.cancel()
	wp.wg.Wait()
	close(wp.taskQueue)
}

func processAnalyticsEvent(event AnalyticsEvent) {
	// This is a placeholder - in a real implementation,
	// you would call your analytics service here
	// For now, we just sleep to simulate processing
	time.Sleep(10 * time.Millisecond)
}

// Rate Limiter
type RequestInfo struct {
	Time time.Time
}

type RateLimiter struct {
	requests map[string][]RequestInfo
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]RequestInfo),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if ipRequests, exists := rl.requests[ip]; exists {
		// Check if the oldest request is outside the window
		if len(ipRequests) > 0 && now.Sub(ipRequests[0].Time) > rl.window {
			rl.requests[ip] = ipRequests[1:]
		}
	}

	if len(rl.requests[ip]) >= rl.limit {
		return false
	}

	rl.requests[ip] = append(rl.requests[ip], RequestInfo{Time: now})
	return true
}

// Simple in-memory rate limiter for demo purposes
// In production, consider using Redis for distributed rate limiting

























































































































// In production, consider using Redis for distributed rate limiting// Simple in-memory rate limiter for demo purposes}	Time time.Timetype RequestInfo struct {}	return true	rl.requests[ip] = append(rl.requests[ip], RequestInfo{Time: now})	}		return false	if len(rl.requests[ip]) >= rl.limit {	}		}			rl.requests[ip] = ipRequests[1:]		if now.Sub(ipRequests[0].Time) > rl.window {		// Check if the oldest request is outside the window	if ipRequests, exists := rl.requests[ip]; exists {	now := time.Now()	defer rl.mu.Unlock()	rl.mu.Lock()func (rl *RateLimiter) Allow(ip string) bool {}	}		window:   window,		limit:    limit,		requests: make(map[string]int),	return &RateLimiter{func NewRateLimiter(limit int, window time.Duration) *RateLimiter {}	window   time.Duration	limit    int	mu       sync.Mutex	requests map[string]inttype RateLimiter struct {// Rate Limiter}	time.Sleep(10 * time.Millisecond)	// For now, we just sleep to simulate processing	// you would call your analytics service here	// This is a placeholder - in a real implementation,func processAnalyticsEvent(event AnalyticsEvent) {}	close(wp.taskQueue)	wp.wg.Wait()	wp.cancel()func (wp *WorkerPool) Shutdown() {}	}		// In production, you might want to log this		// Queue is full, drop the event	default:		// Successfully queued	case wp.taskQueue <- event:	select {func (wp *WorkerPool) Submit(event AnalyticsEvent) {}	}		}			processAnalyticsEvent(event)			// Process analytics event		case event := <-wp.taskQueue:			return		case <-wp.ctx.Done():		select {	for {	defer wp.wg.Done()func (wp *WorkerPool) worker() {}	}		go wp.worker()		wp.wg.Add(1)	for i := 0; i < wp.workerCount; i++ {func (wp *WorkerPool) Start() {}	}		cancel:      cancel,		ctx:         ctx,		taskQueue:   make(chan AnalyticsEvent, 1000),		workerCount: workerCount,	return &WorkerPool{	ctx, cancel := context.WithCancel(context.Background())func NewWorkerPool(workerCount int) *WorkerPool {}	cancel      context.CancelFunc	ctx         context.Context	wg          sync.WaitGroup	taskQueue   chan AnalyticsEvent	workerCount inttype WorkerPool struct {}	Platform  string	Browser   string	Device    string	Country   string	Referrer  string	UserAgent string	IPAddress string	LinkID    inttype AnalyticsEvent struct {)	"time"	"sync"	"context"import (