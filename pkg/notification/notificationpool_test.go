//go:build unit

package notification

import (
	"errors"
	"log"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iboware/http-notifier/pkg/notification/mock"
)

func TestWorkerPool_NewNotificationPool(t *testing.T) {
	log.SetFlags(0)

	if _, err := NewNotificationPool(log.Default(), 0, 0); err != ErrNoWorkers {
		t.Fatalf("expected error with 0 workers, got: %v", err)
	}
	if _, err := NewNotificationPool(log.Default(), -1, 0); err != ErrNoWorkers {
		t.Fatalf("expected error with -1 workers, got: %v", err)
	}
	if _, err := NewNotificationPool(log.Default(), 1, -1); err != ErrNegativeChannelSize {
		t.Fatalf("expected error with -1 channel size, got: %v", err)
	}

	p, err := NewNotificationPool(log.Default(), 5, 0)
	if err != nil {
		t.Fatalf("expected no error creating notification pool, got: %v", err)
	}
	if p == nil {
		t.Fatal("NewNotificationPool returned nil for valid input")
	}
}

func TestWorkerPool_Work(t *testing.T) {
	var tasks []*mock.MockTask
	wg := &sync.WaitGroup{}

	// set logger to 0 so it will act as a NOP Logger
	log.SetFlags(0)

	// create a mock controller for the mock tasks
	ctrl := gomock.NewController(t)

	for i := 0; i < 20; i++ {
		// create mock tasks and add one to wait group
		mt := mock.NewMockTask(ctrl)
		wg.Add(1)

		if i == 10 {
			// failure case
			mt.EXPECT().Execute().DoAndReturn(func() error {
				return errors.New("error")
			}).Times(1)
			mt.EXPECT().OnFailure(gomock.Any()).DoAndReturn(func(err error) {
				wg.Done()
			}).Times(1)
		} else {
			// success case
			mt.EXPECT().Execute().Do(func() error {
				wg.Done()
				return nil
			}).Times(1)
		}
		tasks = append(tasks, mt)
	}

	//initilize and start the pool
	p, err := NewNotificationPool(log.Default(), 1, 1)
	if err != nil {
		t.Fatal("error making notification pool:", err)
	}
	p.Start()

	for _, j := range tasks {
		// add tasks to the pool in non-blocking way
		p.AddWorkNonBlocking(j)
	}

	// block until the tasks are done
	wg.Wait()
}
