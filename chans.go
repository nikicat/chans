package chans

import (
	"context"
	"iter"
)

// Receive one item from channel or timeout on ctx
func RecvChan[T any](ctx context.Context, ch <-chan T) (msg T, err error) {
	select {
	case msg = <-ch:
		return
	case <-ctx.Done():
		// if ctx is Done then try to read without blocking anyway
		select {
		case msg = <-ch:
			return
		default:
		}
		err = context.Cause(ctx)
		return
	}
}

func SendChan[T any](ctx context.Context, ch chan<- T, msg T) (err error) {
	select {
	case ch <- msg:
		return
	case <-ctx.Done():
		err = context.Cause(ctx)
		return
	}
}

// Read items from channel until channel is closed or ctx is cancelled
func ReadChan[T any](ctx context.Context, ch <-chan T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-ch:
				if !ok {
					return
				}
				if !yield(val) {
					return
				}
			}
		}
	}
}

func SafeClose[T any](ch chan<- T) {
	if ch != nil {
		close(ch)
	}
}

func TrySend[T any](ch chan<- T, t T) bool {
	select {
	case ch <- t:
		return true
	default:
		return false
	}
}
