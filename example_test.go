package simple_go_redis_test

import (
	"fmt"
	simple_go_redis "simple-go-redis"
)

const (
	IP   = "127.0.0.1"
	PORT = 6379
)

var r, _ = simple_go_redis.New(IP, PORT)

func ExampleRedisConn_Interactive() {
	r.Interactive()
}

func ExampleRedisConn_Select() {
	r.Select(2)
	r.Do("set", "k1", "v1")
	a, _ := r.Do("get", "v1")
	fmt.Printf("%s\n", a)
}
