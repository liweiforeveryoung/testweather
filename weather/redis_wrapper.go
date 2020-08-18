package weather

import (
	"github.com/gomodule/redigo/redis"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 对 redis 进行一个简单的封装
// 来自 https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/05.6.md

var pool *redis.Pool

func init() {
	pool = newPool(redisAddr)
	close()
}

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("ping")
			return err
		},
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		Wait:        wait,
	}
}

func close() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		pool.Close()
		os.Exit(0)
	}()
}
