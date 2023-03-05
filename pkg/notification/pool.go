package notification

//go:generate mockgen -package mock -destination=./mock/pool_mock.go . Pool
type Pool interface {
	// Start gets the workerpool ready to process jobs, and should only be called once
	Start()
	// Stop stops the workerpool, tears down any required resources,
	// and should only be called once
	Stop()
	// AddWork adds a task for the worker pool to process. It is only valid after
	// Start() has been called and before Stop() has been called.
	AddWork(Task)
}
