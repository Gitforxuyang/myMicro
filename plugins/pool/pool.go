package pool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type NeoPool interface {
	Ping(ctx context.Context) error
}

type neoPool struct {
	maxActive     int
	minActive     int
	maxIdle       int
	currentActive int
	//最大排队数
	maxWait        int
	currentWait    int
	waitTimeOut    time.Duration
	requestTimeOut time.Duration
	//工厂函数
	factory func() (Conn, error)
	sync.Mutex
	ch chan Conn
}

var (
	ContextDoneError = errors.New("context done 错误")
	WaitTimeOutError = errors.New("等待连接池超时")
	ReqTimeOutError  = errors.New("请求超时")
	WaitMaxError     = errors.New("排队数量超标")
)

//从连接池中获取
//TODO:还要做连接池未满时自动创建
func (m *neoPool) Get() (Conn, error) {
	select {
	case conn := <-m.ch:
		return conn, nil
	default:

	}
	//如果当前连接数小于最大连接数。则创建新连接
	m.Lock()
	//fmt.Println(m.currentActive)
	if m.currentActive < m.maxActive {
		conn, err := m.factory()
		if err != nil {
			m.Unlock()
			return nil, err
		}
		m.currentActive++
		m.Unlock()
		return conn, nil
	}
	//fmt.Println(m.currentWait)
	if m.currentWait > m.maxWait {
		m.Unlock()
		return nil, WaitMaxError
	}
	m.currentWait++
	m.Unlock()
	//如果当前连接数大于最大连接数，则等待
	select {
	case <-time.After(m.waitTimeOut):
		m.currentWait--
		return nil, WaitTimeOutError
	case conn := <-m.ch:
		m.currentWait--
		return conn, nil
	}
}

func (m *neoPool) Release(conn Conn) {
	fmt.Println("release")
	m.Lock()
	//如果当前空闲数大于等于最大空闲数。则释放连接
	if len(m.ch) >= m.maxIdle {
		conn.Close()
		m.Unlock()
		return
	}
	m.ch <- conn
	m.Unlock()
}
func (m *neoPool) Ping(ctx context.Context) error {
	conn, err := m.Get()
	if err != nil {
		return err
	}
	//执行的最后释放连接
	defer m.Release(conn)
	//检查一下上级上下文是否正常存活。如果已经结束则直接返回
	select {
	case <-ctx.Done():
		return ContextDoneError
	default:
	}
	err = conn.Doit()
	return err
}

func NewNeoPool(factory func() (Conn, error)) NeoPool {
	pool := &neoPool{maxActive: 10, minActive: 2, maxIdle: 4,
		requestTimeOut: time.Second * 1, maxWait: 2,
		factory:     factory,
		waitTimeOut: time.Second * 4, ch: make(chan Conn, 10)}
	if pool.minActive > 0 {
		for i := 0; i < pool.minActive; i++ {
			conn := NewConn()
			pool.ch <- conn
			pool.currentActive++
		}
	}
	return pool
}
