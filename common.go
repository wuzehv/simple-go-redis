package simple_go_redis

import (
	"fmt"
	"net"
	"strconv"
)

type redisConn struct {
	addr string
	port int
	net.Conn
}

func New(addr string, port int) (*redisConn, error) {
	r := &redisConn{addr: addr, port: port}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return nil, err
	}

	r.Conn = conn
	return r, nil
}

func (r *redisConn) Do(params ...string) (interface{}, error) {
	res := new(redisProtocol)
	res.numOfParams = len(params)
	for _, v := range params {
		res.params = append(res.params, paramLen{
			len: len(v),
			val: v,
		})
	}

	_, err := r.Write([]byte(res.String()))
	if err != nil {
		return nil, err
	}

	return parseResponse(r.Conn)
}

func (r *redisConn) Select(db int) (interface{}, error) {
	return r.Do("select", strconv.Itoa(db))
}

func (r *redisConn) String() string {
	return r.addr + ":" + strconv.Itoa(r.port)
}
