// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

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
