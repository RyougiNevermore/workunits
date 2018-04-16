package workunits

import "context"

// worker group
type WorkerGroup interface {
	Start() (err error)
	Close() (err error)
	Sync(ctx context.Context)
	Send(u Unit) (err error)
}
