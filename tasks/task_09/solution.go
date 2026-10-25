package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type Result[R any] struct {
	Value R
	Err   error
}

var ErrInvalidWorkers = errors.New("workers must be positive")

func ParallelMap[T any, R any](ctx context.Context, workers int, in []T, fn func(context.Context, T) (R, error)) ([]Result[R], error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkers
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	results := make([]Result[R], len(in))
	var next atomic.Int64
	var wg sync.WaitGroup

	for range min(workers, len(in)) {
		wg.Go(func() {
			for ctx.Err() == nil {
				i := int(next.Add(1)) - 1
				if i >= len(in) {
					return
				}
				v, err := fn(ctx, in[i])
				if err != nil {
					results[i].Err = err
				} else {
					results[i].Value = v
				}
			}
		})
	}
	wg.Wait()

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return results, nil
}
