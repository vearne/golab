package mylock

import (
	"github.com/vearne/second-realize/rwlock"
	"time"
)

type MyRWLock struct {
	count int
	mu    *rwlock.RWLocker
}

func NewMyRWLock() *MyRWLock {
	lock := MyRWLock{}
	lock.mu = rwlock.NewRWLocker(rwlock.WithMaxWaitWriteGoroutine(10), rwlock.WithMaxReadLocked(10))
	return &lock
}

func (l *MyRWLock) Write() {
	l.mu.WLock()
	l.count++
	time.Sleep(cost)
	l.mu.WUnLock()
}

func (l *MyRWLock) Read() {
	l.mu.RLock()
	_ = l.count
	time.Sleep(cost)
	l.mu.RUnLock()
}
