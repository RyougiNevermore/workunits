package workunits

import (
	"context"
	"sync"
)

type Unit interface {
	Context() context.Context
	Worker() Worker
}

type defaultUnit struct {
	ctx    context.Context
	worker Worker
}

func (u *defaultUnit) Context() context.Context {
	return u.ctx
}

func (u *defaultUnit) Worker() Worker {
	return u.worker
}

var unitPool = sync.Pool{
	New: func() interface{} {
		return new(defaultUnit)
	},
}

func getUnit() Unit {
	return unitPool.Get().(Unit)
}

func releaseUnit(u *defaultUnit) {
	u.worker = nil
	u.ctx = nil
	unitPool.Put(u)
}
