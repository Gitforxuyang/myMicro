package conn

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

var (
	ErrInvalidConfig = errors.New("invalid config")
	ErrPoolClosed    = errors.New("pool closed")
)

type CustomCloser struct {
}

func (m *CustomCloser) Close() error {
	return nil
}

type factory func() (io.Closer, error)

type Pool interface {
	New() (io.Closer, error)
	Release(io.Closer) error
	Close(io.Closer) error
	Shutdown() error
}

type GenericPool struct {
	sync.Mutex
	pool        chan io.Closer
	maxOpen     int
	numOpen     int
	minOpen     int
	closed      bool
	maxLifetime time.Duration
	factory     factory
}

func (m *GenericPool) Shutdown() error {
	panic("implement me")
}

func NewPool(minOpen int, maxOpen int, maxLifeTime time.Duration, factory factory) (Pool, error) {
	p := &GenericPool{
		maxLifetime: maxLifeTime,
		maxOpen:     maxOpen,
		minOpen:     minOpen,
		factory:     factory,
		pool:        make(chan io.Closer, maxOpen),
	}
	for i := 0; i < minOpen; i++ {
		closer, _ := factory()
		p.numOpen++
		p.pool <- closer
	}

	return p, nil
}

func (m *GenericPool) New() (io.Closer, error) {
	if m.closed {
		return nil, ErrPoolClosed
	}
	closer, _ := m.getOrCreate()
	return closer, nil
}
func (m *GenericPool) getOrCreate() (io.Closer, error) {
	select {
	case closer := <-m.pool:
		fmt.Println(closer)
		fmt.Println("直接拿老的")

		return closer, nil
	default:

	}
	//超过最大连接数时派对。应该加上等待时间
	if m.numOpen >= m.maxOpen {
		select {
		case closer := <-m.pool:
			return closer, nil
		case <-time.After(time.Second * 3):
			fmt.Println("超时消息")
			return nil, errors.New("等待超时")
		}

	}
	m.Lock()
	closer, _ := m.factory()
	m.numOpen++
	m.Unlock()
	fmt.Println("建立新的")
	return closer, nil
}

func (m *GenericPool) Release(closer io.Closer) error {
	m.Lock()
	m.pool <- closer
	m.Unlock()
	return nil
}
func (m *GenericPool) Close(closer io.Closer) error {
	m.Lock()
	closer.Close()
	m.numOpen--
	m.Unlock()
	return nil
}
