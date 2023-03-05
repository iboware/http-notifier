package notification

import (
	"fmt"
	"log"
	"sync"
)

type NotificationPool struct {
	numWorkers int
	tasks      chan Task
	logger     *log.Logger

	// ensure the pool can only be started or stopped once
	start sync.Once
	stop  sync.Once

	// channel to make workers stop
	quit chan struct{}
}

var ErrNoWorkers = fmt.Errorf("no workers")
var ErrNegativeChannelSize = fmt.Errorf("negative channel size is not allowed")

func NewNotificationPool(logger *log.Logger, numWorkers, channelSize int) (*NotificationPool, error) {
	if numWorkers <= 0 {
		return nil, ErrNoWorkers
	}
	if channelSize < 0 {
		return nil, ErrNegativeChannelSize
	}

	tasks := make(chan Task, channelSize)

	return &NotificationPool{
		numWorkers: numWorkers,
		tasks:      tasks,
		logger:     logger,

		start: sync.Once{},
		stop:  sync.Once{},

		quit: make(chan struct{}),
	}, nil
}

func (p *NotificationPool) Start() {
	p.start.Do(func() {
		p.logger.Println("starting notification pool")
		p.startWorkers()
	})
}

func (p *NotificationPool) Stop() {
	p.stop.Do(func() {
		p.logger.Println("stopping notification pool")
		close(p.quit)
	})
}

// AddWork adds work to the NotificationPool. If the channel buffer is full (or 0) and
// all workers are occupied, this will hang until work is consumed or Stop() is called.
func (p *NotificationPool) AddWork(t Task) {
	select {
	case p.tasks <- t:
	case <-p.quit:
	}
}

// AddWorkNonBlocking adds work to the NotificationPool and returns immediately
func (p *NotificationPool) AddWorkNonBlocking(t Task) {
	go p.AddWork(t)
}

func (p *NotificationPool) startWorkers() {
	for i := 0; i < p.numWorkers; i++ {
		go func(workerNum int) {
			p.logger.Printf("starting worker %d", workerNum)
			for {
				select {
				case <-p.quit:
					p.logger.Printf("stopping worker %d with quit channel\n", workerNum)
					return
				case task := <-p.tasks:
					if err := task.Execute(); err != nil {
						task.OnFailure(err)
					}
				}
			}
		}(i)
	}
}
