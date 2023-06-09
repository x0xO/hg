package surf

import (
	"context"
	"runtime"
	"sync"

	"github.com/x0xO/hg/surf/hsyscall"
	"go.uber.org/ratelimit"
)

// Requests a struct that manages concurrent HTTP requests.
type Requests struct {
	rateLimiter ratelimit.Limiter
	ctx         context.Context
	jobs        chan *asyncRequest
	maxWorkers  int
}

// Do performs all queued requests concurrently, returning channels with results and errors.
func (reqs *Requests) Do() (chan *asyncResponse, chan error) {
	maxWorkers := _maxWorkers

	if reqs.maxWorkers != 0 {
		if runtime.GOOS != "windows" {
			reqs.maxWorkers = hsyscall.RlimitStack(reqs.maxWorkers)
		}

		maxWorkers = reqs.maxWorkers
	}

	if reqs.rateLimiter == nil {
		reqs.rateLimiter = ratelimit.NewUnlimited()
	}

	results := make(chan *asyncResponse)
	errors := make(chan error)

	wg := sync.WaitGroup{}

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for job := range reqs.jobs {
				reqs.rateLimiter.Take()

				resp, err := job.
					SetHeaders(job.setHeaders).
					AddHeaders(job.addHeaders).
					AddCookies(job.addCookies...).
					WithContext(reqs.ctx).
					Do()
				if err != nil {
					errors <- err
					continue
				}

				results <- &asyncResponse{resp, job.context}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	return results, errors
}

// RateLimiter sets a rate limiter for the concurrent requests, limiting the number of requests per
// second.
func (reqs *Requests) RateLimiter(maxRequestsPerSecond int) *Requests {
	if maxRequestsPerSecond > 0 {
		reqs.rateLimiter = ratelimit.New(maxRequestsPerSecond)
	}

	return reqs
}

// Pool sets the number of worker goroutines for the concurrent requests.
func (reqs *Requests) Pool(workers int) *Requests {
	reqs.maxWorkers = workers
	return reqs
}

// WithContext associates the provided context with the concurrent requests.
func (reqs *Requests) WithContext(ctx context.Context) *Requests {
	reqs.ctx = ctx
	return reqs
}
