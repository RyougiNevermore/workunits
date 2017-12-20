package workunits

import "context"

type WorkerGroup interface {
	Start() (err error)
	Shutdown(ctx context.Context) (err error)
	Execute(ctx context.Context, worker Worker) (err error)
}
