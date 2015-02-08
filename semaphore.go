package goesl

type Semaphore chan bool

func (s Semaphore) Lock() {
	<-s
}

func (s Semaphore) Unlock() {
	s <- true
}

func NewSemaphore(concurrency uint) Semaphore {
	s := make(Semaphore, concurrency)

	var i uint

	for i = 0; i < concurrency; i++ {
		s.Unlock()
	}

	return s
}
