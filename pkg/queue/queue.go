package queue

import (
	"context"
	"sync"
	"time"
)

type Queue[T any] struct {
	items []T
	mu    sync.Mutex
}

// NewQueue creates a new Queue instance
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0),
	}
}

// Enqueue adds an item to the queue (thread-safe)
func (q *Queue[T]) Enqueue(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

// Dequeue removes an item from the queue (thread-safe)
// It returns the dequeued item and a bool indicating success or failure
func (q *Queue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Handle empty queue case
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}

	// Dequeue the first item
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// StartWorker starts a worker to process items from the queue asynchronously
func (q *Queue[T]) StartWorker(ctx context.Context, processFunc func(context.Context, T)) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				item, ok := q.Dequeue()
				if !ok {
					time.Sleep(1 * time.Second)
					continue
				}
				processFunc(ctx, item)
			}
		}
	}()
}
