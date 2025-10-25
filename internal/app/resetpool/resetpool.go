package resetpool

import (
	"sync"
)

// Resetter объекты которые можно сбрасывать.
type Resetter interface {
	Reset()
}

// Pool пул объектов Resetter.
type Pool[T Resetter] struct {
	pool sync.Pool
}

// New создает новый пул для типа Resetter.
func New[T Resetter]() *Pool[T] {
	p := &Pool[T]{}
	p.pool.New = func() interface{} {
		return new(T)
	}
	return p
}

// Get возвращает Resetter из пула.
func (p *Pool[T]) Get() *T {
	return p.pool.Get().(*T)
}

// Put помещает Resetter в пул.
func (p *Pool[T]) Put(x *T) {
	p.pool.Put(x)
}
