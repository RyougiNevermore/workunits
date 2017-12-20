package workunits

import (
	"context"
	"fmt"
	"time"
)

var loopErrStarted = fmt.Errorf("loop is started")
var loopErrStopped = fmt.Errorf("loop is stopped")
var loopErrNilWorker = fmt.Errorf("worker is nil")
var loopErrTimeout = fmt.Errorf("get worker failed, timeout")
var loopErrEmpty = fmt.Errorf("loop is empty")

type Loop interface {
	Put(u Unit) (err error)
	Get() (u Unit, err error)
	GetTimeout(timeout time.Duration) (u Unit, err error)
	GetDeadline(deadline time.Time) (u Unit, err error)
	Start() (err error)
	Shutdown(ctx context.Context) (err error)
}
