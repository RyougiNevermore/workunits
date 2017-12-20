package workunits

import (
	"context"
	"time"
)

// at call
func WithFuture(ctx context.Context) (context.Context, error) {
	// ctx.withValue(some key, chan)
	return nil, nil
}

// at do, work.Do(...) > ctx = WithPromise(ctx), ctx.ack(v)
func WithPromise(ctx context.Context) (context.Context, error) {
	// ctx.Value(some key)
	return nil, nil
}

type promiseContext struct {
	context.Context
	ch chan interface{}
}

func (p *promiseContext) Put(v interface{}) error {

	return nil
}

type futureContext struct {
	context.Context
	ch chan interface{}
}

func (p *futureContext) Get() (interface{}, error) {

	return nil, nil
}

func (p *futureContext) GetTimeout(timeout time.Duration) (interface{}, error) {

	return nil, nil
}

func (p *futureContext) GetDeadline(timeout time.Time) (interface{}, error) {

	return nil, nil
}
