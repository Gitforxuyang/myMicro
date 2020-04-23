package main

import (
	"context"
	"fmt"
	pool2 "myMicro/plugins/pool"
	"time"
)

type Data struct {
	Name int `json:"name"`
}

func main() {
	pool := pool2.NewNeoPool(func() (conn pool2.Conn, e error) {
		return pool2.NewConn(), nil
	})
	go func() {
		err := pool.Ping(context.TODO())
		fmt.Println(err)
	}()
	go func() {
		err := pool.Ping(context.TODO())
		fmt.Println(err)
	}()
	go func() {
		err := pool.Ping(context.TODO())
		fmt.Println(err)
	}()
	go func() {
		err := pool.Ping(context.TODO())
		fmt.Println(err)
	}()
	go func() {
		err := pool.Ping(context.TODO())
		fmt.Println(err)
	}()
	go func() {
		err := pool.Ping(context.TODO())
		fmt.Println(err)
	}()
	go func() {
		err := pool.Ping(context.TODO())
		fmt.Println(err)
	}()
	go func() {
		err := pool.Ping(context.TODO())
		fmt.Println(err)
	}()
	go func() {
		err := pool.Ping(context.TODO())
		fmt.Println(err)
	}()
	go func() {
		err := pool.Ping(context.TODO())
		fmt.Println(err)
	}()

	go func() {
		err := pool.Ping(context.TODO())
		fmt.Println(err)
	}()

	go func() {
		err := pool.Ping(context.TODO())
		fmt.Println(err)
	}()

	go func() {
		err := pool.Ping(context.TODO())
		fmt.Println(err)
	}()

	time.Sleep(time.Second * 1000)
}
