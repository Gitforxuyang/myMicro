package pool

import (
	"fmt"
	"time"
)

type Conn interface {
	Doit() error
	Close()
}

type conn struct {
}

func (m *conn) Close() {
	fmt.Println("conn close")
}

func (m *conn) Doit() error {
	fmt.Println("doit")
	time.Sleep(time.Second * 3)
	return nil
}
func NewConn() Conn {
	//fmt.Println("new")
	return &conn{}
}
