package workerpool

type Pool struct {
	Job  func(func())
	Wait func()
}

func New(limit int) *Pool {
	self := &Pool{}

	sem := make(chan struct{}, limit)

	self.Job = func(fn func()) {
		sem <- struct{}{}
		go func() {
			fn()
			<-sem
		}()
	}

	self.Wait = func() {
		for i := 0; i < limit; i++ {
			sem <- struct{}{}
		}
		for i := 0; i < limit; i++ {
			<-sem
		}
	}

	return self
}
