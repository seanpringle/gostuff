package jobqueue

import (
	"sync"
)

type Queue struct {
	Run func()
	Job func(func())
}

func New() *Queue {

	var mutex sync.Mutex
	var list []func()

	return &Queue{
		Run: func() {
			mutex.Lock()
			batch := list
			list = nil
			mutex.Unlock()
			for _, fn := range batch {
				fn()
			}
		},
		Job: func(fn func()) {
			mutex.Lock()
			list = append(list, fn)
			mutex.Unlock()
		},
	}
}
