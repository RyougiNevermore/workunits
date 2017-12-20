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
	Put(worker Worker) (err error)
	Get() (worker Worker, err error)
	GetTimeout(timeout time.Duration) (worker Worker, err error)
	GetDeadline(deadline time.Time) (worker Worker, err error)
	Start() (err error)
	Shutdown(ctx context.Context) (err error)
}
