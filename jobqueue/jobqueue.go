package jobqueue

import (
	"sync"
)

type Queue struct {
	list  []func()
	mutex sync.Mutex
}

func New() *Queue {
	return &Queue{}
}

func (self *Queue) Run() {
	self.mutex.Lock()
	batch := self.list
	self.list = nil
	self.mutex.Unlock()
	for _, fn := range batch {
		fn()
	}
}

func (self *Queue) Job(fn func()) {
	self.mutex.Lock()
	self.list = append(self.list, fn)
	self.mutex.Unlock()
}
