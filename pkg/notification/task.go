package notification

//go:generate mockgen -package mock -destination=./mock/task_mock.go . Task
type Task interface {
	// Execute performs the work
	Execute() error
	// OnFailure handles any error returned from Execute()
	OnFailure(error)
}
