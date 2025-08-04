package storage

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

func RedisInitConn() *redis.Pool {

	dialFunc := func() (redis.Conn, error) {
		conn, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			panic(fmt.Sprintf("redis Dial error: %v\n", err))
		}
		return conn, err
	}
	return &redis.Pool{
		Dial:        dialFunc,
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
	}

}
