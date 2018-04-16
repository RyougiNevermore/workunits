package workunits

// worker group
type WorkerGroup interface {
	Start() (err error)
	Close() (err error)
	Sync() (err error)
	Send(u Unit) (err error)
}

// job unit
type Unit interface {
	Process()
}

// worker
type Worker interface {
	Work()
}
