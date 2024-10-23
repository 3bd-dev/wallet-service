package rest

import (
	"context"
	"time"
)

type Retrier interface {
	Do(context.Context, func(ctx context.Context, attempt int) error) error
}

func NewRetrier(maxAttempts int, retryDelay time.Duration) Retrier {
	if maxAttempts < 1 {
		panic("maxAttempts should be at least 1")
	}
	return &retrier{
		maxAttempts: maxAttempts,
		retryDelay:  retryDelay,
	}
}

type retrier struct {
	maxAttempts int
	retryDelay  time.Duration
}

func (r *retrier) Do(ctx context.Context, fn func(ctx context.Context, attempt int) error) error {
	for attempt := 1; attempt <= r.maxAttempts; attempt++ {
		err := fn(ctx, attempt)

		if err == nil || attempt == r.maxAttempts {
			return err
		}

		t := time.NewTimer(r.retryDelay)
		select {
		case <-t.C:
		case <-ctx.Done():
			if !t.Stop() {
				<-t.C
			}
			return ctx.Err()
		}
	}
	return nil
}
